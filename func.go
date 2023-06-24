package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func checkIfFileExists(fileName string) bool {
	stat, err := os.Stat(fileName)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func getTmpConfigFileName() string {
	return path.Join(starterConfig.ConfigDir, "tmp.config.yaml")
}

func genClashConfigFileName() string {
	return fmt.Sprintf("config.%s.yaml", time.Now().Format("20060102_150405"))
}

func reloadConfig(api string, fileName string) error {
	if strings.Index(fileName, "/") != -1 {
		workingDir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "get working dir failed")
		}
		fileName = path.Join(workingDir, fileName)
		log.Println("full file name:", fileName)
	}
	data, _ := json.Marshal(map[string]interface{}{
		"path": fileName,
	})
	buffer := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPut, api, buffer)
	if err != nil {
		return errors.Wrap(err, "create reload request failed")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sent request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 == 2 {
		_, _ = ioutil.ReadAll(resp.Body)
		return nil
	}
	return errors.Errorf("reload failed, status: %v", resp.Status)
}
