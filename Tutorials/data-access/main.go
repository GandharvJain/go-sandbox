package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/url"
	"os"
)

type Album struct {
	ID int64
	Title, Artist string
	Price float32
}

var ctxBg context.Context = context.Background()

func main() {
	db_url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		url.QueryEscape(os.Getenv("DBUSER")),
		url.QueryEscape(os.Getenv("DBPASS")),
		"127.0.0.1",
		"5432",
		"recordings",
	)
	conn, err := pgx.Connect(ctxBg, db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctxBg)
	pingErr := conn.Ping(ctxBg)
	if pingErr != nil {
		fmt.Fprintf(os.Stderr, "Couldn't ping database: %v\n", pingErr)
		os.Exit(1)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist(conn, "Taylor Swift")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Albums not found: %v\n", err)
	}
	fmt.Printf("Albums found: %#v\n", albums)
}

func albumsByArtist(conn *pgx.Conn, name string) ([]Album, error) {
	var albums []Album

	rows, err := conn.Query(ctxBg, "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, err
	}
	albums, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (Album, error) {
		var alb Album
		err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		return alb, err
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}
