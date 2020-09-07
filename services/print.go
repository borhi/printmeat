package services

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	setName             = "print"
	backupSetName       = "print_backup"
	maxTimeWindow int64 = 10
)

// Print service
type Print struct {
	repo massageRepository
}

// NewPrintService create print service
func NewPrintService(repo massageRepository) *Print {
	return &Print{repo: repo}
}

type massageRepository interface {
	Fetch(setName string) (redis.ZWithKey, error)
	Add(setName string, timestamp float64, msg string) error
	Remove(setName string, jobName string) (int64, error)
	FindByTime(setName string, timestamp float64) ([]redis.Z, error)
}

// Run print service
func (s *Print) Run() error {
	fmt.Println("[INFO] Starting prcessor")
	for {
		msg, _ := s.repo.Fetch(setName)

		if msg.Member == nil {
			continue
		}

		if int64(msg.Score)-time.Now().Unix() > maxTimeWindow {
			s.reSchedule(msg.Score, fmt.Sprintf("%v", msg.Member))
			time.Sleep(1 * time.Second)
		} else {
			for {
				if int64(msg.Score)-time.Now().Unix() <= 0 {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}

			fmt.Println("Running job :: Time:", msg.Score, " Job:", msg.Member)
			s.repo.Remove(backupSetName, fmt.Sprintf("%v", msg.Member))
		}
	}
}

// Schedule massage to sets
func (s *Print) Schedule(timestamp float64, msg string) error {
	if err := s.repo.Add(setName, timestamp, msg); err != nil {
		return err
	}

	if err := s.repo.Add(backupSetName, timestamp, msg); err != nil {
		return err
	}

	return nil
}

// FeedBack old massages
func (s *Print) FeedBack() {
	for {
		msgs, _ := s.repo.FindByTime(backupSetName, float64(time.Now().Unix()-maxTimeWindow))
		for _, msgWithScore := range msgs {
			s.reSchedule(msgWithScore.Score, fmt.Sprintf("%v", msgWithScore.Member))
		}

		time.Sleep(5 * time.Second)
	}
}

func (s *Print) reSchedule(timestamp float64, msg string) error {
	if err := s.repo.Add(setName, timestamp, msg); err != nil {
		return err
	}

	return nil
}
