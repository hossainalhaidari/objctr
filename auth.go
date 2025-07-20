package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func hash(key string) string {
	hashBytes := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hashBytes[:])
}

func hasPermission(permissions *[]string, path string) bool {
	if permissions == nil {
		return false
	}

	for _, permission := range *permissions {
		prefixPath := permission
		if !strings.HasSuffix(permission, "/") {
			prefixPath += "/"
		}
		if permission == path || strings.HasPrefix(path, prefixPath) {
			return true
		}
	}

	return false
}

func canRead(key string, path string) bool {
	if path == "" {
		return false
	}

	hashed := hash(key)

	if hasPermission(config.Users[0].Read, path) {
		return true
	}

	for _, user := range config.Users {
		if hashed == user.Key {
			return hasPermission(user.Read, path)
		}
	}

	return false
}

func canWrite(key string, path string) bool {
	if path == "" {
		return false
	}

	hashed := hash(key)

	if hasPermission(config.Users[0].Write, path) {
		return true
	}

	for _, user := range config.Users {
		if hashed == user.Key {
			return hasPermission(user.Write, path)
		}
	}

	return false
}
