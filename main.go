package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// Custom HTTP client
var client = http.Client{
	Timeout: 10 * time.Second,
}

func Main(Context openruntimes.Context) openruntimes.Response {
	// Initialize Fiber app
	app := fiber.New()

	// Define the route handler
	app.Get("/", func(c *fiber.Ctx) error {
		// Make an external HTTP request
		resp, err := client.Get("https://dummyjson.com/products/1")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		defer resp.Body.Close()

		// Check the HTTP status
		if resp.StatusCode != http.StatusOK {
			return c.Status(resp.StatusCode).JSON(&fiber.Map{
				"success": false,
				"error":   "Failed to fetch data",
			})
		}

		// Copy the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.Status(http.StatusOK).Send(body)
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

	// Execute Fiber app handler
	err := app.Test(&http.Request{
		Method: string(req.Header.Method()),
		URL:    req.URI(),
		Body:   io.NopCloser(bytes.NewReader(req.Body())),
	}, res)
	if err != nil {
		log.Println("Error handling request:", err)
		return Context.Res.Text("Internal Server Error", 500)
	}

	// Convert Fiber response to Appwrite response
	return openruntimes.Response{
		Body:    res.Body(),
		Headers: res.Header.Header(),
		Status:  res.StatusCode(),
	}
}
