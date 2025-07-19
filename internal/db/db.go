package db

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

/*
   Package db provides SQLite connection management with automatic schema initialization.

   Example:
       db, err := db.Open("app.db")
       if err != nil {
           log.Fatal(err)
       }
       defer db.Close()

       // Use with SQLC
       queries := sqlc.New(db)
       user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
           Name:     "John Doe",
           Email:    "john@example.com", 
           Password: "hashed_password",
       })

   Notes:
   - Automatically executes embedded schema.sql on connection
   - Creates database file if it doesn't exist
   - Compatible with SQLC generated code
*/

//go:embed schema.sql
var ddl string

func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		return nil, err
	}

	return db, nil
}
