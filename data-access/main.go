package main

import (
	"context"
	"data_access/configuration"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	config, err := configuration.GetConfiguration()
	if err != nil {
		slog.Error("Error loading configuration", "err", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), config.Database.ConnectionString())
	if err != nil {
		slog.Error("Failed to establish connection to database", "err", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		slog.Error("Failed to get greeting from database.", "err", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
