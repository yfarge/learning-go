package main

import (
	"context"
	"data_access/configuration"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	config, err := configuration.GetConfiguration()
	if err != nil {
		slog.Error("Error loading configuration", "err", err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.New(
		context.Background(),
		config.Database.ConnectionString(),
	)

	if err != nil {
		slog.Error("Failed to establish connection to database", "err", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	pingErr := dbpool.Ping(context.Background())
	if pingErr != nil {
		slog.Error("Failed to get greeting from database.", "err", err)
		os.Exit(1)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane", dbpool)
	if err != nil {
		slog.Error("Failed to get albums from John Coltrane from database.", "err", err)
		os.Exit(1)
	}

	fmt.Println(albums)
}

func albumsByArtist(name string, dbpool *pgxpool.Pool) ([]Album, error) {
	var albums []Album

	rows, err := dbpool.Query(
		context.Background(),
		"SELECT * FROM album WHERE artist = $1",
		name,
	)

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}
