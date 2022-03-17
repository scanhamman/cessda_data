package database

import (
	"context"
	_ "encoding/json"
	"fmt"
	_ "io/ioutil"
	"os"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func init() {
	settings = "../internal/database/db_settings.json"
}

func ObtainIds() {

}

func ScrapeDetails(idents []Identifier) {

	db_connstring := GetConnectionString("cessda")
	conn, err := pgx.Connect(context.Background(), db_connstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}
