package calculator

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/magiconair/properties/assert"
)

func newFiberApp() *fiber.App {
	app := fiber.New()

	return app
}

func TestPlusMuskOK(t *testing.T) {
	commandFunc := NewCommandFunc()
	params := []byte(`
		{
			"command" :"plus",
			"number_1" : 2,
			"number_2" : 1
		}`)

	var input Input
	_ = json.Unmarshal(params, &input)
	result, _ := commandFunc(input.Command, input.Number1, input.Number2)
	expectResult := int64(3)

	assert.Equal(t, expectResult, result)
}

func TestNewHandler(t *testing.T) {
	getRedisFunc := func(ctx context.Context, key string) (int64, error) {
		return 5, nil
	}

	setRedisFunc := func(ctx context.Context, key string, value int64) error {
		return nil
	}

	params := []byte(`
		{
			"command" :"plus",
			"number_1" : 2,
			"number_2" : 1
		}`)

	body := bytes.NewReader(params)

	handler := NewHandler(NewCommandFunc(), getRedisFunc, setRedisFunc)

	app := newFiberApp()
	app.Post("/", handler)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)

	assert.Equal(t, err, nil)
	assert.Equal(t, res.StatusCode, http.StatusOK)

}
