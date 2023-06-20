package main

import (
	"flag"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

var (
	configFile    string
	verbose       bool
	currentConfig StarterConfig
)

func init() {
	flag.StringVar(&configFile, "c", "config.yaml", "config file")
	flag.BoolVar(&verbose, "v", false, "output debug info")
}

func initConfig() {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("open config file failed:", err)
	}
	err = yaml.Unmarshal(configData, &currentConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("load config file success: ", configFile)
	log.Println("clash config dir:", currentConfig.ConfigDir)
	log.Println("clash bin file:", currentConfig.Bin)
}

func start() {
	_, err := fetchConfig(false)
	if err != nil {
		log.Fatal(err)
	}
	command := exec.Command(currentConfig.Bin, "-d", currentConfig.ConfigDir, "-f", path.Join(currentConfig.ConfigDir, "config.yaml"))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	initConfig()
	checkAndPrepare()
	start()
}
