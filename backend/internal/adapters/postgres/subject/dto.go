package subject

import (
	"backend/internal/entities"
)

type dto struct {
	ID     string `db:"id"`
	Slug   string `db:"slug"`
	NameRu string `db:"name_ru"`
	NameKz string `db:"name_kz"`
}

func (d *dto) toEntity() entities.Subject {
	return entities.Subject{
		ID:     d.ID,
		Slug:   d.Slug,
		NameRu: d.NameRu,
		NameKz: d.NameKz,
	}
}
