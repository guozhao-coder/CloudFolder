package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

func GetFileSize(filename string) string {
	var result float32
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = float32(f.Size())
		return nil
	})
	return strconv.FormatFloat(float64(result/1000), 'f', 1, 64)
}
