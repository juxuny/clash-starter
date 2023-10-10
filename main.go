package main

import (
	"flag"
	"github.com/juxuny/fs"
	"gopkg.in/yaml.v3"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultConfigFile := path.Join(homeDir, ".config", "clash-starter", "config.yaml")
	flag.StringVar(&starterConfigFile, "c", defaultConfigFile, "config file")
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
	if starterConfig.WorkingDirectory != "" {
		err = os.Chdir(starterConfig.WorkingDirectory)
		if err != nil {
			log.Println(err)
			os.Exit(255)
		}
	}
}

func start() {
	remoteConfig, err := fetchConfig(starterConfig.ForceRefreshFirst)
	if err != nil {
		log.Fatal(err)
	}
	clashConfigFileName := genClashConfigFileName()
	remoteConfig.Patch(starterConfig.GetAutoProxyGroup(), starterConfig.Override, starterConfig.Merge)
	remoteConfig.RunFilter(starterConfig.ProxyFilter)
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
	for range time.NewTicker(time.Second * time.Duration(starterConfig.ReloadInterval)).C {
		remoteConfig, err := fetchConfig(true)
		if err != nil {
			log.Println(err)
			continue
		}
		clashConfigFileName := genClashConfigFileName()
		remoteConfig.Patch(starterConfig.GetAutoProxyGroup(), starterConfig.GetOverride(), starterConfig.GetMerge())
		remoteConfig.RunFilter(starterConfig.ProxyFilter)
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
		// auto-remove oldest file
		autoCleanerHelper := fs.CreateFileCleaner(starterConfig.ConfigDir, func(fileName string, createTime time.Time, modifiedTime time.Time) bool {
			return path.Ext(fileName) == "yaml" || path.Ext(fileName) == "yml"
		})
		err = autoCleanerHelper.Execute(starterConfig.KeepNumOfFile, func(fileName string, createTime time.Time, modifiedTime time.Time) bool {
			return time.Now().Sub(createTime).Seconds() > float64(starterConfig.KeepDurationInSeconds)
		})
		if err != nil {
			log.Println(err)
		}
	}
}
