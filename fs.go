package main

import (
	"fmt"
	"io"
	"io/fs"
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
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return !info.IsDir()
}

func list(path string) []string {
	entries := []string{}
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err == nil {
			relativePath := strings.TrimPrefix(path, config.Path)
			if info.IsDir() {
				relativePath += "/"
			}
			entries = append(entries, relativePath)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return entries
}

func createDir(path string) bool {
	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func uploadFile(file *multipart.FileHeader, outputPath string) bool {
	src, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer src.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func move(fromPath string, toPath string) bool {
	err := os.Rename(fromPath, toPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func copy(fromPath string, toPath string) bool {
	if isFile(fromPath) {
		originalFile, err := os.Open(fromPath)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer originalFile.Close()

		newFile, err := os.Create(toPath)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, originalFile)

		if err != nil {
			fmt.Println(err)
			return false
		}

		return true
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

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func delete(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
