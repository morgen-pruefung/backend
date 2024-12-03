package blogstore

import "backend/internal/blog"

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) GetArticles() ([]blog.Article, error) {
	return nil, nil
}

func (s *Store) GetArticle(articleID string) (*blog.Article, error) {
	return nil, nil
}
