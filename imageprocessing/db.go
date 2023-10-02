package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	IMG_SVR_URL = "image-db:6379"
	PASS        = "12345"
	USER        = "image"
)

type ImageDBClient struct {
	client *redis.Client
	ctx    context.Context
}

var ctx = context.Background()

func (this *ImageDBClient) NewClient() error {
	this.client = redis.NewClient(&redis.Options{
		Addr:     IMG_SVR_URL,
		Password: PASS,
		Username: USER,
	})
	this.ctx = context.Background()
	pong, err := this.client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println(pong, err)

	log.Println("Connected with image db!")
	return nil
}

func (this *ImageDBClient) WriteImage(username string, img Image) error {
	rawImg, err := json.Marshal(img)
	if err != nil {
		return err
	}
	stat, err := this.client.Set(this.ctx, username, rawImg, 0).Result()
	if err != nil {
		return err
	}

	log.Println(stat)

	return nil
}

func (this *ImageDBClient) ReadImage(username string) (*Image, error) {

	rawImg, err := this.client.Get(this.ctx, username).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("Key not found in Redis")
		}
		return nil, err
	}

	var data *Image
	if err := json.Unmarshal([]byte(rawImg), &data); err != nil {
		return nil, err
	}
	return data, nil
}
