package main

import (
	"embed"
	"log"
	"os"
)

//go:embed bin
var binFs embed.FS

const (
	clashCommandLineFile = "bin/clash-linux_amd64"
)

func checkAndPrepare() {
	if stat, err := os.Stat(starterConfig.Bin); err == nil && !stat.IsDir() {
		log.Println("clash verified")
		return
	}
}
