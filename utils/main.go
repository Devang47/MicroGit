package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

const (
	DEFAULT_PATH = ".microgit"
)

type SavePoint struct {
	Message   string            `json:"message"`
	Timestamp string            `json:"timestamp"`
	Parent    string            `json:"parent"`
	Files     map[string]string `json:"files"`
}

// hashContent returns the SHA-256 hash of the file content
func HashContent(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

// writeObject saves the file content to objects/<hash>
func WriteObject(hash string, content []byte) error {
	objectPath := filepath.Join(DEFAULT_PATH, "objects", hash)
	return os.WriteFile(objectPath, content, 0644)
}
