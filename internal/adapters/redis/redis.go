package redis

import "github.com/YelzhanWeb/schoolWithAi/internal/ports"

type redis struct {
}

func NewRedis() ports.Redis {
	return &redis{}
}
