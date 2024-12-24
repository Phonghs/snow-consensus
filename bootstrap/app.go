package bootstrap

import (
	"SnowConsensus/p2p"
	"SnowConsensus/util"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"sync"
)

type Application struct {
	Env   *Env
	Redis *redis.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Redis = NewRedisClient(app.Env)
	return *app
}

func (app *Application) CloseRedisConnection() {
	CloseRedisClient(app.Redis)
}

func (app *Application) SetupNodes() []*p2p.SnowNode {
	countNode := app.Env.CountNode
	profileMap := make(map[string]*p2p.Profile, countNode)
	maxPort := 8000 + countNode
	for i := 8000; i < maxPort; i++ {
		log.Println("Setup node", i)
		prive, pub, err := util.GenerateRSAKeys(512)

		if err != nil {
			log.Println("Can't generate RSA keys", err)
			os.Exit(1)
		}
		nodeProfile := p2p.Profile{
			ID:   fmt.Sprint(i),
			IP:   "localhost",
			Port: fmt.Sprint(i),
			Key: p2p.Key{
				PublicKey:  pub,
				PrivateKey: prive,
				Algorithm:  "RSA",
				Hash:       "SHA256",
			},
		}
		profileMap[fmt.Sprint(i)] = &nodeProfile
	}
	nodes := make([]*p2p.SnowNode, countNode)
	for i := 0; i < countNode; i++ {

		nodeProfile := profileMap[fmt.Sprint(8000+i)]
		nodes[i] = &p2p.SnowNode{
			NodeProfile:     nodeProfile,
			PeerNodeInfoMap: util.GetMapExcludingKey(profileMap, nodeProfile.ID),
			Transaction:     make(map[string]*p2p.Transaction),
			Env: p2p.Env{
				SampleSize:        app.Env.SampleSize,
				QuorumSize:        app.Env.QuorumSize,
				DecisionThreshold: app.Env.DecisionThreshold,
				TimeOut:           app.Env.ContextTimeout,
			},
			Mutex: sync.Mutex{},
		}
	}
	return nodes
}
