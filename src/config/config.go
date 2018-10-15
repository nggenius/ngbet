package config

import (
	"encoding/json"
	"log"

	"github.com/mysll/toolkit"
)

type Config struct {
	Recommend []string `json:"recommend"`
	Broadcast []string `json:"broadcast"`
}

var (
	Setting Config
)

func LoadConfig() {
	data, err := toolkit.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(data, &Setting)
	if err != nil {
		log.Fatalln(err)
	}
}
