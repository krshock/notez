package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Article struct {
	ID     int
	Title  string
	Source string
}

var ViewArticleWithID = regexp.MustCompile(`/view_article`)

func main() {
	db := OpenDb()
	defer db.Close()

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " - " + r.URL.Path)

		//fmt.Fprint(w, "Hello world")
		articles, err := GetAllArticles(db)

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		html_articles := make([]map[string]interface{}, 0)

		for _, art := range articles {
			data := map[string]interface{}{
				"ID":     art.ID,
				"Title":  art.Title,
				"Source": template.HTML(art.Source),
			}
			html_articles = append(html_articles, data)
		}

		tmpl, err := template.ParseFiles("templates/index.html")

		if err != nil {
			fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, html_articles)
	})

	http.HandleFunc("/view_article/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " - " + r.URL.Path)

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/view_article/"))

		if err != nil {
			fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("Reading article id: " + strconv.Itoa(id))

		article, err := GetArticleByID(db, id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
			return
		}

		tmpl, err := template.ParseFiles("templates/view_article.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
			return
		}

		html_article := map[string]interface{}{
			"ID":     article.ID,
			"Title":  article.Title,
			"Source": template.HTML(article.Source),
		}
		tmpl.Execute(w, html_article)
	})

	http.HandleFunc("/edit_article/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " - " + r.URL.Path)
		if r.Method == http.MethodPost {

			r.ParseForm()

			art_id, err := strconv.Atoi(r.FormValue("id_article"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
				return
			}

			article := new(Article)
			article.ID = art_id
			article.Title = r.FormValue("title")
			article.Source = r.FormValue("source")

			err = UpdateArticle(db, article)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
				return
			}

			http.Redirect(w, r, "/view_article/"+fmt.Sprint(article.ID), 301)
			return

		} else {
			id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/edit_article/"))

			if err != nil {
				fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			article, err := GetArticleByID(db, id)

			if err != nil {
				fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			html_article := map[string]interface{}{
				"ID":     article.ID,
				"Title":  article.Title,
				"Source": template.HTML(article.Source),
			}

			tmpl, err := template.ParseFiles("templates/edit_article.html")

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error: "+fmt.Sprint(err))
				return
			}

			tmpl.Execute(w, html_article)
		}
	})

	http.ListenAndServe(":8000", nil)
}
