package main

import (
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
    // Create a new Fiber app for testing
    app := fiber.New()
    
    // Add the health route
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "OK",
            "message": "Backend is running!",
        })
    })

    // Define test cases
    tests := []struct {
        description    string
        route          string
        expectedCode   int
        expectedBody   string
    }{
        {
            description:  "get HTTP status 200 from health endpoint",
            route:        "/health",
            expectedCode: 200,
            expectedBody: `{"message":"Backend is running!","status":"OK"}`,
        },
        {
            description:  "get HTTP status 404 from non-existing endpoint",
            route:        "/non-existing",
            expectedCode: 404,
            expectedBody: "Cannot GET /non-existing",
        },
    }

    // Iterate through test cases
    for _, test := range tests {
        // Create a new HTTP request
        req := httptest.NewRequest("GET", test.route, nil)
        
        // Perform the request
        resp, err := app.Test(req, -1)
        
        // Check for errors
        assert.NoError(t, err, test.description)
        
        // Check status code
        assert.Equal(t, test.expectedCode, resp.StatusCode, test.description)
    }
}

func TestAPITestEndpoint(t *testing.T) {
    // Create a new Fiber app for testing
    app := fiber.New()
    
    // Add the API test route
    app.Get("/api/test", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello from Go backend!",
            "mongo":   "Connected",
            "redis":   "Connected",
        })
    })

    // Create test request
    req := httptest.NewRequest("GET", "/api/test", nil)
    
    // Perform the request
    resp, err := app.Test(req, -1)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
