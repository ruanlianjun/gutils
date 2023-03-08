package gutils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !stat.IsDir()
}
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	return stat.IsDir()
}

func MkdirAll(path string, options ...MkDirOptions) error {
	exists := FileExists(path)
	if !exists {

		if dir := IsDir(path); !dir {
			path = filepath.Dir(path)
		}

		fmt.Println("======path:", path)
		fileMode := newMkFileMode()
		for _, item := range options {
			item(fileMode)
		}

		if fileMode.recursion {
			if err := os.MkdirAll(path, os.FileMode(fileMode.fileMode)); err != nil {
				return err
			}
		} else {
			if err := os.Mkdir(path, os.FileMode(fileMode.fileMode)); err != nil {
				return err
			}
		}

	}
	return nil
}

type MkDirOptions func(*MkFileMode)

type MkFileMode struct {
	fileMode  uint32
	recursion bool //是否递归创建
}

const defaultMkFileMode = 0777

func newMkFileMode() *MkFileMode {
	return &MkFileMode{fileMode: defaultMkFileMode, recursion: true}
}

func WithFileMode(fileMode uint32) MkDirOptions {
	return func(mode *MkFileMode) {
		mode.fileMode = fileMode
	}
}

func WithMkNotRecursion() MkDirOptions {
	return func(mode *MkFileMode) {
		mode.recursion = false
	}
}
