package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func UpdateArticle(db *sql.DB, article *Article) error {
	updateSQL := `UPDATE articles SET title=?, source=? WHERE id_article=?`

	statement, err := db.Prepare(updateSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(article.Title, article.Source, article.ID)
	return err
}
func CreateArticle(db *sql.DB, article *Article) (int64, error) {

	insertSQL := `INSERT INTO articles (title, source) VALUES (?,?);`

	statement, err := db.Prepare(insertSQL)

	if err != nil {
		return 0, err
	}

	defer statement.Close()

	fmt.Println(article)

	result, err := statement.Exec(article.Title, article.Source)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetAllArticles(db *sql.DB) ([]*Article, error) {
	query, err := db.Prepare("SELECT * FROM articles ORDER BY id_article DESC;")

	if err != nil {
		return nil, err
	}

	defer query.Close()

	result, err := query.Query()

	if err != nil {
		return nil, err
	}

	articles := make([]*Article, 0)

	for result.Next() {

		data := new(Article)

		err := result.Scan(&data.ID,
			&data.Title,
			&data.Source,
		)

		if err != nil {
			return nil, err
		}

		articles = append(articles, data)
	}

	return articles, nil
}

func GetArticleByID(db *sql.DB, id int) (*Article, error) {
	query, err := db.Prepare("SELECT * FROM articles WHERE id_article=?;")

	if err != nil {
		return nil, err
	}

	defer query.Close()

	result := query.QueryRow(id)

	data := new(Article)

	if result.Scan(&data.ID, &data.Title, &data.Source) != nil {
		return nil, err
	}

	return data, nil
}

func OpenDb() *sql.DB {
	log.Println("Opening ./.notez.db")

	db, err := sql.Open("sqlite3", "./.notez.db")

	if err != nil {
		panic(err)
	}

	CheckDbTables(db)

	return db
}

func CheckDbTables(db *sql.DB) {
	createArticlesSQL := `CREATE  TABLE IF NOT EXISTS articles(
		"id_article" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"source" TEXT);`

	if _, err := db.Exec(createArticlesSQL); err != nil {
		log.Fatal("Cannot create articles table")
	}
	log.Println("Article table check OK")
}
