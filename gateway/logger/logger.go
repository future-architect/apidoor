package logger

import (
	"context"
	"encoding/csv"
	"github.com/Songmu/flextime"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type BillingStatus int

const (
	NotBilling BillingStatus = iota
	Billing
)

var billingStatuses = [...]string{"not billing", "billing"}

func (bs BillingStatus) String() string {
	return billingStatuses[bs]
}

type Appender interface {
	UpdateDB(ctx context.Context)
	Do(key, path string, r *http.Request,
		apiResp *http.Response, calcBillingStatus func(resp *http.Response) BillingStatus) error
}

func UpdateDBRoutine(ctx context.Context, appender Appender, interval time.Duration, kill, finish chan bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			appender.UpdateDB(ctx)
		case <-kill:
			appender.UpdateDB(ctx)
			finish <- true
			return
		}
	}
}

func CleanupUpdateDBTask(kill, finish chan bool) {
	kill <- true
	<-finish
}

type DefaultAppender struct {
	Writer   io.Writer
	LogItems LogItems
}

func (a *DefaultAppender) Do(key, path string, r *http.Request,
	apiResp *http.Response, calcBillingStatus func(resp *http.Response) BillingStatus) error {
	logItem, err := NewLogItem(key, path, apiResp, calcBillingStatus)
	record := make([]string, 0, len(LogOptionPattern))
	for _, logOption := range LogOptionPattern {
		logOption(&record, &logItem, r)
	}
	if err != nil {
		return err
	}
	a.LogItems.Append(logItem)

	_, err = a.Writer.Write([]byte(strings.Join(record, ",") + "\n"))
	return err
}

func (a *DefaultAppender) UpdateDB(ctx context.Context) {
	logItems := a.LogItems.ReadAndDeleteAll()
	for _, item := range logItems {
		if err := db.postAccessLogDB(ctx, item); err != nil {
			log.Printf("putting log info, %v, failed: %v", item, err)
		}
	}
}

type CSVAppender struct {
	Writer *csv.Writer

	LogItems LogItems
}

func NewCSVAppender(writer *csv.Writer) CSVAppender {
	return CSVAppender{
		Writer:   writer,
		LogItems: NewLogItems(),
	}
}

func (a *CSVAppender) Do(key, path string, r *http.Request,
	apiResp *http.Response, calcBillingStatus func(resp *http.Response) BillingStatus) error {
	logItem, err := NewLogItem(key, path, apiResp, calcBillingStatus)
	record := make([]string, 0, len(LogOptionPattern))
	for _, logOption := range LogOptionPattern {
		logOption(&record, &logItem, r)
	}
	if err != nil {
		return err
	}
	a.LogItems.Append(logItem)
	return a.Writer.Write(record)
}

func (a *CSVAppender) UpdateDB(ctx context.Context) {
	logItems := a.LogItems.ReadAndDeleteAll()
	for _, item := range logItems {
		if err := db.postAccessLogDB(ctx, item); err != nil {
			log.Printf("putting log info, %v, failed: %v", item, err)
		}
	}
}

type LogItem struct {
	TimeStamp     string        `dynamo:"timestamp"`
	Key           string        `dynamo:"api_key"`
	Path          string        `dynamo:"path"`
	StatusCode    int           `dynamo:"status_code"`
	BillingStatus BillingStatus `dynamo:"billing_status"`
}

func NewLogItem(key, path string, apiResp *http.Response, calcBillingStatus func(resp *http.Response) BillingStatus) (LogItem, error) {
	return LogItem{
		TimeStamp:     flextime.Now().Format(time.RFC3339),
		Key:           key,
		Path:          path,
		StatusCode:    apiResp.StatusCode,
		BillingStatus: calcBillingStatus(apiResp),
	}, nil
}

type LogItems struct {
	sync.Mutex
	Items []LogItem
}

func NewLogItems() LogItems {
	return LogItems{
		Items: make([]LogItem, 0),
	}
}

func (li *LogItems) Append(item LogItem) {
	li.Lock()
	defer li.Unlock()
	li.Items = append(li.Items, item)
}

func (li *LogItems) ReadAndDeleteAll() []LogItem {
	li.Lock()
	defer li.Unlock()
	items := make([]LogItem, len(li.Items))
	for i, v := range li.Items {
		items[i] = v
	}
	li.Items = li.Items[:0]
	return items
}
