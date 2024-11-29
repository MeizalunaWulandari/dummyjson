kpackage handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// Struct for API response
type ApiResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`

	Data    interface{} `json:"data,omitempty"`
}

// Dummy user data
var dummyUser = map[string]interface{}{
	"id":       1,
	"name":     "John Doe",
	"email":    "johndoe@example.com",
	"username": "johndoe",
}

func Main(Context openruntimes.Context) openruntimes.Response {
	// Initialize Fiber app
	app := fiber.New()

	// Define endpoint: /ping
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(ApiResponse{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "pong",
		})
	})

	// Define endpoint: /user
	app.Get("/user", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(ApiResponse{
			Code:   fiber.StatusOK,
			Status: "success",
			Data:   dummyUser,
		})
	})

	// Fiber context handling
	// Create a Fiber context for the incoming Appwrite request
	c := app.AcquireCtx(&fiber.Ctx{})
	defer app.ReleaseCtx(c)

	// Map Appwrite request data to Fiber context
	c.Request().Header.SetMethod(Context.Req.Method)
	c.Request().SetRequestURI(Context.Req.Url)
	for key, value := range Context.Req.Headers {
		c.Request().Header.Set(key, value)
	}
	if bodyBytes, ok := Context.Req.Body().([]byte); ok {
		c.Request().SetBody(bodyBytes)
	}

	// Handle the request using Fiber
	app.Handler()(c)

	// Convert Fiber response to Appwrite response
	headers := make(map[string]string)
	c.Response().Header.VisitAll(func(k, v []byte) {
		headers[string(k)] = string(v)
	})

	return openruntimes.Response{
		Body:    string(c.Response().Body()),
		Headers: headers,
		Status:  c.Response().StatusCode(),
	}
}

