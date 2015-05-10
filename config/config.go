package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Smtp_Server     string
	Port           int
	Sender_Email    string
	Sender_Name     string
	Sender_Password string
	Execl_Path   string
	Subject string
}

func ReadConfig(configFilePath string) (*Config, error) {
	file, err := os.Open(configFilePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()


	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil
	}


	cfg := &Config{}
	err = yaml.Unmarshal([]byte(content), cfg)
	return cfg, err
}
