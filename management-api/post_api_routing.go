package managementapi

import (
	"net/http"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

// PostAPI godoc
// @Summary Post API routing
// @Description Post a new API routing
// @produce json
// @Success 201 {object} Routing
// @Router /api [post]
func PostAPIRouting(w http.ResponseWriter, r *http.Request) {

}
