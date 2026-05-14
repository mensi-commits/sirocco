package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"sirocco-cluster/LoadClusterMetadata"
)

func main() {

	dsn := "admin:admin@tcp(127.0.0.1:3306)/sirocco_metadata?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	shardMap, err := cluster.LoadClusterMetadata(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cluster loaded:", len(shardMap.Shards), "shards")

	// -----------------------------
	// TEST ROUTING HERE
	// -----------------------------
	testKey := "user_42"

	route, err := cluster.Route(testKey, shardMap, false) // false = read
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ROUTE RESULT:")
	fmt.Println("Shard:", route.ShardID)
	fmt.Println("Node:", route.Host, route.Port)
	fmt.Println("Role:", route.Role)

	// keep alive
	select {}
}