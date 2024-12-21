package blogstore

import (
	"backend/internal/blog"
	"backend/internal/github"
	"encoding/json"
)

const BibliothekRepo = "bibliothek"

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}
func (s *Store) GetArticles() ([]blog.Article, error) {
	files, err := github.ListFiles(BibliothekRepo, "blog/articles")
	if err != nil {
		return nil, err
	}

	articles := []blog.Article{}
	for _, file := range files {
		a, err := s.GetArticle(file)
		if err != nil {
			return nil, err
		}
		articles = append(articles, *a)
	}

	return articles, nil
}

func (s *Store) GetArticle(articleID string) (*blog.Article, error) {
	metadata, err := github.ReadFile(BibliothekRepo, "blog/articles/"+articleID+"/"+articleID+".json")
	if err != nil {
		return nil, err
	}

	var article blog.Article
	err = json.Unmarshal(metadata, &article)
	if err != nil {
		return nil, err
	}

	content, err := github.ReadFile(BibliothekRepo, "blog/articles/"+articleID+"/"+articleID+".md")
	if err != nil {
		return nil, err
	}

	article.Body = string(content)

	return &article, nil
}
