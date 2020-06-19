package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kacejot/resize-service/pkg/utils"

	"github.com/kacejot/resize-service/pkg/api/graph/model"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/kacejot/resize-service/pkg/storage/cloud"
)

// ArangoConfig stores necessary info for establishing connection
// with database
type ArangoConfig struct {
	Endpoint string
	User     string
	Password string
}

// RecordStorage is an abstraction for database storage
// It encapsilates connection with the database
type RecordStorage struct {
	records driver.Collection
}

// Record stores resize result after upload in database
type Record struct {
	Images cloud.UploadResult
	User   string
}

// LoadArangoConfig fills ArangoConfig with cmd args
func LoadArangoConfig() *ArangoConfig {
	return &ArangoConfig{
		Endpoint: utils.EnvOr("ARANGO_ENDPOINT", "http://localhost:8529"),
		User:     utils.EnvOrDie("ARANGO_USER"),
		Password: utils.EnvOrDie("ARANGO_PASSWORD"),
	}
}

// OpenRecords returns collection of resize results to have access to them
func OpenRecords(config *ArangoConfig) (*RecordStorage, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{config.Endpoint},
	})
	if err != nil {
		return nil, err
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.User, config.Password),
	})
	if err != nil {
		return nil, err
	}

	// Open "examples_books" database
	db, err := c.Database(ctx, Database)
	if driver.IsNotFound(err) {
		db, err = c.CreateDatabase(ctx, Database, nil)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	collection, err := db.Collection(ctx, Collection)
	if driver.IsNotFound(err) {
		collection, err = db.CreateCollection(ctx, Collection, nil)
	}

	return &RecordStorage{
		records: collection,
	}, nil
}

// StoreRecord stores links of uploaded images to db and returns them with unique record ID
// This ID could be used later to extract the record
func (rs *RecordStorage) StoreRecord(ctx context.Context, result cloud.UploadResult) (*model.ResizeResult, error) {
	record := Record{
		Images: result,
		User:   ctx.Value("user").(string),
	}
	meta, err := rs.records.CreateDocument(ctx, record)
	if err != nil {
		return nil, err
	}

	return &model.ResizeResult{
		ID:       meta.Key,
		Original: result.Original,
		Resized:  result.Resized,
	}, nil
}

// FindRecordByID extracts one record from db by its id, if it exists
func (rs *RecordStorage) FindRecordByID(ctx context.Context, id string) (*model.ResizeResult, error) {
	out := new(Record)
	meta, err := rs.records.ReadDocument(ctx, id, out)
	if err != nil {
		return nil, err
	}

	if out.User != ctx.Value("user").(string) {
		return nil, errors.New("access denied")
	}

	return &model.ResizeResult{
		ID:       meta.Key,
		Original: out.Images.Original,
		Resized:  out.Images.Resized,
	}, nil
}

// FindRecordsForUser returns all performed resizes by single user
func (rs *RecordStorage) FindRecordsForUser(ctx context.Context) ([]*model.ResizeResult, error) {
	user := ctx.Value("user").(string)

	query := fmt.Sprintf(`
		FOR record in %s
			FILTER record.user == %s
			RETURN record`,
		Collection,
		user)

	cursor, err := rs.records.Database().Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	out := []*model.ResizeResult{}
	for {
		record := new(Record)
		meta, err := cursor.ReadDocument(ctx, record)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		out = append(out, &model.ResizeResult{
			ID:       meta.Key,
			Original: record.Images.Original,
			Resized:  record.Images.Resized,
		})
	}

	return out, nil
}
