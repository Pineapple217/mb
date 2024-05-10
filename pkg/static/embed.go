package static

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

var (
	//go:embed public/*
	PublicFS    embed.FS
	excludeExts = []string{".woff2"}
	StaticMap   map[string]string
)

func hashFile(file fs.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := hash.Sum(nil)

	return hex.EncodeToString(sum)[:12], nil
}

func hashFiles(f embed.FS) (map[string]string, error) {
	filesDetails := map[string]string{}

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := filepath.Ext(path)
			for _, e := range excludeExts {
				if ext == e {
					return nil
				}
			}

			file, err := f.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			hash, err := hashFile(file)
			if err != nil {
				return err
			}

			fileName := strings.TrimSuffix(d.Name(), ext)

			newFilename := fmt.Sprintf("%s-%s%s", fileName, hash, ext)
			basePath := strings.TrimSuffix(path, d.Name())
			basePath = strings.ReplaceAll(basePath, "public", "/static")

			p := strings.ReplaceAll(path, "public", "/static")
			filesDetails[p] = basePath + newFilename
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesDetails, nil
}

func HashPublicFS() map[string]string {
	slog.Info("Hashing static files")
	files, err := hashFiles(PublicFS)
	if err != nil {
		panic(err)
	}

	for k, v := range files {
		slog.Debug("file hashed", "old", k, "new", v)
	}
	StaticMap = files
	swappedMap := make(map[string]string)
	for k, v := range files {
		swappedMap[v] = k
	}
	return swappedMap
}
