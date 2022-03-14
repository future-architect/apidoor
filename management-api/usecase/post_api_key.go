package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

const APIKeyBytes = 16

func PostAPIKey(ctx context.Context, req model.PostAPIKeyReq) (*model.PostAPIKeyResp, error) {
	userID, err := fetchUserID(ctx, req.UserAccountID)
	if err != nil {
		log.Printf("fetch user id error: %v", err)
		if errors.Is(err, ErrNotFound) {
			return nil, ClientError{fmt.Errorf("account_id %s does not exist", req.UserAccountID)}
		} else {
			return nil, ServerError{err}
		}
	}
	key := generateKey(APIKeyBytes)

	apiKey := model.APIKey{
		UserID:    userID,
		AccessKey: key,
	}

	keyDescription, err := db.postAPIKey(ctx, apiKey)
	if err != nil {
		log.Printf("db insert api key error: %v", err)
		// this error occurs when the api_user is deleted after fetching its id
		if constraintErr, ok := err.(*dbConstraintErr); ok {
			var errMsg string
			switch constraintErr.field {
			case "user_id":
				errMsg = fmt.Sprintf("account_id %s does not exist", req.UserAccountID)
			}
			return nil, ClientError{errors.New(errMsg)}
		} else {
			return nil, ServerError{err}
		}
	}
	return &model.PostAPIKeyResp{
		UserAccountID: req.UserAccountID,
		AccessKey:     keyDescription.AccessKey,
		CreatedAt:     keyDescription.CreatedAt,
		UpdatedAt:     keyDescription.UpdatedAt,
	}, nil
}

func generateKey(byteLength int) string {
	key := make([]byte, byteLength)
	if _, err := rand.Read(key); err != nil {
		// unreachable because rand.Read always returns nil error
		return ""
	}
	return hex.EncodeToString(key)
}
