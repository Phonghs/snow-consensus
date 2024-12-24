package main

import (
	"SnowConsensus/api/route"
	"SnowConsensus/bootstrap"
	"SnowConsensus/p2p"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	rd := app.Redis
	defer app.CloseRedisConnection()

	var wg sync.WaitGroup
	nodes := app.SetupNodes()
	for _, n := range nodes {
		wg.Add(1)
		go func(n *p2p.SnowNode) {
			defer wg.Done()
			r := gin.Default()
			route.Setup(env, rd, r, n)

			addr := fmt.Sprintf("%s:%s", n.NodeProfile.IP, n.NodeProfile.Port)
			log.Println("Node", n.NodeProfile.ID, "is listening on", addr)
			if err := r.Run(addr); err != nil {
				log.Printf("Error starting server on %s: %v\n", addr, err)
			}
		}(n)
	}
	wg.Wait()
}
