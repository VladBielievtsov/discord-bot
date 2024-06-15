package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func ConnectDatabase() {
	connStr := "user=postgres dbname=discordbot password=password sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	createTableDiscordMessages(DB)
	// dropTabel("discord_messages")
}

func dropTabel(table string) {
	query := "DROP TABLE IF EXISTS " + table + ";"
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to drop table: %v", err)
	}
}

func createTableDiscordMessages(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS discord_messages(
		id SERIAL PRIMARY KEY,
		payload JSON NOT NULL,
		user_id BIGINT NOT NULL);
	`
	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func AddDiscordMessages(jbytes []byte, userID string) (int64, error) {
	query := "INSERT INTO discord_messages (payload, user_id) VALUES ($1,$2) RETURNING id;"

	var lastInserted int64
	err := DB.QueryRow(query, string(jbytes), userID).Scan(&lastInserted)
	if err != nil {
		return 0, err
	}

	return lastInserted, nil
}
