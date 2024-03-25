package media

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Pineapple217/mb/config"
	"github.com/Pineapple217/mb/database"
)

const UploadDir = config.DataDir + "/uploads"
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var invalidCharsR = regexp.MustCompile(`[^\w\s-.]+`)

func CreateUploadDir() {
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		err := os.Mkdir(UploadDir, 0755)
		if err != nil {
			slog.Error("Failed to create directory",
				"dir", UploadDir,
				"error", err,
			)
			return
		}
		fmt.Println("Upload directory created successfully")
	}
}

func SaveFile(ctx context.Context, f *multipart.FileHeader, customName string) error {
	// TODO: handle corrupt or malicious files
	src, err := f.Open()
	if err != nil {
		return err
	}
	// Name
	var name string
	var ext string
	if customName == "" {
		name, ext = splitFileNameAndExtension(f.Filename)
	} else {
		name = customName
		ext = getFileExtension(f.Filename)
	}

	// Destination
	fPath := makeValidFileName(name) + "." + ext
	fPathFull := filepath.Join(UploadDir, fPath)

	_, err = os.Stat(fPathFull)

	i := 0
	for !os.IsNotExist(err) && i < 100 {
		fPath = makeValidFileName(name) + "_" + generateRandom(3) + "." + ext
		fPathFull = filepath.Join(UploadDir, fPath)
		_, err = os.Stat(fPathFull)
		i++
	}

	dst, err := os.Create(fPathFull)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	defer src.Close()

	// Add to DB
	// TODO: use inmemory file instead of rereading file
	imgF, err := os.Open(fPathFull)
	if err != nil {
		return err
	}
	queries := database.GetQueries()
	_, err = queries.CreateMediafile(
		ctx,
		database.CreateMediafileParams{
			FileName:      name,
			FilePath:      fPath,
			FileType:      getFileType(f.Filename),
			FileExtention: ext,
			Thumbnail:     thumbnail(imgF),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFile(filename string) error {
	f := filepath.Join(UploadDir, filename)
	err := os.Remove(f)
	if os.IsNotExist(err) {
		slog.Warn("File to delete not found", "file", f)
		return nil
	}
	if err != nil {
		slog.Info("Delete failed", "file", f, "error", err)
		return err
	}
	slog.Debug("Upload deleted", "file", f)
	return nil
}

func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	// Remove the dot from the extension
	if ext != "" {
		ext = ext[1:]
	}
	return ext
}

func splitFileNameAndExtension(filename string) (name, ext string) {
	name = filename[:len(filename)-len(filepath.Ext(filename))]
	ext = filepath.Ext(filename)
	// Remove the dot from the extension
	if ext != "" {
		ext = ext[1:]
	}
	return name, ext
}

func getFileType(file string) string {
	ext := strings.ToLower(filepath.Ext(file))
	// TODO add gif and bmp support
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return "image"
	case ".mp4", ".avi", ".mkv", ".mov", ".wmv":
		return "video"
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "audio"
	default:
		return "unknown"
	}
}

func makeValidFileName(input string) string {
	processedString := invalidCharsR.ReplaceAllString(input, "_")
	processedString = strings.Trim(processedString, "_")

	return processedString
}

func generateRandom(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
