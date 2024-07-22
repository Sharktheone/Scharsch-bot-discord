package conf

import (
	"encoding/json"
	"fmt"
	"github.com/Sharktheone/ScharschBot/config"
	"github.com/Sharktheone/ScharschBot/flags"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	confPath = flags.String("configPath")
	Config   *Format
)

func init() {
	GetConf()
}

func GetConf() *Format {
	ymlConf, err := os.ReadFile(*confPath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := config.GetDefaultConf()
			if err != nil {
				log.Fatalf("Failed to get default config: %v", err)
			}
			if err := os.WriteFile(*confPath, f, 0644); err != nil {
				log.Fatalf("Failed to write default config: %v", err)
			}
			fmt.Printf("No config found, created default config at %s", *confPath)
			os.Exit(0)
		}
		log.Fatalf("Failed to get config: %v", err)
	}
	if err := yaml.Unmarshal(ymlConf, &Config); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	configJSON, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	log.Println("Config loaded:\n", string(configJSON))

	return Config
}
