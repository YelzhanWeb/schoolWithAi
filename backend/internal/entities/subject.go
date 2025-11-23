package entities

import (
	"errors"
)

type Subject struct {
	ID     string
	Slug   string
	NameRu string
	NameKz string
}

func NewSubject(id, slug, nameRu, nameKz string) (*Subject, error) {
	if slug == "" || nameRu == "" || nameKz == "" {
		return nil, errors.New("all subject fields are required")
	}
	return &Subject{
		ID:     id,
		Slug:   slug,
		NameRu: nameRu,
		NameKz: nameKz,
	}, nil
}
