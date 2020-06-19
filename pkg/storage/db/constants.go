package db

type key int

const (
	// Database is database name for arangodb
	Database = "resize_storage"

	// Collection is a collection inside Database
	Collection = "resize_records"

	// UserContextKey is used to store and extract user value from context
	UserContextKey key = iota
)
