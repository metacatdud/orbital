package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BACKUP for files. Files are backup in place
// - BackupWithTimestamp - will back up a file with format: <name>_<timestamp>_origExt.bkp
// - Backup - will back up a file in place with format: <name>.origExt.bkp
// - CleanupBackups - will remove ".bkp" files on a given path

func BackupWithTimestamp(srcPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("%w:[%s]", ErrStatFile, err.Error())
	}

	if info.IsDir() {
		return fmt.Errorf("%w:[%s]", ErrSrcIsDir, srcPath)
	}

	ts := time.Now().Format("20060102150405") // YYYYMMDDhhmmss
	dir, base := filepath.Dir(srcPath), filepath.Base(srcPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	ext = strings.TrimPrefix(ext, ".")

	backupName := fmt.Sprintf("%s_%s_%s.bkp", name, ts, ext)
	backupPath := filepath.Join(dir, backupName)

	return copyFile(srcPath, backupPath, info.Mode())
}

func Backup(srcPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("%w:[%s]", ErrStatFile, err.Error())
	}

	if info.IsDir() {
		return fmt.Errorf("%w:[%s]", ErrSrcIsDir, srcPath)
	}

	backupPath := srcPath + ".bkp"
	return copyFile(srcPath, backupPath, info.Mode())
}

func CleanupBackups(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) == ".bkp" {
			if rmErr := os.Remove(path); rmErr != nil {
				return fmt.Errorf("%w:[%s] %v", ErrRmBackup, path, rmErr)
			}
		}
		return nil
	})
}

// copyFile helper that copies a file from src to dst, preserving permissions.
func copyFile(src, dst string, perm os.FileMode) error {

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("%w:[%s]: %v", ErrSrcOpen, src, err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("%w:[%s]: %v", ErrCreateBackup, dst, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("%w:[%s] %v", ErrCopyFile, dst, err)
	}
	return nil
}
