package controller

import (
	"SnowConsensus/bootstrap"
	"SnowConsensus/p2p"
	"github.com/redis/go-redis/v9"
)

type Base struct {
	Env   *bootstrap.Env
	Node  p2p.Validator
	Redis *redis.Client
}

func NewBase(env *bootstrap.Env, node p2p.Validator, redis *redis.Client) *Base {
	return &Base{
		Env:   env,
		Node:  node,
		Redis: redis,
	}
}
