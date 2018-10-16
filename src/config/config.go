package config

import (
	"encoding/json"
	"log"

	"github.com/mysll/toolkit"
)

type Rule struct {
	Update []string            `json:"update"`
	State  map[string][]string `json:"state"`
}

type Bet365 struct {
	WSURL string `json:"wsurl"`
	Host  string `json:"host"`
}

type Config struct {
	Bet365    Bet365   `json:"bet365"`
	Recommend []string `json:"recommend"`
	Broadcast []string `json:"broadcast"`
	Rule      Rule     `json:"rule"`
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

	log.Println(Setting)
}
