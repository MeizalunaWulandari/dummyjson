kkpackage handler

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

	// Endpoint: /ping
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(ApiResponse{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "pong",
		})
	})

	// Endpoint: /user
	app.Get("/user", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(ApiResponse{
			Code:   fiber.StatusOK,
			Status: "success",
			Data:   dummyUser,
		})
	})

	// Convert Appwrite request to Fiber request
	req := fiber.AcquireRequest()
	defer fiber.ReleaseRequest(req)

	req.Header.SetMethod(Context.Req.Method)
	req.SetRequestURI(Context.Req.Url)
	for key, value := range Context.Req.Headers {
		req.Header.Set(key, value)
	}
	bodyBytes := Context.Req.Body().([]byte)
	req.SetBody(bodyBytes)

	// Prepare Fiber response
	res := fiber.AcquireResponse()
	defer fiber.ReleaseResponse(res)

	// Execute Fiber handler
	err := app.Test(req, res)
	if err != nil {
		// Error handling
		return Context.Res.Json(ApiResponse{
			Code:    500,
			Status:  "fail",
			Message: "Internal Server Error",
		})
	}

	// Convert Fiber response to Appwrite response
	headers := make(map[string]string)
	res.Header.VisitAll(func(k, v []byte) {
		headers[string(k)] = string(v)
	})

	return openruntimes.Response{
		Body:    res.Body(),
		Headers: headers,
		Status:  res.StatusCode(),
	}
}

