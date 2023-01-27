package config

import (
	"encoding/json"
	"os"
)

var (
	instance = loadConfig()
)

func PortNumber() string      { return instance.PortNumber }
func BackupDirectory() string { return instance.BackupDirectory }

type config struct {
	PortNumber      string `json:"portNumber"`
	BackupDirectory string `json:"backupDirectory"`
}

func loadConfig() *config {
	file, err := os.Open("./config.json")
	if err != nil {
		panic("`config.json` doesn't exist")
	}

	c := new(config)
	json.NewDecoder(file).Decode(&c)
	return c
}
