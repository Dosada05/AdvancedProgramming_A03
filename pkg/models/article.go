package models

import "database/sql"

type Article struct {
	id      int
	Title   string
	Content string
}

func CreateArticle(db *sql.DB, article Article) error {
	_, err := db.Exec("INSERT INTO articles (title, content) VALUES (?, ?)", article.Title, article.Content)
	if err != nil {
		return err
	}
	return nil
}

func GetLatestArticles(db *sql.DB) ([]Article, error) {
	rows, err := db.Query("SELECT title, content FROM articles ORDER BY id DESC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.Title, &article.Content); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
