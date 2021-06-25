package config

import "log"

type Application struct {
	Port    int
	Env     string
	Logger  *log.Logger
	Version string
	DB      struct {
		Driver string
		DSN    string
	}
	JWT struct {
		Secret string
	}
}
