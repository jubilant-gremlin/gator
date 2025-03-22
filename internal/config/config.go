package config

import (
	"encoding/json"
	"fmt"
	"os"
)
const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getConfigFilePath() (string, error){
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path + "/" + configFileName, nil
}


func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.WriteFile(path, encoded, 0666)
	return nil
}


func Read() Config {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}
	read, err := os.ReadFile(path)
	if err != nil {
		return Config{}
	}
	var cfg Config
	json.Unmarshal(read, &cfg)
	return cfg
}

func (c *Config) SetUser(name string) {
	c.Current_user_name = name
	err := write(*c)
	if err != nil {
		fmt.Printf("ERROR SETTING USER:%v", err)
	}
}
