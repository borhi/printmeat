package main

import (
	"log"
	"net/http"
	"printMeAt/repositories"
	"printMeAt/services"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:       "redis:6379",
		PoolSize:   5,
		MaxRetries: 2,
		DB:         0,
	})

	repo := repositories.NewMassageRepo(redisClient)
	service := services.NewPrintService(repo)
	go service.FeedBack()
	go service.Run()

	r := mux.NewRouter()
	r.HandleFunc("/printMeAt", func(w http.ResponseWriter, r *http.Request) {
		strTime := r.FormValue("time")
		msg := r.FormValue("massage")

		if r.FormValue("time") == "" {
			http.Error(w, "time parameter missed", http.StatusBadRequest)
			return
		}

		if r.FormValue("massage") == "" {
			http.Error(w, "massage parameter missed", http.StatusBadRequest)
			return
		}

		flTime, err := strconv.ParseFloat(strTime, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if int64(flTime) < time.Now().Unix() {
			http.Error(w, "time must be in the future", http.StatusBadRequest)
			return
		}

		if err := service.Schedule(flTime, msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
