package managementapi

import "net/http"

// PostUser godoc
// @Summary Create a user
// @Description Create a user
// @produce json
// @Param product body PostUserReq true "user description"
// @Success 201 {string} string
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Router /users [post]
func PostUser(w http.ResponseWriter, r *http.Request) {

}
