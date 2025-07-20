package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var config Config

type Config struct {
	Port  string `yaml:"port"`
	Path  string `yaml:"path"`
	Users []User `yaml:"users"`
}

type User struct {
	Key   string    `yaml:"key"`
	Read  *[]string `yaml:"read"`
	Write *[]string `yaml:"write"`
}

func findConfigFile() string {
	args := os.Args[1:]
	if len(args) > 0 && exists(args[0]) {
		return args[0]
	}

	if exists("objctr.yml") {
		return "objctr.yml"
	}

	if exists("~/objctr.yml") {
		return "~/objctr.yml"
	}

	return ""
}

func loadConfig() {
	configFile := findConfigFile()
	if configFile == "" {
		log.Fatal("Config file not found")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
	}

	println("< \033[1mOBJCTR SERVER\033[0m >")
	println()
	println("📄 Config Path: \t\033[1m", configFile, "\033[0m")
	println("🌐 Server Port: \t\033[1m", config.Port, "\033[0m")
	println("💾 Storage Path: \t\033[1m", config.Path, "\033[0m")
	println("--------------------------------")
}
