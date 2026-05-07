package main

import (
	"sirocco-agent/api"
	"sirocco-agent/config"
)

func main() {
	cfg := config.Load()

	server := api.NewServer(cfg)
	server.Start()
}