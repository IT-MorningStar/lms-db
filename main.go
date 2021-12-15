package main

import (
	"lms-db/config"
	"lms-db/engine"
)

func main() {
	_ = engine.NewLmsDB(config.NewConfig())
}
