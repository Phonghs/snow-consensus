package bootstrap

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	AppEnv            string `mapstructure:"APP_ENV"`
	RedisHost         string `mapstructure:"REDIS_HOST"`
	RedisPort         string `mapstructure:"REDIS_PORT"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	ContextTimeout    int    `mapstructure:"CONTEXT_TIMEOUT"`
	SampleSize        int    `mapstructure:"SAMPLE_SIZE"`
	QuorumSize        int    `mapstructure:"QUORUM_SIZE"`
	DecisionThreshold int    `mapstructure:"DECISION_THRESHOLD"`
	CountNode         int    `mapstructure:"COUNT_NODE"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile("../.env,example")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env,example : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
