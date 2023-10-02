package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	IMG_SVR_URL = "image-db:6379"
	PASS        = "12345"
	USER        = "image"
)

type ImageDBClient struct {
	client *redis.Client
}

var ctx = context.Background()

func (this *ImageDBClient) NewClient() error {
	this.client = redis.NewClient(&redis.Options{
		Addr:     IMG_SVR_URL,
		Password: PASS,
		Username: USER,
	})
	pong, err := this.client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)

	fmt.Println("Connected with image db!")
	return nil
}
