package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"net/url"
	"os"
)

type Album struct {
	ID            int64
	Title, Artist string
	Price         float32
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
		os.Exit(1)
	}
	fmt.Printf("Albums found: %#v\n", albums)

	alb, err := albumByID(conn, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Album not found: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Album found: %#v\n", alb)

	albID, err := addAlbum(conn, Album{
		Title:  "The Days / Nights",
		Artist: "Avicii",
		Price:  49.99,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Album not added: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ID of added album: %v\n", albID)

	albums, err = allAlbums(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Albums not found: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All Albums:")
	for _, album := range albums {
		fmt.Printf("%#v\n", album)
	}
}

func albumsByArtist(conn *pgx.Conn, name string) ([]Album, error) {
	var albums []Album

	rows, err := conn.Query(ctxBg, "SELECT * FROM album WHERE artist = $1", name)
	//if err != nil {
	//	return nil, err
	//}
	//albums, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (Album, error) {
	//	var alb Album
	//	err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	//	return alb, err
	//})
	albums, err = pgx.CollectRows(rows, pgx.RowToStructByName[Album])
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func albumByID(conn *pgx.Conn, id int64) (Album, error) {
	var alb Album

	rows, _ := conn.Query(ctxBg, "SELECT * FROM album WHERE id = $1", id)
	alb, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Album])
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("No rows found with id %d", id)
			return alb, err
		}
		return alb, err
	}
	return alb, nil
}

func addAlbum(conn *pgx.Conn, alb Album) (int64, error) {
	var lastID int64
	err := conn.QueryRow(
		ctxBg,
		"INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		alb.Title,
		alb.Artist,
		alb.Price,
	).Scan(&lastID)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func allAlbums(conn *pgx.Conn) ([]Album, error) {
	var albums []Album

	rows, err := conn.Query(ctxBg, "SELECT * FROM album")
	albums, err = pgx.CollectRows(rows, pgx.RowToStructByName[Album])
	if err != nil {
		return nil, err
	}
	return albums, nil
}
