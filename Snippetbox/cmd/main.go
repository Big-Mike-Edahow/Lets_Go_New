/* main.go */

package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"snippetbox/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	snippets *models.SnippetModel
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
        snippets: &models.SnippetModel{DB: db},
    }

	http.HandleFunc("/", app.indexHandler)
	http.HandleFunc("/view", app.viewHandler)
	http.HandleFunc("/create", app.createHandler)
	http.HandleFunc("/save", app.saveHandler)
	http.HandleFunc("/edit", app.editHandler)
	http.HandleFunc("/update", app.updateHandler)
	http.HandleFunc("/delete", app.deleteHandler)
	http.HandleFunc("/about", app.aboutHandler)

	log.Printf("Serving HTTP on port%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, logRequest(http.DefaultServeMux)))
}
