package examstore

import (
	"backend/internal/exam"
	"backend/internal/github"
	"encoding/json"
)

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) GetExams() ([]exam.Exam, error) {
	files, err := github.ListFiles(github.BibliothekRepo, "exams")
	if err != nil {
		return nil, err
	}

	exams := []exam.Exam{}
	for _, file := range files {
		examID := file[:len(file)-5] // Remove ".json" suffix
		e, err := s.GetExam(examID)
		if err != nil {
			return nil, err
		}
		exams = append(exams, *e)
	}

	return exams, nil
}

func (s *Store) GetExam(examID string) (*exam.Exam, error) {
	metadata, err := github.ReadFile(github.BibliothekRepo, "exams/"+examID+".json")
	if err != nil {
		return nil, err
	}

	var e exam.Exam
	err = json.Unmarshal(metadata, &e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
