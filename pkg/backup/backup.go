package backup

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/database"
)

func GetAllBackups() ([]string, error) {
	var backups []string

	directoryPath := "./data/backups"
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			backups = append(backups, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(backups, func(i, j int) bool {
		return extractDateFromFilename(backups[j]).Before(extractDateFromFilename(backups[i]))
	})

	return backups, nil
}

const layout = "2006-01-02_15-04-05"

func extractDateFromFilename(filename string) time.Time {
	s := strings.TrimLeft(filename, "backup_")
	s = strings.TrimRight(s, "tar.gz")
	t, err := time.Parse(layout, s)
	if err != nil {
		slog.Warn("failed to parse date", "error", err)
	}
	return t
}

func Backup(ctx context.Context, q *database.Queries) error {
	file, err := q.Backup(ctx)
	if err != nil {
		return err
	}
	defer os.Remove(file)

	timeString := time.Now().Format("2006-01-02_15-04-05")
	tarballName := fmt.Sprintf("%s/backup_%s.tar.gz", config.BackupDir, timeString)

	tarballFile, err := os.Create(tarballName)
	if err != nil {
		return err
	}
	defer tarballFile.Close()

	gzipWriter := tar.NewWriter(tarballFile)
	defer gzipWriter.Close()

	if err := addFileToTar(gzipWriter, file, ""); err != nil {
		return err
	}

	folderPath := config.UploadDir
	baseUploadDir := filepath.Base(filepath.Clean(folderPath))

	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filePath == folderPath {
			return nil
		}

		relPath, err := filepath.Rel(folderPath, filePath)
		if err != nil {
			return err
		}
		headerName := filepath.Join(baseUploadDir, filepath.ToSlash(relPath))
		return addFileToTar(gzipWriter, filePath, headerName)
	})

	return err
}

func addFileToTar(gzipWriter *tar.Writer, filePath, headerName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}

	if headerName != "" {
		header.Name = headerName
	} else {
		header.Name = fileInfo.Name()
	}

	if err := gzipWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(gzipWriter, file)
	return err
}
