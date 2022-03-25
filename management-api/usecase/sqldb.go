package usecase

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/lib/pq"
	"log"
	"os"
	"text/template"

	"github.com/jmoiron/sqlx"
)

var db *sqlDB

var (
	//go:embed sql/search_product.sql
	searchAPISQLTemplateStr string
	searchAPISQLTemplate    *template.Template

	//go:embed sql/fetch_products.sql
	fetchProductsSQLTemplateStr string
	fetchProductsSQLTemplate    *template.Template

	//go:embed sql/fetch_products_linked_to_contracts.sql
	fetchProductsLinkedToContractsSQLTemplateStr string
	fetchProductsLinkedToContractsSQLTemplate    *template.Template

	foreignKeyErrCode pq.ErrorCode   = "23503"
	foreignKeyErr     constraintType = "foreign key constraint"

	ErrNotFound = errors.New("db: item not found")
)

func init() {
	// setup DB driver
	var err error
	if db, err = NewSqlDB(); err != nil {
		log.Fatalf("setup postgreSQL failed: %v", err)
	}

	// setup sql template
	searchAPISQLTemplate, err = template.New("search API SQL template").Parse(searchAPISQLTemplateStr)
	if err != nil {
		log.Fatalf("creating searchAPISQL template failed: %v", err)
	}
	fetchProductsSQLTemplate, err = template.New("fetch API products by product names SQL template").Parse(fetchProductsSQLTemplateStr)
	if err != nil {
		log.Fatalf("creating fetchProductsSQL template failed: %v", err)
	}

	fetchProductsLinkedToContractsSQLTemplate, err = template.New("fetch contract_products").Parse(fetchProductsLinkedToContractsSQLTemplateStr)
	if err != nil {
		log.Fatalf("creating fetchProductsLinkedToContractsSQL template failed: %v", err)
	}

}

type sqlDB struct {
	driver *sqlx.DB
}

func NewSqlDB() (*sqlDB, error) {
	dbDriver := os.Getenv("DATABASE_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSLMODE"))

	db, err := sqlx.Open(dbDriver, dbSource)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	return &sqlDB{
		driver: db,
	}, nil
}

func (sd sqlDB) getProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := sd.driver.QueryxContext(ctx, "SELECT * from product")
	if err != nil {
		return nil, fmt.Errorf("sql execution error: %w", err)
	}

	var list []model.Product
	for rows.Next() {
		var row model.Product

		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("scanning record error: %w", err)
		}

		list = append(list, row)
	}

	return list, nil
}

func (sd sqlDB) postProduct(ctx context.Context, product *model.PostProductDB) (*model.Product, error) {
	ret := new(model.Product)
	stmt, err := sd.driver.PrepareNamedContext(ctx,
		`INSERT INTO product(name, source, display_name, description, thumbnail, base_path, swagger_url, is_available, created_at, updated_at)
			VALUES(:name, :source, :display_name, :description, :thumbnail, :base_path, :swagger_url, :is_available, current_timestamp, current_timestamp) RETURNING *`)
	err = stmt.QueryRowxContext(ctx, product).StructScan(ret)
	if err != nil {
		return nil, fmt.Errorf("sql execution error: %w", err)
	}
	return ret, nil
}

func (sd sqlDB) deleteProduct(ctx context.Context, productID int) error {
	_, err := sd.driver.ExecContext(ctx, `DELETE FROM product WHERE id = $1`, productID)
	if err != nil {
		return fmt.Errorf("sql execution error: %w", err)
	}
	return nil
}

func (sd sqlDB) searchProduct(ctx context.Context, params *model.SearchProductParams) (*model.SearchProductResp, error) {
	var query bytes.Buffer
	if err := searchAPISQLTemplate.Execute(&query, params); err != nil {
		return nil, fmt.Errorf("generate SQL error: %w", err)
	}
	targetValues := make(map[string]interface{}, len(params.Q)+2)
	for i, q := range params.Q {
		key := fmt.Sprintf("q%d", i)
		targetValues[key] = q
	}
	targetValues["limit"] = params.Limit
	targetValues["offset"] = params.Offset

	rows, err := sd.driver.NamedQueryContext(ctx, query.String(), targetValues)

	if err != nil {
		return nil, fmt.Errorf("sql execution error: %w", err)
	}

	list := make([]model.Product, 0)
	count := 0
	for rows.Next() {
		var row model.SearchProductResult
		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("scanning record error: %w", err)
		}

		list = append(list, row.Product)
		count = row.Count
	}

	return &model.SearchProductResp{
		ProductList: list,
		SearchProductMetaData: model.SearchProductMetaData{
			ResultSet: model.ResultSet{
				Count:  count,
				Limit:  params.Limit,
				Offset: params.Offset,
			},
		},
	}, nil
}

