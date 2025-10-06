package postgre

import (
	"database/sql"

	"github.com/YelzhanWeb/schoolWithAi/internal/ports"
)

type Pool struct {
	DB *sql.DB
}

func NewPoolDB(db *sql.DB) ports.Postgre {
	return &Pool{
		DB: db,
	}
}
