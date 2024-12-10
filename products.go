package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Products struct {
	ProductId   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	RetailPrice float64 `json:"source"`
}

type JsonResponse struct {
	Data   []Products `json:"data"`
	Source string     `json:"source"`
}

func getProducts(db *gorm.DB, redisClient *redis.Client) (*JsonResponse, error) {

	cachedProducts, err := redisClient.Get("products").Bytes()
	response := JsonResponse{}
	var products []Products

	// this error to check whether cachedProducts exists is REDIS CACHE
	if err != nil {

		// if doesnt exist in redis, fetch from DB
		productsDBresult := db.Find(&products)

		if productsDBresult.Error != nil {
			return nil, err
		}

		// store the DB-fetched-products
		productsBytes, err := json.Marshal(products)

		// check for error while unmarshalling (Products[] -> json encoding)
		if err != nil {
			return nil, err
		}

		// cache the DB-fetched-products in REDIS
		err = redisClient.Set("products", productsBytes, 10*time.Second).Err()

		// check for error while storing into redis cache
		if err != nil {
			return nil, err
		}

		response = JsonResponse{Data: products, Source: "PostgreSQL"}
		return &response, err
	}

	err = json.Unmarshal(cachedProducts, &products)

	// check for error while unmarshalling (json encoding -> Products[])
	if err != nil {
		return nil, err
	}

	response = JsonResponse{Data: products, Source: "Redis Cache"}
	return &response, nil
}

// func fetchFromDb(db *gorm.DB) ([]Products, error) {

// 	// dbUser := ""
// 	// dbPassword := ""
// 	dbName := "sample_company"

// 	conString := fmt.Sprintf("host=localhost dbname=%s sslmode=disable", dbName)

// 	db, err := sql.Open("postgres", conString)

// 	if err != nil {
// 		return nil, err
// 	}

// 	queryString := `select product_id, product_name, retail_price from products`

// 	rows, err := db.Query(queryString)

// 	if err != nil {
// 		return nil, err
// 	}

// 	var records []Products

// 	for rows.Next() {
// 		var p Products
// 		err = rows.Scan(&p.ProductId, &p.ProductName, &p.RetailPrice)

// 		records = append(records, p)

// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return records, nil
// }
