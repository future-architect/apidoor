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

func CreateField(ctx context.Context, key, hkey, forwardURL string) (model.Field, error) {
	var schema string
	if strings.HasPrefix(forwardURL, "http://") {
		schema = "http"
		forwardURL = strings.Replace(forwardURL, "http://", "", 1)
		fmt.Println("Redis After:", forwardURL)
	} else if strings.HasPrefix(forwardURL, "https://") {
		schema = "https"
		forwardURL = strings.Replace(forwardURL, "https://", "", 1)
	} else {
		// スキーマが存在しない(tcpなどのスキーマは非対応)
		schema = "http"
	}
	path := model.NewURITemplate(forwardURL)
	template := model.NewURITemplate(hkey)

	count, err := logger.APICounter.GetCount(ctx, key, template)
	if err != nil {
		return model.Field{}, fmt.Errorf("fetch api count error: %w", err)
	}

	return model.Field{
		Template:      template,
		ForwardSchema: schema,
		Path:          path,
		Num:           count,
		Max:           defaultAPICallMaxLimit,
	}, nil

}
