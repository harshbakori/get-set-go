package main

import (
	"encoding/json"
	"fmt"
	_ "get-set-go/docs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Config struct {
	Redis struct {
		Address string `json:"address"`
	} `json:"redis"`
	Kafka struct {
		BootstrapServers string `json:"bootstrap_servers"`
	} `json:"kafka"`
}

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	initRedis(config.Redis.Address)

	initKafka(config.Kafka.BootstrapServers)
}

func main() {
	logFile, err := os.OpenFile("requests.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	e := echo.New()
	e.GET("/api/verve/accept", acceptHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	go logUniqueRequests()
	e.Logger.Fatal(e.Start(":8081"))
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
		return c.String(http.StatusBadRequest, "Failed")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed")
	}

	endpoint := c.QueryParam("endpoint")

	if isUniqueRequest(id) {
		if endpoint != "" {
			go sendPostRequest(endpoint)
		} else {
			log.Printf("No endpoint provided for notification")
		}
		return c.String(http.StatusOK, "ok")
	} else {
		return c.String(http.StatusOK, "Duplicate request")
	}
}

// func isUniqueRequest(id int) bool {
// 	val, err := rdb.SetNX(ctx, fmt.Sprintf("request_id_%d", id), 1, 1*time.Minute).Result()
// 	if err != nil {
// 		log.Printf("Failed to check/set unique ID in Redis: %v", err)
// 		return false
// 	}
// 	return val
// }

// func logUniqueRequests() {
// 	for {
// 		time.Sleep(1 * time.Minute)

// 		count, err := rdb.DBSize(ctx).Result()
// 		if err != nil {
// 			log.Printf("Failed to get unique request count from Redis: %v", err)
// 			continue
// 		}
// 		log.Printf("Unique request count in the last minute: %d", count)

// 		sendToKafka(count)

// 		rdb.FlushDB(ctx)
// 	}
// }

func sendPostRequest(endpoint string) {
	count, err := rdb.DBSize(ctx).Result()
	if err != nil {
		log.Printf("Failed to get unique request count from Redis: %v", err)
		return
	}

	resp, err := http.PostForm(endpoint, url.Values{"count": {fmt.Sprintf("%d", count)}})
	if err != nil {
		log.Printf("Failed to send POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("POST request to %s returned status code: %d", endpoint, resp.StatusCode)
}
