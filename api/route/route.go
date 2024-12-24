package route

import (
	"SnowConsensus/api/controller"
	"SnowConsensus/bootstrap"
	"SnowConsensus/p2p"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Setup(env *bootstrap.Env, rd *redis.Client, gin *gin.Engine, node p2p.Validator) {
	publicRouter := gin.Group("")
	base := controller.NewBase(env, node, rd)

	publicRouter.GET("/ping", base.Ping)

	publicRouter.POST("/transaction/create", base.CreateTransaction)

	publicRouter.POST("/transaction/query", base.VerifyTransaction)
}
