package calculator

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	Command string `json:"command"`
	Number1 int64  `json:"number_1"`
	Number2 int64  `json:"number_2"`
}

func NewHandler(commandFunc CommandFunc, getRedisfunc GetRedisFunc, setRedisFunc SetRedisFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input Input

		err := c.BodyParser(&input)

		if err != nil {
			return fiber.NewError(999, "BodyParser")
		}
		key := fmt.Sprintf("%s:%v:%v", input.Command, input.Number1, input.Number2)

		result, err := getRedisfunc(c.Context(), key)
		if err != nil {
			fmt.Println("redis not found")
			result, err = commandFunc(input.Command, input.Number1, input.Number2)
			if err != nil {
				return fiber.NewError(999, err.Error())
			}
			err = setRedisFunc(c.Context(), key, result)
			if err != nil {
				fmt.Println("Cannot set setRedisFunc")
			}
			fmt.Println("set redis success")
		}
		fmt.Println("set found ", key, result)

		return c.SendString(fmt.Sprintf("%v", result))
	}
}

type CommandFunc func(command string, number1 int64, number2 int64) (int64, error)

type GetRedisFunc func(ctx context.Context, key string) (int64, error)

type SetRedisFunc func(ctx context.Context, key string, value int64) error

func NewCommandFunc() CommandFunc {
	return func(command string, number1 int64, number2 int64) (int64, error) {
		switch command {
		case "plus":
			return number1 + number2, nil
		case "minus":
			return number1 - number2, nil
		case "mul":
			return number1 * number2, nil
		case "devide":
			if number2 == 0 {
				return 0, errors.New("number2 much not equal 0")
			}
			return number1 + number2, nil
		default:
			return 0, errors.New("invalid command")
		}
	}
}

func NewGetRedisFunc(redisClient *redis.Client) GetRedisFunc {
	return func(ctx context.Context, key string) (int64, error) {
		return redisClient.Get(ctx, key).Int64()
	}
}

func NewSetRedisFunc(redisClient *redis.Client) SetRedisFunc {
	return func(ctx context.Context, key string, value int64) error {
		return redisClient.Set(ctx, key, value, 0).Err()
	}
}
