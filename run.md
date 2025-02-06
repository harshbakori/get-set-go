
## install echo
`go get -u github.com/labstack/echo/v4`

## generate doc
`go run github.com/swaggo/swag/cmd/swag@latest init`

## install redis package : `Extension 1`
`go get github.com/go-redis/redis/v8`

## install kafka : `Extension 2`
`go get github.com/confluentinc/confluent-kafka-go/kafka`

## run redpanda for kafka mimic using docker
`docker run --network host redpandadata/redpanda`

## run go server
`go run main.go`