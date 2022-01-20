package main

import (
	"context"
	"encoding/csv"
	"github.com/future-architect/apidoor/gateway"
	"github.com/future-architect/apidoor/gateway/datasource/redis"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	updateDBInterval = 10 * time.Second
)

// gateway entry point @localhost
func main() {

	logPath := os.Getenv("LOG_PATH")
	if len(logPath) == 0 {
		logPath = "./log.csv"
	}

	// open log file
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	// write to log file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	h := gateway.DefaultHandler{
		Appender: &logger.CSVAppender{
			Writer: writer,
		},
		DataSource: redis.New(),
	}

	ctx := context.Background()

	// update access logging db
	routineKill := make(chan bool)
	routineFinish := make(chan bool)
	go logger.UpdateDBRoutine(ctx, h.Appender, updateDBInterval, routineKill, routineFinish)
	defer logger.CleanupUpdateDBTask(routineKill, routineFinish)

	// capturing keyboard interrupt
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("keyboard interrupt occurs")
		logger.CleanupUpdateDBTask(routineKill, routineFinish)
		os.Exit(2)
	}()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", h.Handle)
		r.Put("/*", h.Handle)
		r.Delete("/*", h.Handle)
		r.Post("/*", h.Handle)
	})

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
