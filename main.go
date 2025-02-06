package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	_ "get-set-go/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // swagger middleware
)

// @title Verve Service API
// @version 1.0
// @description This is a sample API for Verve service.
// @host localhost:8080
// @BasePath /

var (
	uniqueRequests = make(map[int]struct{})
	mu             sync.Mutex
)

func main() {
	// Create a log file
	logFile, err := os.OpenFile("requests.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set the logger to write to the log file
	log.SetOutput(logFile)

	e := echo.New()
	e.GET("/api/verve/accept", acceptHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	go logUniqueRequests()
	e.Logger.Fatal(e.Start(":8080"))
}

// acceptHandler handles the GET request
// @Summary Accept a request
// @Description Accepts an integer ID and an optional endpoint to notify
// @Param id query int true "Request ID"
// @Param endpoint query string false "Notification Endpoint"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "id is required or must be an integer"
// @Router /api/verve/accept [get]
func acceptHandler(c echo.Context) error {
	log.Printf("Received a request with id: %s", c.QueryParam("id"))
	idParam := c.QueryParam("id")
	if idParam == "" {
		// return c.String(http.StatusBadRequest, "id is required")
		return c.String(http.StatusBadRequest, "Failed")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		// return c.String(http.StatusBadRequest, "id must be an integer")
		return c.String(http.StatusBadRequest, "Failed")
	}

	endpoint := c.QueryParam("endpoint")

	mu.Lock()
	uniqueRequests[id] = struct{}{}
	mu.Unlock()

	if endpoint != "" {
		go sendPostRequest(endpoint)
	} else {
		log.Printf("No endpoint provided for notification")
	}

	return c.String(http.StatusOK, "ok")
}

func logUniqueRequests() {
	for {
		time.Sleep(1 * time.Minute)

		mu.Lock()
		count := len(uniqueRequests)
		log.Printf("Unique request count in the last minute: %d", count)

		uniqueRequests = make(map[int]struct{})
		mu.Unlock()
	}
}

func sendPostRequest(endpoint string) {
	mu.Lock()
	count := len(uniqueRequests)
	mu.Unlock()
	fmt.Println("count", count)

	resp, err := http.PostForm(endpoint, url.Values{"count": {fmt.Sprintf("%d", count)}})
	if err != nil {
		log.Printf("Failed to send POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("POST request to %s returned status code: %d", endpoint, resp.StatusCode)
}
