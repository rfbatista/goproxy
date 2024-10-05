package shared

import (
	"path/filepath"
	"runtime"
)

func FindProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")
	return string(rootPath)
}
