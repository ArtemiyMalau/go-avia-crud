package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
)

func GetDBConnection(config *config) (db *sqlx.DB, err error) {
	err = DoWithTries(func() (errInternal error) {
		db, errInternal = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?multiStatements=true", config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Database))
		log.Println(errInternal)
		return
	}, 5, time.Second*3)

	return
}

func InitDB(config *config, db *sqlx.DB) {
	migrateDb := flag.Bool("migratedb", false, "Initialize database's structure")
	seedDb := flag.Bool("seeddb", false, "Seeding database's data")
	flag.Parse()

	if *migrateDb {
		log.Printf("Initialize database's structure")
		mustExecSQLScript(db, filepath.Join(config.DB.Scripts, "structure.sql"))
	}

	if *seedDb {
		log.Printf("Seeding database's data")
		mustExecSQLScript(db, filepath.Join(config.DB.Scripts, "seeder.sql"))
	}
}

func mustExecSQLScript(db *sqlx.DB, path string) {
	tx := db.MustBegin()
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	tx.MustExec(string(buf))
	tx.Commit()
}
