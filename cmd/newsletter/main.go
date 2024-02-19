package main

import (
	"context"
	"genuinebnt/newsletter/config"
	lib "genuinebnt/newsletter/internal"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	connectionString := config.Database.ConnectionString()

	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalln("Failed to connect to database: ", err)
	}
	defer dbpool.Close()

	address := config.Application.Host + ":" + strconv.Itoa(config.Application.Port)

	server := lib.Server(dbpool)
	server.Run(address)
}
