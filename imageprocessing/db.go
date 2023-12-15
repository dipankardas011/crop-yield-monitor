package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	USER = "image"
)

var (
	IMG_SVR_URL = ""
	PASS        = ""
)

type ImageDBClient struct {
	client *redis.Client
	ctx    context.Context
	mx     sync.RWMutex
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
	this.mx.Lock()
	defer this.mx.Unlock()

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
	this.mx.RLock()
	defer this.mx.RUnlock()

	rawImg, err := this.client.Get(this.ctx, username).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("Key not found in Redis")
		}
		return nil, err
	}

	var data *Image
	if err := json.Unmarshal(rawImg, &data); err != nil {
		return nil, err
	}
	return data, nil
}
