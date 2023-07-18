package main

import (
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

func fetchConfigData(force bool) (data []byte, err error) {
	if !force {
		exists := checkIfFileExists(getTmpConfigFileName())
		if exists {
			data, err = ioutil.ReadFile(getTmpConfigFileName())
			return
		}
	}
	resp, err := http.Get(starterConfig.Link)
	if err != nil {
		return nil, errors.Wrap(err, "fetch subscribe config failed")
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.Errorf("error http status, code: %d, error: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "network I/O error")
	}
	err = ioutil.WriteFile(getTmpConfigFileName(), data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func fetchConfig(force bool) (config *ClashConfig, err error) {
	data, err := fetchConfigData(force)
	if err != nil {
		return nil, errors.Wrap(err, "fetch config file failed")
	}
	config = &ClashConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, errors.Wrap(err, "invalid yaml data")
	}
	return
}

func saveConfig(config *ClashConfig, fileName string) (err error) {
	if config == nil {
		return errors.Wrap(err, "invalid config data")
	}
	log.Println("found proxy: ", len(config.GetProxies()))
	for _, proxy := range config.GetProxies() {
		log.Println("found proxy name: ", proxy.Name)
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "marshal yaml failed")
	}
	return ioutil.WriteFile(fileName, data, 0644)
}
