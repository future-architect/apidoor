package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
	"sort"
)

func PostAPIKeyProducts(ctx context.Context, req *model.PostAPIKeyProductsReq) error {
	apiKeyID := *req.ApiKeyID
	userID, err := db.fetchAPIKeyUser(ctx, apiKeyID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ClientError{fmt.Errorf("apikey not found, id %d", *req.ApiKeyID)}
		}
		log.Printf("fetching apikey user failed: %v", err)
		return ServerError{err}
	}

	contractProducts, err := db.fetchContractProductToAuth(ctx, userID, req.Contracts)
	if err != nil {
		log.Printf("fetching contract product to auth failed: %v", err)
		return ServerError{err}
	}

	// check whether all contracts are the user's
	if err := checkAllProductsFetched(req, contractProducts); err != nil {
		return ClientError{err}
	}

	// TODO: gatewayのルーティングの追加

	if err = db.postAPIKeyContractProductAuthorized(ctx, apiKeyID, contractProducts); err != nil {
		log.Printf("inser apikey_contract_product_authorized db error: %v", err)
		if err, ok := err.(dbConstraintErr); ok {
			return ClientError{err}
		}
		return ServerError{err}
	}

	return nil
}

func checkAllProductsFetched(req *model.PostAPIKeyProductsReq, contractProducts []model.ContractProductDB) error {
	reqContractIDMap := req.ContractIDMap()
	gotContractIDs := make(map[int][]int)
	for _, cp := range contractProducts {
		gotContractIDs[cp.ContractID] = append(gotContractIDs[cp.ContractID], cp.ProductID)
	}

	for contractID, item := range reqContractIDMap {
		gotProducts, ok := gotContractIDs[contractID]
		if !ok {
			return fmt.Errorf("contract %d does not exist or is not yours, or some products in contract %d are wrong", contractID, contractID)
		}
		wantProducts := item.ProductIDs

		if wantProducts != nil {
			missedProductIDs := searchMissedIDs(wantProducts, gotProducts)
			if missedProductIDs != nil {
				return fmt.Errorf("following product ids is not found in ids linked to contract %d: %v", contractID, missedProductIDs)
			}
		}
	}
	return nil
}

func searchMissedIDs(want []int, got []int) []int {
	sort.Ints(want)
	sort.Ints(got)

	var ret []int
	idx := 0
	for _, val := range want {
		if idx == len(got) || got[idx] != val {
			ret = append(ret, val)
			continue
		}
		idx++
	}
	return ret
}
