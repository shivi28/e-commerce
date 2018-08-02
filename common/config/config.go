package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	gcfg "gopkg.in/gcfg.v1"
)

var CF *Config

type DatabaseConfig struct {
	Master string
	Slave  string
}

type ServerConfig struct {
	Host         string
	Port         int
	TemplatePath string
	Timeout      time.Duration
}

type Config struct {
	Server   ServerConfig
	Database map[string]*DatabaseConfig
}

func GetConfig() *Config {
	return CF
}

func init() {
	CF = &Config{}
	GOPATH := os.Getenv("GOPATH")
	ok := ReadConfig(CF, GOPATH+"/src/github.com/e-commerce/files/etc", "e-commerce") || ReadConfig(CF, "files/etc", "e-commerce")
	if !ok {
		log.Fatal("Failed to read config file")
	}
}

func ReadConfig(cfg *Config, path string, module string) bool {

	var configString []string

	fname := path + "/" + module + "/main." + "ini"

	config, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Println("common/config.go function ReadConfig", err)
		return false
	}

	configString = append(configString, string(config))

	err = gcfg.ReadStringInto(cfg, strings.Join(configString, "\n\n"))
	if err != nil {
		log.Println("common/config.go function ReadConfig", err)
		return false
	}

	return true
}
