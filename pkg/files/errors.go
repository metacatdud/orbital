package files

import "errors"

var (
	ErrStatFile     = errors.New("stat file failed")
	ErrSrcIsDir     = errors.New("source is a directory")
	ErrSrcOpen      = errors.New("source open failed")
	ErrCreateBackup = errors.New("create backup failed")
	ErrRmBackup     = errors.New("remove backup failed")
	ErrCopyFile     = errors.New("copy file failed")
)
