package main

import (
    "os"
    "context"
    "fmt"
    "net/http"
    "github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379" // Default value
    }

    rdb = redis.NewClient(&redis.Options{
        Addr: redisAddr, // Redis server address
    })

    // Test Redis connection
    if err := rdb.Ping(ctx).Err(); err != nil {
        panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
    }
}

func hello(w http.ResponseWriter, r *http.Request) {
    incr := rdb.Incr(ctx, "counter")
    count, err := incr.Result()
    if err != nil {
        http.Error(w, "Error incrementing counter", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "hello, my name is Falcon. You've visited %d times!", count)
}

func main() {
    http.HandleFunc("/", hello)
    port := ":4000"
    fmt.Println("Server running at http://localhost" + port)
    http.ListenAndServe(port, nil)
}
