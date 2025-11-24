package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Subject struct {
	ID     string
	Slug   string
	NameRu string
	NameKz string
}

func NewSubject(slug, nameRu, nameKz string) (*Subject, error) {
	if slug == "" || nameRu == "" || nameKz == "" {
		return nil, errors.New("all subject fields are required")
	}
	return &Subject{
		ID:     uuid.NewString(),
		Slug:   slug,
		NameRu: nameRu,
		NameKz: nameKz,
	}, nil
}
