package main

import (
	"github.com/dev-mackan/gowiki/internal/apiserver"
	"github.com/dev-mackan/gowiki/internal/db"
	"github.com/dev-mackan/gowiki/internal/repos"
	"github.com/dev-mackan/gowiki/internal/reposervice"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	dbConfig := db.DbConfig{
		Addr:         "file:store.sqlite?cache=shared&journal_mode=WAL&busy_timeout=3000",
		MaxOpenConns: 30,
		MaxIdleConns: 30,
		MaxIdleTime:  "15m",
	}
	db, err := db.New(&dbConfig)
	if err != nil {
		panic(err)
	}
	repo := repos.NewSqlRepository(db)
	repoService := reposervice.NewRepoService(repo)
	api := apiserver.NewAPIServer(":3000", repoService)
	log.Fatal(api.Run())
}
