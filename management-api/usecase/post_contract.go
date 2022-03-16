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

	products := make([]*model.ContractProductsDB, len(req.Products))
	for i, product := range req.Products {
		productID, err := fetchProductID(ctx, product.ProductName)
		if err != nil {
			log.Printf("fetch product id error: %v", err)
			if errors.Is(err, ErrNotFound) {
				return ClientError{fmt.Errorf("product_name %s does not exist", product.ProductName)}
			}
			return ServerError{err}
		}
		products[i] = &model.ContractProductsDB{
			ProductID:   productID,
			Description: product.Description,
		}
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

func fetchProductID(ctx context.Context, productName string) (int, error) {
	product, err := db.fetchProduct(ctx, productName)
	if err != nil {
		return 0, fmt.Errorf("fetch product id db error: %w", err)
	}

	return product.ID, nil
}
