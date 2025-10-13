package redis

import "backend/internal/ports"

type redis struct {
}

func NewRedis() ports.Redis {
	return &redis{}
}
