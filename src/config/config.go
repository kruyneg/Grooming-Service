package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env string `yaml:"env"`
	LogPath string `yaml:"log-path"`
	DBConnection string `yaml:"dbconnection"`
	Port int `yaml:"port"`
	TemplatesPath string `yaml:"templates-path"`
	StaticPath string `yaml:"static-path"`
}

func Load() (conf *Config) {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "./configs/dev.yml", "path to config file")
	flag.Parse()
	
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		log.Panicf("error while opening config file: %s", err)
	}

	// Parse file
	conf = &Config{
		Env: "dev",
		Port: 8080,
		TemplatesPath: "/templates",
		StaticPath: "/static",
	}

	if err = yaml.NewDecoder(cfgFile).Decode(&conf); err != nil {
		log.Panic(err)
	}

	if conf.DBConnection == "" {
		log.Panic("Configuration field 'dbconnection' is requred!")
	}

	return
}