/* handlers.go */

package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		log.Println(err)
	}

	indexTemplate := template.Must(template.ParseFiles("./templates/index.html"))
	err = indexTemplate.Execute(w, snippets)
	if err != nil {
		log.Println(err)
	}
}

func (app *application) viewHandler(w http.ResponseWriter, r *http.Request) {
	snippetId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(snippetId)

	snippet, err := app.snippets.Get(id)
	if err != nil {
		log.Println(err)
	}

	viewTemplate := template.Must(template.ParseFiles("./templates/view.html"))
	err = viewTemplate.Execute(w, snippet)
	if err != nil {
		log.Println(err)
	}
}

func (app *application) createHandler(w http.ResponseWriter, r *http.Request) {
	createTemplate := template.Must(template.ParseFiles("./templates/create.html"))
	createTemplate.Execute(w, nil)
}

func (app *application) saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		type Message struct {
			Title   string
			Content string
			Errors  map[string]string
		}
		msg := &Message{
			Title:   r.PostFormValue("title"),
			Content: r.PostFormValue("content"),
		}
		msg.Errors = make(map[string]string)

		if strings.TrimSpace(title) == "" {
			msg.Errors["Title"] = "Title required"
		}
		if strings.TrimSpace(content) == "" {
			msg.Errors["Content"] = "Content required"
		}
		if msg.Errors["Title"] != "" || msg.Errors["Content"] != "" {
			createTemplate, _ := template.ParseFiles("./templates/create.html")
			createTemplate.Execute(w, msg)
		} else {
			_, err := app.snippets.Insert(title, content)
			if err != nil {
				log.Println(err)

			}
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (app *application) editHandler(w http.ResponseWriter, r *http.Request) {
	snippetId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(snippetId)

	snippet, err := app.snippets.Get(id)
	if err != nil {
		log.Println(err)
	}

	type Message struct {
		Id      int
		Title   string
		Content string
		Errors  map[string]string
	}
	msg := &Message{
		Id:      snippet.Id,
		Title:   snippet.Title,
		Content: snippet.Content,
		Errors:  make(map[string]string),
	}

	editTemplate := template.Must(template.ParseFiles("./templates/edit.html"))
	editTemplate.Execute(w, msg)
}

func (app *application) updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		snippetId := r.FormValue("id")
		id, _ := strconv.Atoi(snippetId)

		title := r.FormValue("title")
		content := r.FormValue("content")

		type Message struct {
			Id      int
			Title   string
			Content string
			Errors  map[string]string
		}
		msg := &Message{
			Id:      id,
			Title:   r.PostFormValue("title"),
			Content: r.PostFormValue("content"),
		}
		msg.Errors = make(map[string]string)

		if strings.TrimSpace(title) == "" {
			msg.Errors["Title"] = "Title required"
		}
		if strings.TrimSpace(content) == "" {
			msg.Errors["Content"] = "Content required"
		}
		if msg.Errors["Title"] != "" || msg.Errors["Content"] != "" {
			editTemplate, _ := template.ParseFiles("./templates/edit.html")
			editTemplate.Execute(w, msg)
		} else {
			err := app.snippets.Update(id, title, content)
			if err != nil {
				log.Println(err)
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (app *application) deleteHandler(w http.ResponseWriter, r *http.Request) {
	snippetId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(snippetId)

	err := app.snippets.Delete(id)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request) {
	aboutTemplate, _ := template.ParseFiles("./templates/about.html")
	aboutTemplate.Execute(w, nil)
}
