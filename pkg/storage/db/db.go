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
	ID     string
	User   string
	Images cloud.UploadResult
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
func (rs *RecordStorage) StoreRecord(ctx context.Context, result cloud.UploadResult) (*Record, error) {
	user := ctx.Value(UserContextKey)
	if user == nil {
		return nil, errors.New("user was not specified")
	}

	record := Record{
		Images: result,
		User:   user.(string),
	}
	meta, err := rs.records.CreateDocument(ctx, record)
	if err != nil {
		return nil, err
	}

	record.ID = meta.Key
	return &record, nil
}

// FindRecordByID extracts one record from db by its id, if it exists
func (rs *RecordStorage) FindRecordByID(ctx context.Context, id string) (*Record, error) {
	record := new(Record)
	meta, err := rs.records.ReadDocument(ctx, id, record)
	if err != nil {
		return nil, err
	}

	user := ctx.Value(UserContextKey)
	if user == nil {
		return nil, errors.New("user was not specified")
	}

	if record.User != user.(string) {
		return nil, errors.New("access denied")
	}

	record.ID = meta.Key
	return record, nil
}

// FindRecordsForUser returns all performed resizes by single user
func (rs *RecordStorage) FindRecordsForUser(ctx context.Context) ([]*Record, error) {
	user := ctx.Value(UserContextKey)
	if user == nil {
		return nil, errors.New("user was not specified")
	}
	userStr := user.(string)

	query := fmt.Sprintf(`
		FOR record in %s
			FILTER record.User == @value
			RETURN record`,
		Collection)

	bindvars := map[string]interface{}{"value": userStr}
	cursor, err := rs.records.Database().Query(ctx, query, bindvars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	out := []*Record{}
	for {
		record := new(Record)
		meta, err := cursor.ReadDocument(ctx, record)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		record.ID = meta.Key
		out = append(out, record)
	}

	return out, nil
}

func imageToGraphQl(from *cloud.UploadedImage) *model.Image {
	return &model.Image{
		ImageLink: from.ImageLink,
		ExpiresAt: from.ExpiresAt,
		Width:     from.Width,
		Height:    from.Height,
	}
}
