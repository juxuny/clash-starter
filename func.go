package main

import (
	"os"
	"path"
)

func checkIfFileExists(fileName string) bool {
	stat, err := os.Stat(fileName)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func getTmpConfigFileName() string {
	return path.Join(currentConfig.ConfigDir, "tmp.config.yaml")
}
