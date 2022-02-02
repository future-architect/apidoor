package managementapi

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/lib/pq"
	"log"
	"os"
	"text/template"

	"github.com/jmoiron/sqlx"
)

var db *sqlDB

var (
	//go:embed sql/search_api.sql
	searchAPISQLTemplateStr string
	searchAPISQLTemplate    *template.Template

	foreignKeyErrCode pq.ErrorCode   = "23503"
	foreignKeyErr     constraintType = "foreign key constraint"
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
		log.Fatalf("create searchAPISQL template %v", err)
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

func (sd sqlDB) getAPIInfo(ctx context.Context) ([]APIInfo, error) {
	rows, err := sd.driver.QueryxContext(ctx, "SELECT * from apiinfo")
	if err != nil {
		return nil, fmt.Errorf("sql execution error: %w", err)
	}

	var list []APIInfo
	for rows.Next() {
		var row APIInfo

		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("scanning record error: %w", err)
		}

		list = append(list, row)
	}

	return list, nil
}

func (sd sqlDB) postAPIInfo(ctx context.Context, info *PostAPIInfoReq) error {
	_, err := sd.driver.NamedExecContext(ctx,
		"INSERT INTO apiinfo(name, source, description, thumbnail, swagger_url) VALUES(:name, :source, :description, :thumbnail, :swagger_url)",
		info)
	if err != nil {
		return fmt.Errorf("sql execution error: %w", err)
	}
	return nil
}

func (sd sqlDB) searchAPIInfo(ctx context.Context, params *SearchAPIInfoParams) (*SearchAPIInfoResp, error) {
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

	list := make([]APIInfo, 0)
	count := 0
	for rows.Next() {
		var row SearchAPIInfoResult
		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("scanning record error: %w", err)
		}

		list = append(list, row.APIInfo)
		count = row.Count
	}
	metaData := SearchAPIInfoMetaData{
		ResultSet: ResultSet{
			Count:  count,
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	}

	return &SearchAPIInfoResp{
		APIList:               list,
		SearchAPIInfoMetaData: metaData,
	}, nil
}

func (sd sqlDB) postUser(ctx context.Context, user *PostUserReq) error {
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

func (sd sqlDB) postProduct(ctx context.Context, product *PostProductReq) error {
	tx, err := sd.driver.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	stmt, err := tx.PrepareNamedContext(ctx,
		`INSERT INTO product(name, display_name, source, description, thumbnail, is_available, created_at, updated_at)
				VALUES (:name, :display_name, :source, :description, :thumbnail, :is_available, current_timestamp, current_timestamp)
				RETURNING id`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare sql to insert product failed: %w", err)
	}

	var productID int
	err = stmt.QueryRowxContext(ctx, product).Scan(&productID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("execute sql to insert product failed: %w", err)
	}

	for _, content := range product.Contents {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO product_api_content(product_id, api_id, description, created_at, updated_at)
					VALUES ($1, $2, $3, current_timestamp, current_timestamp)`,
			productID, content.ID, content.Description)
		if err != nil {
			tx.Rollback()
			if postgresErr, ok := err.(*pq.Error); ok {
				if postgresErr.Code == foreignKeyErrCode {
					return &dbConstraintErr{
						constraintType: foreignKeyErr,
						field:          "api_id",
						value:          content.ID,
						message:        fmt.Sprintf("insert content, api_id = %d, failed: foreign key constraint", content.ID),
					}
				}
			}
			return fmt.Errorf("insert content, api_id = %d, failed: %w", content.ID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction faield: %w", err)
	}
	return nil
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
