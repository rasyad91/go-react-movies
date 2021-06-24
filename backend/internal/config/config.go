package config

import "log"

type Application struct {
	Port    int
	Env     string
	Logger  *log.Logger
	Version string
}