func (sd sqlDB) postUser(ctx context.Context, user *model.PostUserReq) error {
	_, err := sd.driver.NamedExecContext(ctx,
		`INSERT INTO apiuser(account_id, email_address, login_password_hash, name, created_at, updated_at)
				VALUES(:account_id, :email_address, crypt(:password, gen_salt('bf')),
			    :name, current_timestamp, current_timestamp)`,
		user)
	if err != nil {
		return fmt.Errorf("sql execution error: %w", err)
	}
	return nil
}

func (sd sqlDB) fetchUser(ctx context.Context, accountID string) (*model.User, error) {
	rows, err := sd.driver.QueryxContext(ctx, "SELECT * FROM apiuser WHERE account_id = $1", accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	cnt := 0
	var user model.User
	for rows.Next() {
		if cnt > 0 {
			// unreachable because apiuser.account_id has a unique constraint
			return nil, fmt.Errorf("multiple users have an account_id %s", accountID)
		}

		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("failed to scan result as user: %w", err)
		}
		cnt++
	}

	if cnt == 0 {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (sd sqlDB) fetchProduct(ctx context.Context, productName string) (*model.Product, error) {
	rows, err := sd.driver.QueryxContext(ctx, "SELECT * FROM product WHERE name = $1", productName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	cnt := 0
	var product model.Product
	for rows.Next() {
		if cnt > 0 {
			// unreachable because product.name has a unique constraint
			return nil, fmt.Errorf("multiple products have a name %s", productName)
		}

		if err := rows.StructScan(&product); err != nil {
			return nil, fmt.Errorf("failed to scan result as product: %w", err)
		}
		cnt++
	}

	if cnt == 0 {
		return nil, ErrNotFound
	}
	return &product, nil
}

func (sd sqlDB) fetchProducts(ctx context.Context, productNames []string) (map[string]model.Product, error) {
	var query bytes.Buffer
	if err := fetchProductsSQLTemplate.Execute(&query, productNames); err != nil {
		return nil, fmt.Errorf("generate SQL error: %w", err)
	}
	parameters := make(map[string]interface{})
	for i, name := range productNames {
		parameters[fmt.Sprintf("name%d", i)] = name
	}

	rows, err := sd.driver.NamedQueryContext(ctx, query.String(), parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	products := make(map[string]model.Product)
	var product model.Product
	for rows.Next() {
		if err := rows.StructScan(&product); err != nil {
			return nil, fmt.Errorf("failed to scan result as product: %w", err)
		}
		products[product.Name] = product
	}

	return products, nil
}

func (sd sqlDB) postContract(ctx context.Context, contract *model.PostContractDB) error {
	tx, err := sd.driver.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	stmt, err := tx.PreparexContext(ctx,
		`INSERT INTO contract(user_id, created_at, updated_at)
				VALUES ($1, current_timestamp, current_timestamp) RETURNING id`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare sql to insert contract failed: %w", err)
	}

	var contractID int
	err = stmt.QueryRowxContext(ctx, contract.UserID).Scan(&contractID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("execute sql to insert contract failed: %w", err)
	}

	for _, product := range contract.Products {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO contract_product_content(contract_id, product_id, description, created_at, updated_at)
					VALUES ($1, $2, $3, current_timestamp, current_timestamp)`,
			contractID, product.ProductID, product.Description)
		if err != nil {
			tx.Rollback()
			if postgresErr, ok := err.(*pq.Error); ok {
				if postgresErr.Code == foreignKeyErrCode {
					return &dbConstraintErr{
						constraintType: foreignKeyErr,
						field:          "product_id",
						value:          product.ProductID,
						message:        fmt.Sprintf("insert content, product_id = %d, failed: foreign key constraint", product.ProductID),
					}
				}
			}
			return fmt.Errorf("insert content, api_id = %d, failed: %w", product.ProductID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction faield: %w", err)
	}
	return nil
}

func (sd sqlDB) postAPIKey(ctx context.Context, apiKey model.APIKey) (*model.APIKey, error) {
	ret := new(model.APIKey)
	stmt, err := sd.driver.PrepareNamedContext(ctx,
		`INSERT INTO apikey(user_id, access_key, created_at, updated_at)
				VALUES (:user_id, :access_key, current_timestamp, current_timestamp) returning *`)
	if err != nil {
		return nil, fmt.Errorf("preparing sql query failed: %w", err)
	}

	err = stmt.QueryRowx(apiKey).StructScan(ret)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == foreignKeyErrCode {
				return nil, &dbConstraintErr{
					constraintType: foreignKeyErr,
					field:          postgresErr.Column,
					message:        "insert content failed: foreign key constraint",
				}
			}
		}
		return nil, fmt.Errorf("execute sql to insert api key failed: %w", err)
	}
	return ret, nil
}

func (sd sqlDB) fetchAPIKeyAndUser(ctx context.Context, apiKeyId int) (int, string, error) {
	var userId int
	var key string
	rows, err := sd.driver.QueryxContext(ctx,
		` SELECT user_id, access_key FROM apikey WHERE id = $1`, apiKeyId)
	if err != nil {
		return 0, "", fmt.Errorf("executing sql query failed: %w", err)
	}

	cnt := 0
	for rows.Next() {
		if err = rows.Scan(&userId, &key); err != nil {
			return 0, "", fmt.Errorf("scanning result into api_key failed: %v", err)
		}
		cnt++
	}
	if cnt == 0 {
		return 0, "", ErrNotFound
	}
	return userId, key, nil
}

func (sd sqlDB) fetchContractProductToAuth(ctx context.Context, userID int, contractProducts []model.AuthorizedContractProducts) ([]model.ContractProductDB, error) {
	var query bytes.Buffer
	if err := fetchProductsLinkedToContractsSQLTemplate.Execute(&query, contractProducts); err != nil {
		return nil, fmt.Errorf("generate SQL error: %w", err)
	}
	parameters := make(map[string]interface{})
	parameters["user_id"] = userID
	for i, contract := range contractProducts {
		parameters[fmt.Sprintf("contract_id_%d", i)] = contract.ContractID
		for j, productID := range contract.ProductIDs {
			parameters[fmt.Sprintf("product_id_%d_%d", i, j)] = productID
		}
	}

	rows, err := sd.driver.NamedQueryContext(ctx, query.String(), parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	products := make([]model.ContractProductDB, 0)
	var product model.ContractProductDB
	for rows.Next() {
		if err := rows.StructScan(&product); err != nil {
			return nil, fmt.Errorf("failed to scan result as contract_product: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (sd sqlDB) postAPIKeyContractProductAuthorized(ctx context.Context, apiKeyID int, contractProducts []model.ContractProductDB) ([]int, error) {
	ids := make([]int, len(contractProducts))

	tx, err := sd.driver.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}
	stmt, err := tx.PreparexContext(ctx,
		`INSERT INTO apikey_contract_product_authorized(apikey_id, contract_product_id, created_at, updated_at)
					VALUES ($1, $2, current_timestamp, current_timestamp) RETURNING id`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("prepare statement failed: %w", err)
	}
	for i, cp := range contractProducts {
		var id int
		err = stmt.QueryRowxContext(ctx, apiKeyID, cp.ID).Scan(&id)
		if err != nil {
			tx.Rollback()
			if postgresErr, ok := err.(*pq.Error); ok {
				if postgresErr.Code == foreignKeyErrCode {
					switch postgresErr.Column {
					case "apikey_id":
						return nil, &dbConstraintErr{
							constraintType: foreignKeyErr,
							field:          "apikey_id",
							value:          apiKeyID,
							message:        fmt.Sprintf("insert item, apikey_id = %d, failed: foreign key constraint", apiKeyID),
						}
					case "contract_product_id":
						return nil, &dbConstraintErr{
							constraintType: foreignKeyErr,
							field:          "contract_product_id",
							value:          cp.ID,
							message:        fmt.Sprintf("insert item, contract_product_id = %d, failed: foreign key constraint", cp.ID),
						}
					}
				}
			}
			return nil, fmt.Errorf("insert item, contract_product_id = %d, failed: %w", cp.ID, err)
		}
		ids[i] = id
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction faield: %w", err)
	}

	return ids, nil
}

type constraintType string

type dbConstraintErr struct {
	constraintType constraintType
	field          string
	value          interface{}
	message        string
}

func (dc dbConstraintErr) Error() string {
	return dc.message
}
