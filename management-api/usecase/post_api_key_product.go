package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
	"sort"
)

func PostAPIKeyProducts(ctx context.Context, req *model.PostAPIKeyProductsReq) error {
	apiKeyID := *req.ApiKeyID
	keyAndUserID, err := db.fetchAPIKeyAndUser(ctx, apiKeyID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ClientError{fmt.Errorf("apikey not found, id %d", *req.ApiKeyID)}
		}
		log.Printf("fetching apikey user failed: %v", err)
		return ServerError{err}
	}

	contractProducts, err := db.fetchContractProductToAuth(ctx, keyAndUserID.userID, req.Contracts)
	if err != nil {
		log.Printf("fetching contract product to auth failed: %v", err)
		return ServerError{err}
	}

	// check whether all contracts are the user's
	if err := checkAllProductsFetched(req, contractProducts); err != nil {
		return ClientError{err}
	}

	// check whether some products are registered in multiple contracts
	if err := checkProductAddedDuplicately(contractProducts); err != nil {
		return ClientError{err}
	}

	// TODO: 中途失敗時のrollback処理
	err = db.postAPIKeyContractProductAuthorized(ctx, apiKeyID, contractProducts)
	if err != nil {
		log.Printf("insert apikey_contract_product_authorized db error: %v", err)
		if err, ok := err.(dbConstraintErr); ok {
			return ClientError{err}
		}
		return ServerError{err}
	}

	productIDs := productIDs(contractProducts)
	log.Printf("products %v", productIDs)
	swaggers, err := apirouting.ApiDBDriver.BatchGetSwagger(ctx, productIDs)
	if err != nil {
		log.Printf("get swagger info list db error: %v", err)
		return ServerError{err}
	}
	routings, err := generateRoutings(keyAndUserID.apiKey, contractProducts, swaggers)
	if err != nil {
		log.Printf("in generating routings, data consistency error: %v", err)
		return ServerError{err}
	}

	log.Println(routings)

	_, err = apirouting.ApiDBDriver.BatchPostRouting(ctx, routings)
	if err != nil {
		log.Printf("post api routing db error: %v", err)
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

func checkProductAddedDuplicately(contractProducts []model.ContractProductDB) error {
	appeared := make(map[int]int)
	for _, cp := range contractProducts {
		product := cp.ProductID
		contract, ok := appeared[product]
		if ok {
			return fmt.Errorf("product, id %d, are registered in multiple contracts, id = %d, %d",
				product, contract, cp.ContractID)
		}
		appeared[product] = cp.ContractID
	}
	return nil
}

// searchMissedIDs compares
// if products_ids is missed in the request, that means all products related to the contract will be added,
// searchMissedIds always returns nil, since want array is empty.
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

func productIDs(items []model.ContractProductDB) []int {
	ret := make([]int, len(items))
	for i, v := range items {
		ret[i] = v.ProductID
	}
	return ret
}

func generateRoutings(apikey string, products []model.ContractProductDB, swaggers []model.Swagger) ([]model.Routing, error) {
	swaggerMap := make(map[int]model.Swagger)
	for _, v := range swaggers {
		swaggerMap[v.ProductID] = v
	}
	routings := make([]model.Routing, 0, len(products)*20) // assume each product has 20 APIs on average
	for _, v := range products {
		swagger, ok := swaggerMap[v.ProductID]
		if !ok {
			return nil, fmt.Errorf("swagger info related to product, id %d, not found", v.ProductID)
		}
		for _, scheme := range swagger.Schemes {
			for _, api := range swagger.APIList {
				routings = append(routings, model.Routing{
					APIKey:     apikey,
					Path:       swagger.PathBase + api.Path,
					ForwardURL: fmt.Sprintf("%s://%s%s", scheme, swagger.ForwardURLBase, api.ForwardURL),
					ContractID: v.ContractID,
				})
			}
		}
	}
	return routings, nil
}
