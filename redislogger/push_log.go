package redislogger

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

func PushLog() {
	file, err := os.OpenFile("./log/log.csv", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	count := make(map[string]map[string]int)
	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		key, path := line[1], line[2]
		if _, ok := count[key]; ok {
			count[key][path]++
		} else {
			count[key] = make(map[string]int)
			count[key][path] = 1
		}
	}

	ctx := context.Background()
	for k, v := range count {
		for p, n := range v {
			if rdb.HExists(ctx, k, p).Val() {
				now, err := strconv.Atoi(rdb.HGet(ctx, k, p).Val())
				if err != nil {
					log.Fatal(err)
				}
				rdb.HSet(ctx, k, p, now+n)
			} else {
				rdb.HSet(ctx, k, p, n)
			}
		}
	}
}
