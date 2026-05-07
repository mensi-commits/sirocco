package main

import (
	"log"
	"sirocco-cluster/api"
	"sirocco-cluster/cluster"
)

func main() {
	cluster.Init()

	server := api.NewServer()
	log.Println("Sirocco Cluster running on :9000")
	server.Start(":9000")
}