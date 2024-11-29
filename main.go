package handler

import (
	"io"
	"log"
	"net/http"
	"os"
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
		resp, err := client.Get("https://dummyjson.com/products/1")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return c.Status(resp.StatusCode).JSON(&fiber.Map{
				"success": false,
				"error":   "Failed to fetch data",
			})
		}

		if _, err := io.Copy(c.Response().BodyWriter(), resp.Body); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return nil
	})

	// Convert Appwrite request to Fiber context
	req := fiber.Request{}
	req.SetRequestURI(Context.Req.URL)
	req.Header.SetMethod(Context.Req.Method)
	for key, value := range Context.Req.Headers {
		req.Header.Set(key, value)
	}
	req.SetBody(Context.Req.Body)

	// Handle the request using Fiber
	res := fiber.Response{}
	err := app.Handler()(&req, &res)
	if err != nil {
		log.Println("Error handling request:", err)
		return Context.Res.Text("Internal Server Error", http.StatusInternalServerError)
	}

	// Convert Fiber response to Appwrite response
	headers := make(map[string]string)
	res.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return openruntimes.Response{
		Body:    string(res.Body()),
		Headers: headers,
		Status:  res.StatusCode(),
	}
}
