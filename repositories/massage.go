package repositories

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Massage repository
type Massage struct {
	client *redis.Client
}

// NewMassageRepo create massage repository
func NewMassageRepo(c *redis.Client) *Massage {
	return &Massage{client: c}
}

// Add massage to redis
func (r *Massage) Add(setName string, timestamp float64, msg string) error {
	if err := r.client.ZAdd(setName, redis.Z{
		Score:  timestamp,
		Member: msg,
	}).Err(); err != nil {
		return err
	}

	return nil
}

// Fetch massage from redis
func (r *Massage) Fetch(setName string) (redis.ZWithKey, error) {
	val, err := r.client.BZPopMin(3*time.Second, setName).Result()
	if err != nil {
		return val, err
	}
	return val, nil
}

// FindByTime from redis
func (r *Massage) FindByTime(setName string, timestamp float64) ([]redis.Z, error) {
	val, err := r.client.ZRangeByScoreWithScores(setName, redis.ZRangeBy{
		Min: string(0),
		Max: fmt.Sprintf("%f", timestamp),
	}).Result()
	if err != nil {
		return nil, err
	}
	return val, err
}

// Remove from redis
func (r *Massage) Remove(setName string, jobName string) (int64, error) {
	val, err := r.client.ZRem(setName, jobName).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}
