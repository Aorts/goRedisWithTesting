package main

import (
	"calculator/calculator"
	"calculator/config"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	cfg, err := initConfig()
	if err != nil {
		fmt.Println("cannot init config")
		panic(err.Error())
	}

	redisClient := initRedisClient(cfg.Redis)

	app := fiber.New()

	app.Get("/version", func(c *fiber.Ctx) error {
		return c.SendString(cfg.Version)
	})

	app.Get("/calculator", calculator.NewHandler(calculator.NewCommandFunc(), calculator.NewGetRedisFunc(redisClient), calculator.NewSetRedisFunc(redisClient)))

	app.Listen(cfg.Server.Port)
}

func initConfig() (*config.Config, error) {
	viper.SetDefault("Version", "v0.0.1")
	viper.SetDefault("SERVER.PORT", ":8080")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg config.Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initRedisClient(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Host + cfg.Port,
		DB:   cfg.DB,
	})

	return client
}
