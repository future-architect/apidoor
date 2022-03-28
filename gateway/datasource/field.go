package datasource

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/future-architect/apidoor/gateway/model"
	"strings"
)

var (
	defaultAPICallMaxLimit = 100
)

func CreateField(ctx context.Context, routing *Routing) (model.Field, error) {
	var schema string
	if strings.HasPrefix(routing.ForwardURL, "http://") {
		schema = "http"
		routing.ForwardURL = strings.Replace(routing.ForwardURL, "http://", "", 1)
		fmt.Println("Redis After:", routing.ForwardURL)
	} else if strings.HasPrefix(routing.ForwardURL, "https://") {
		schema = "https"
		routing.ForwardURL = strings.Replace(routing.ForwardURL, "https://", "", 1)
	} else {
		// スキーマが存在しない(tcpなどのスキーマは非対応)
		schema = "http"
	}
	path := model.NewURITemplate(routing.ForwardURL)
	template := model.NewURITemplate(routing.Path)

	count, err := logger.APICounter.GetCount(ctx, routing.ContractID, template)
	if err != nil {
		return model.Field{}, fmt.Errorf("fetch api count error: %w", err)
	}

	return model.Field{
		ContractID:    routing.ContractID,
		Template:      template,
		ForwardSchema: schema,
		Path:          path,
		Num:           count,
		Max:           defaultAPICallMaxLimit,
	}, nil

}
