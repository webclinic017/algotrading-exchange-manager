package dbticks

import (
	"context"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
)

func DbInit() bool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return false
	}
	defer conn.Close(ctx)

	var greeting string

	err = conn.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return false
	}
	fmt.Println(greeting)
	return true
}
