package dbactions

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// CreateAndFillBuckets is a helper function that creates buckets for the API endpoints and fills them with some initial data
func CreateAndFillBuckets(db *bolt.DB) error {
	err := CreateBucket(db, "pingResponses")
	if err != nil {
		return err
	}

	// TODO figure out how to autoincrement keys
	PutStringIntoBucket(db, "pingResponses", "ping_first", "You rang?")    // keep prefix scans in mind
	PutStringIntoBucket(db, "pingResponses", "ping_second", "Hello there") // keep prefix scans in mind

	// TODO more of those

	return nil
}

// CreateBucket creates a single bucket if a bucket with the same name does not already exist
func CreateBucket(db *bolt.DB, bucketName string) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("bucket: %s, err: %s", bucketName, err)
		}
		return nil // return nil to committ transaction
	}); err != nil {
		return err
	}
	log.Printf("Created bucket \"%s\"", bucketName)
	return nil
}

// PutStringIntoBucket inserts a string at the provided key in the provided bucket
func PutStringIntoBucket(db *bolt.DB, bucket string, key string, value string) (err error) {
	if err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))

		// if successful this will be nil, committing the transaction
		return err // else the transaction will be rolled back
	}); err != nil {
		return err
	}
	log.Printf("Inserted \"%s\" at key \"%s\" in bucket \"%s\"", value, key, bucket)
	return nil
}

// GetStringFromBucket retrieves the string stored at the provided key in the provided bucket
func GetStringFromBucket(db *bolt.DB, bucket string, key string) (value []byte, err error) {
	if err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		value = append(value, v...) // inelegant solution to copying by unpacking v using ... notation
		log.Printf("Retrieved: \"%s\" at \"%s\", Will send \"%s\" as response", v, key, value)
		return nil
	}); err != nil {
		return nil, err
	}
	return value, nil
}

// TODO get all entries in a bucket
func GetAllFromBucket(db *bolt.DB, bucket string) {}

// DeleteBucketEntry deletes the entry stored at the provided key in the provided bucket
func DeleteBucketEntry(db *bolt.DB, bucket string, key string) (err error) {
	if err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key)) // if nil is returned deletion succeeded
	}); err != nil {
		return err
	}
	return nil
}
