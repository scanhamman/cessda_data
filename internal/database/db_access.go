package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

var settings string // string for relative path of settings file

func init() {
	settings = "../internal/database/db_settings.json"
}

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
}

func GetConnectionString(db_name string) string {

	content, err := ioutil.ReadFile(settings)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var c Credentials
	err = json.Unmarshal(content, &c)
	if err != nil {
		os.Exit(1)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, db_name)

}

func StoreIdentifiers(idents []Identifier) {

	db_connstring := GetConnectionString("cessda")
	conn, err := pgx.Connect(context.Background(), db_connstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	SliceAsRows := [][]interface{}{}
	for _, i := range idents {
		SliceAsRows = append(SliceAsRows, []interface{}{i.Id, i.Status})
	}
	fmt.Println("created slice as rows")
	var columns = []string{"id", "status"}
	var tablename = []string{"da", "identifiers"}

	_, err = conn.CopyFrom(
		context.Background(),
		tablename,
		columns,
		pgx.CopyFromRows(SliceAsRows))

	if err != nil {
		os.Exit(1)
	}
	fmt.Println("rows stored in db")
}
