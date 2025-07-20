package main

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func join(path string) string {
	return filepath.Join(config.Path, strings.ReplaceAll(path, "..", ""))
}

func isRoot(path string) bool {
	return path == "/"
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func list(path string) (bool, []string) {
	dirs := []string{}
	files := []string{}
	entries, err := os.ReadDir(path)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() {
			dirs = append(dirs, entry.Name()+"/")
		} else {
			files = append(files, entry.Name())
		}
	}
	return err == nil, append(dirs, files...)
}

func createDir(path string) bool {
	err := os.Mkdir(path, 0755)
	return err == nil
}

func uploadFile(file *multipart.FileHeader, outputPath string) bool {
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		return false
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return false
	}

	return true
}

func move(fromPath string, toPath string) bool {
	err := os.Rename(fromPath, toPath)
	return err == nil
}

func copy(fromPath string, toPath string) bool {
	if isFile(fromPath) {
		originalFile, err := os.Open(fromPath)
		if err != nil {
			return false
		}
		defer originalFile.Close()

		newFile, err := os.Create(toPath)
		if err != nil {
			return false
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, originalFile)
		return err == nil
	}

	err := filepath.Walk(fromPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(toPath, path)
		if info.IsDir() {
			os.MkdirAll(destPath, info.Mode())
			return nil
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		if _, err := io.Copy(dest, src); err != nil {
			return err
		}

		if err := os.Chmod(destPath, info.Mode()); err != nil {
			return err
		}

		return nil
	})

	return err == nil
}

func delete(path string) bool {
	err := os.RemoveAll(path)
	return err == nil
}
