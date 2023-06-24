package main

import (
	"flag"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

var (
	starterConfigFile string
	verbose           bool
	starterConfig     StarterConfig
)

func init() {
	flag.StringVar(&starterConfigFile, "c", "config.yaml", "config file")
	flag.BoolVar(&verbose, "v", false, "output debug info")
}

func initConfig() {
	configData, err := ioutil.ReadFile(starterConfigFile)
	if err != nil {
		log.Fatal("open config file failed:", err)
	}
	err = yaml.Unmarshal(configData, &starterConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("load config file success: ", starterConfigFile)
	log.Println("clash config dir:", starterConfig.ConfigDir)
	log.Println("clash bin file:", starterConfig.Bin)
}

func start() {
	remoteConfig, err := fetchConfig(false)
	if err != nil {
		log.Fatal(err)
	}
	clashConfigFileName := genClashConfigFileName()
	remoteConfig.Patch(starterConfig.GetAutoProxyGroup(), starterConfig.Override, starterConfig.Merge)
	err = saveConfig(remoteConfig, path.Join(starterConfig.ConfigDir, clashConfigFileName))
	if err != nil {
		panic(err)
	}
	log.Println("use config:", clashConfigFileName)
	command := exec.Command(starterConfig.Bin, "-d", starterConfig.ConfigDir, "-f", path.Join(starterConfig.ConfigDir, clashConfigFileName))
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
	go start()
	for range time.NewTimer(time.Second * time.Duration(starterConfig.ReloadInterval)).C {
		remoteConfig, err := fetchConfig(true)
		if err != nil {
			log.Println(err)
			continue
		}
		clashConfigFileName := genClashConfigFileName()
		remoteConfig.Patch(starterConfig.GetAutoProxyGroup(), starterConfig.GetOverride(), starterConfig.GetMerge())
		err = saveConfig(remoteConfig, path.Join(starterConfig.ConfigDir, clashConfigFileName))
		if err != nil {
			log.Println(err)
			continue
		}
		clashConfigFullFileName := path.Join(starterConfig.ConfigDir, clashConfigFileName)
		err = reloadConfig(remoteConfig.GetControlPanelEntrypoint()+"/configs", path.Join(starterConfig.ConfigDir, clashConfigFileName))
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("reload config file success:", clashConfigFullFileName)
	}
}
