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

// TODO: cleanup
func Backup(ctx context.Context, q *database.Queries) error {
	file, err := q.Backup(ctx)
	if err != nil {
		return err
	}

	timeString := time.Now().Format("2006-01-02_15-04-05")
	tarballName := fmt.Sprintf("%s/backup_%s.tar.gz", config.BackupDir, timeString)

	tarballFile, err := os.Create(tarballName)
	if err != nil {
		return err
	}
	defer tarballFile.Close()

	gzipWriter := tar.NewWriter(tarballFile)
	defer gzipWriter.Close()

	dbBackup, err := os.Open(file)
	if err != nil {
		return err
	}
	defer dbBackup.Close()

	dbBackupStats, err := dbBackup.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(dbBackupStats, "")
	if err != nil {
		return err
	}
	header.Name = dbBackupStats.Name()

	if err := gzipWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(gzipWriter, dbBackup); err != nil {
		return err
	}

	dbBackup.Close()
	if err = os.Remove(file); err != nil {
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

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(folderPath, filePath)
		if err != nil {
			return err
		}
		header.Name = baseUploadDir + "/" + strings.Replace(relPath, "\\", "/", -1)

		if err := gzipWriter.WriteHeader(header); err != nil {
			return err
		}

		if _, err := io.Copy(gzipWriter, file); err != nil {
			return err
		}

		return nil
	})

	return err
}
