package main

import (
	"proxy_manager/config"
	"proxy_manager/internal/app"
)

func main() {
	cfg := config.NewConfig()
	app.Run(&cfg)
}
