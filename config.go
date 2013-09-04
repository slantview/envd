package main

import (
	"launchpad.net/goyaml"
)

type Config struct {
	servers []string
	key string
	cert string
	cacert string
}

func NewConfig(file string) *Config {
	
}