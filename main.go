package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDatabase() *gorm.DB {
	// connect to database
	// DATABASE_URL := "postgres://localhost:5432/practice?sslmode=disable"
	DATABASE_URL := "host=localhost dbname=sample_company port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{}) // os.Getenv("DATABASE_URL")
	if err != nil {
		log.Fatal(err)
		fmt.Println("Database connection failed ❌")
	}
	fmt.Println("Database connection successful ✅")
	return db
}

func connectCache() *redis.Client {
	// connect to redis cache
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := redisClient.Ping().Result()

	if err != nil {
		fmt.Println("Redis connection failed ❌")
		panic(err)
	}

	fmt.Println("Redis connection successful ✅")
	return redisClient
}

var db *gorm.DB
var redisClient *redis.Client

func main() {
	db = connectDatabase()
	redisClient = connectCache()
	http.HandleFunc("/products", httpHandler)
	fmt.Println("Server has started ✅")
	http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	response, err := getProducts(db, redisClient)

	if err != nil {
		fmt.Fprintf(w, err.Error()+"\r\n")
	} else {
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")

		if err := enc.Encode(response); err != nil {
			fmt.Println(err.Error())
		}
	}

}
