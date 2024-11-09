package conf

import (
	"cmp"
	"fmt"
	"github.com/Sharktheone/ScharschBot/config"
	"github.com/Sharktheone/ScharschBot/flags"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"slices"
)

var (
	confPath = flags.String("configPath")
	Config   *Format
)

func LoadConf() {
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

	sortRoleConf()
}

func sortRoleConf() {
	slices.SortStableFunc(Config.Whitelist.RolesConfig, func(i, j RoleConfig) int {
		return cmp.Compare(i.Priority, j.Priority)
	})
}
