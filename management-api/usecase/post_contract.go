package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostContract(ctx context.Context, req model.PostContractReq) error {
	userID, err := fetchUserID(ctx, req.UserAccountID)
	if err != nil {
		log.Printf("fetch user id error: %v", err)
		if errors.Is(err, ErrNotFound) {
			return ClientError{fmt.Errorf("account_id %s does not exist", req.UserAccountID)}
		}
		return ServerError{err}
	}

	products, err := fetchProductIDs(ctx, req.Products)
	if err != nil {
		log.Printf("fetch ids of products error: %v", err)
		return err
	}

	contract := model.PostContractDB{
		UserID:   userID,
		Products: products,
	}

	if err := db.postContract(ctx, &contract); err != nil {
		log.Printf("db insert contract error: %v", err)
		// this error occurs the api_user or the product is deleted after fetching its id
		if constraintErr, ok := err.(*dbConstraintErr); ok {
			return ClientError{fmt.Errorf("product_name %s does not exist", constraintErr.value)}
		}
		return ServerError{err}
	}
	return nil
}

// TODO: fetch*関数の位置
// おそらくusecase/のファイルをオブジェクトの種類ごとにして、user.goや、product.goあたりを見るのが良さそう

func fetchUserID(ctx context.Context, accountID string) (int, error) {
	// TODO　契約する権限の確認 (管理者ユーザと被管理ユーザのテーブルを作って用いる?)
	user, err := db.fetchUser(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("fetch user id db error: %w", err)
	}

	return user.ID, nil
}

func fetchProductIDs(ctx context.Context, products []*model.ContractProducts) ([]*model.ContractProductContentDB, error) {
	productNames := make([]string, len(products))
	for i, product := range products {
		productNames[i] = product.ProductName
	}

	// TODO: is_availableがtrueなproductからのみ取得する
	productMap, err := db.fetchProducts(ctx, productNames)
	if err != nil {
		return nil, ServerError{err}
	}

	contractProducts := make([]*model.ContractProductContentDB, len(products))
	for i, product := range products {
		matchedProduct, ok := productMap[product.ProductName]
		if !ok {
			return nil, ClientError{fmt.Errorf("product_name %s does not exist", product.ProductName)}
		}
		contractProducts[i] = &model.ContractProductContentDB{
			ProductID:   matchedProduct.ID,
			Description: product.Description,
		}
	}

	return contractProducts, nil
}
