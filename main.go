package main

import (
	"net/http"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 && args[0] == "hash" {
		println(hash(args[1]))
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get
		var query = r.URL.Query().Get
		var path = r.URL.Path
		var fullPath = join(path)
		var toPath = query("to")
		var toFullPath = join(toPath)
		var key = header("Authorization")

		switch r.Method {
		case "GET":
			if !canRead(key, path) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !exists(fullPath) {
				http.Error(w, "Path does not exist", http.StatusNotFound)
				return
			}

			if isFile(fullPath) {
				http.ServeFile(w, r, fullPath)
			} else {
				success, entries := list(fullPath)
				if !success {
					http.Error(w, "Failed to list directory", http.StatusInternalServerError)
					return
				}

				for _, entry := range entries {
					w.Write([]byte(entry + "\n"))
				}
			}
		case "POST":
			if !canWrite(key, path) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if isRoot(path) {
				http.Error(w, "Cannot create root path", http.StatusBadRequest)
				return
			}

			if strings.HasPrefix(header("Content-Type"), "multipart/form-data") {
				r.ParseMultipartForm(10 << 20)
				_, handler, err := r.FormFile("file")

				if err != nil {
					http.Error(w, "File upload error", http.StatusBadRequest)
					return
				}

				success := uploadFile(handler, fullPath)
				if !success {
					http.Error(w, "Failed to delete path", http.StatusInternalServerError)
					return
				}
			} else {
				success := createDir(fullPath)
				if !success {
					http.Error(w, "Failed to delete path", http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusNoContent)
		case "PATCH":
			if !canWrite(key, path) || !canWrite(key, toPath) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if isRoot(path) || isRoot(toPath) {
				http.Error(w, "Cannot move root path", http.StatusBadRequest)
				return
			}

			if !exists(fullPath) {
				http.Error(w, "Path does not exist", http.StatusNotFound)
				return
			}

			success := move(fullPath, toFullPath)
			if !success {
				http.Error(w, "Failed to delete path", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		case "PUT":
			if !canWrite(key, path) || !canWrite(key, toPath) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if isRoot(path) || isRoot(toPath) {
				http.Error(w, "Cannot copy root path", http.StatusBadRequest)
				return
			}

			if !exists(fullPath) {
				http.Error(w, "Path does not exist", http.StatusNotFound)
				return
			}

			success := copy(fullPath, toFullPath)
			if !success {
				http.Error(w, "Failed to delete path", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		case "DELETE":
			if !canWrite(key, path) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if isRoot(path) {
				http.Error(w, "Cannot delete root path", http.StatusBadRequest)
				return
			}

			if !exists(fullPath) {
				http.Error(w, "Path does not exist", http.StatusNotFound)
				return
			}

			success := delete(fullPath)
			if !success {
				http.Error(w, "Failed to delete path", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	})

	loadConfig()
	http.ListenAndServe(":3000", nil)
}
