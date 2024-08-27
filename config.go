package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port           int
	PingIntervalMs int
	MaxFails       int
	ReportUrl      string
}

var config = readConfig()

func readConfig() *Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Could not read config file")
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Could not parse config")
	}
	return &config
}
