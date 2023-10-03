package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	USER = "recommend"
)

var (
	RECOMMEND_SVR_URL = "recommend-db:6379"
	PASS              = "12345"
)

type RecommendDBClient struct {
	client *redis.Client
	ctx    context.Context
	mx     sync.RWMutex
}

var ctx = context.Background()

func (this *RecommendDBClient) NewClient() error {
	this.client = redis.NewClient(&redis.Options{
		Addr:     RECOMMEND_SVR_URL,
		Password: PASS,
		Username: USER,
	})
	this.ctx = context.Background()
	pong, err := this.client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println(pong, err)

	log.Println("Connected with recommend db!")
	return nil
}

func (this *RecommendDBClient) WriteRecommendations(username string, rec Recommendations) error {
	this.mx.Lock()
	defer this.mx.Unlock()

	rawImg, err := json.Marshal(rec)
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

func (this *RecommendDBClient) ReadRecommendations(username string) (*Recommendations, error) {
	this.mx.RLock()
	defer this.mx.RUnlock()

	rawImg, err := this.client.Get(this.ctx, username).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, err
	}

	var data *Recommendations
	if err := json.Unmarshal([]byte(rawImg), &data); err != nil {
		return nil, err
	}
	return data, nil
}
