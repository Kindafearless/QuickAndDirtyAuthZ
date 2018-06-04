package database

import "github.com/hashicorp/go-memdb"

type InMemory struct {
	DB *memdb.MemDB
}

type Object struct {
	Key       string
	Value     string
}

// Create the high level schema for in-memory database based on the Object struct
func CreateSchema(tableName string, usersTableName string) *memdb.DBSchema {

	// This refers to the object value that will be used as the indexer field in the database
	idField := "Key"

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableName: &memdb.TableSchema{
				Name: tableName,
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: idField},
					},
				},
			},
			usersTableName: &memdb.TableSchema{
				Name: usersTableName,
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: idField},
					},
				},
			},
		},
	}

	// Create the DB schema
	return schema
}

// Create a new in-memory database based on passed schema
func CreateDatabase(schema *memdb.DBSchema) (*memdb.MemDB, error) {

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (iM *InMemory) WriteToDatabase(tableName string, key string, value string) error {

	// Create a write transaction
	txn := iM.DB.Txn(true)

	obj := &Object{key, value}

	// Check if object exists in database
	exists, err := txn.First(tableName, "id", key)
	if err != nil {
		println(err.Error())
	}

	if exists == nil {

		// Insert a new object
		if err := txn.Insert(tableName, obj); err != nil {
			return err
		}
	}

	// Commit the transaction
	txn.Commit()

	return nil
}

func (iM *InMemory) GetFromDatabase(tableName string, key string) (*Object, error) {

	// Create read-only transaction
	txn := iM.DB.Txn(false)
	defer txn.Abort()

	// Check if object exists in database
	exists, err := txn.First(tableName, "id", key)
	if err != nil {
		return nil, err
	}

	obj := exists.(*Object)

	return obj, nil
}

func (iM *InMemory) DeleteFromDatabase(tableName string, key string) error {

	// Create a write transaction
	txn := iM.DB.Txn(true)

	// Check if object exists in database
	exists, err := txn.First(tableName, "id", key)
	if err != nil {
		println(err.Error())
	}

	if exists != nil {

		// Delete object
		if err := txn.Delete(tableName, exists); err != nil {
			return err
		}
	}

	// Commit the transaction
	txn.Commit()

	return nil
}

func (iM *InMemory) ListDatabase(tableName string) ([]*Object, error) {

	// Create read-only transaction
	txn := iM.DB.Txn(false)
	defer txn.Abort()

	values, err := txn.Get(tableName, "id")
	if err != nil {
		return nil, err
	}

	var list []*Object

	value := values.Next()
	for value != nil {

		list = append(list, value.(*Object))
		value = values.Next()
	}

	return list, nil
}