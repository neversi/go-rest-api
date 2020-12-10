package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
  	"gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"
)

// RedisCache ...
type RedisCache struct {
	ctx context.Context
	host string
	db int
	expireTime time.Duration
}

// NewRedisCache ...
func NewRedisCache(newHost string, newDB int, expireT time.Duration) *RedisCache {
	return &RedisCache{
		ctx: context.Background(),
		host: newHost,
		db: newDB,
		expireTime: expireT,
	}
}

func (rc *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: rc.host,
		Password: "",
		DB: rc.db,
	})
}

// Set ...
func (rc *RedisCache) Set(key string, value *models.Task) error {
	client := rc.getClient()

	_, err := json.Marshal(value)
	if err != nil {
		return err
	}

	client.Set(rc.ctx, key, value, rc.expireTime * time.Second)
	return nil
}

// Get ...
func (rc *RedisCache) Get(key string) (*models.Task,error) {
	client := rc.getClient()
	
	str, err := client.Get(rc.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	task := new(models.Task)

	err = json.Unmarshal([]byte(str), &task)
	if err != nil {
		return nil, err
	}
	return task, nil
}