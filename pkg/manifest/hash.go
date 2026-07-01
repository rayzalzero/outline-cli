package manifest

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// FileHash computes MD5 hash of file contents
func FileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}
