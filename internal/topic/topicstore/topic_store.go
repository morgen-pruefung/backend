package topicstore

import (
	"backend/internal/github"
	"backend/internal/topic"
	"encoding/json"
)

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) GetTopics() ([]topic.Topic, error) {
	files, err := github.ListFiles(github.BibliothekRepo, "topics")
	if err != nil {
		return nil, err
	}

	topics := []topic.Topic{}
	for _, file := range files {
		t, err := s.GetTopic(file)
		if err != nil {
			return nil, err
		}
		topics = append(topics, *t)
	}

	return topics, nil
}

func (s *Store) GetTopic(topicID string) (*topic.Topic, error) {
	metadata, err := github.ReadFile(github.BibliothekRepo, "topics/"+topicID+"/topic.json")
	if err != nil {
		return nil, err
	}

	var t topic.Topic
	err = json.Unmarshal(metadata, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Store) GetTextContent(topicID string) ([]byte, error) {
	content, err := github.ReadFile(github.BibliothekRepo, "topics/"+topicID+"/text-content.md")
	if err != nil {
		return nil, err
	}

	return content, nil
}
