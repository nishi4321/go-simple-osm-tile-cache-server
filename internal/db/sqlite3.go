package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS tiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			x INTEGER,
			y INTEGER,
			z INTEGER,
			r TEXT,
			tile_data BLOB
		);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTile(x, y, z int, r string, tileData []byte) error {
	insertQuery := `INSERT INTO tiles (x, y, z, r, tile_data) VALUES (?, ?, ?, ?, ?);`

	_, err := db.Exec(insertQuery, x, y, z, r, tileData)
	if err != nil {
		return err
	}

	return nil
}

func GetTile(x, y, z int, r string) ([]byte, error) {
	selectQuery := `SELECT tile_data FROM tiles WHERE x = ? AND y = ? AND z = ? AND r = ?;`

	var tileData []byte
	err := db.QueryRow(selectQuery, x, y, z, r).Scan(&tileData)
	if err != nil {
		return nil, err
	}

	return tileData, nil
}
