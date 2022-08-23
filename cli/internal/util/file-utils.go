package util

import "os"

func FileExist(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}
