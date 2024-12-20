package main

import (
	"encoding/json"
	"os"
)

var Config struct {
	Port      int    `json:"port"`
	Token     string `json:"token"`
	BuyOnly   bool   `json:"buy_only"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Demo      bool   `json:"demo"`
}

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
