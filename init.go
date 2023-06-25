package main

import (
	"embed"
	"io/ioutil"
	"log"
	"os"
	"path"
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
	binDir := path.Dir(starterConfig.Bin)
	if stat, err := os.Stat(binDir); err == nil && !stat.IsDir() {
		log.Println(binDir, " is not a directory")
		os.Exit(255)
	}
	err := os.MkdirAll(binDir, 0755)
	if err != nil {
		panic(err)
	}
	data, err := binFs.ReadFile(clashCommandLineFile)
	if err != nil {
		log.Println("load clash binary failed:", err)
		os.Exit(255)
	}
	err = ioutil.WriteFile(starterConfig.Bin, data, 0755)
	if err != nil {
		log.Println("init clash binary failed:", err)
		os.Exit(255)
	}
}
