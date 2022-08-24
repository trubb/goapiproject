package dbactions

import (
	"log"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

// test creating a db, inserting an item, retrieving it, updating it, and then deleting it using our wrappers
func TestDBFuncs(t *testing.T) {
	// set up connection to db file
	db, err := bolt.Open("dbactions_test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // close to release file lock

	// should succeed
	err = CreateAndFillBuckets(db)
	if err != nil {
		t.Errorf(err.Error())
	}

	// should succeed
	err = PutStringIntoBucket(db, "pingResponses", "ping_first", "new value")
	if err != nil {
		t.Errorf(err.Error())
	}

	// should succeed
	retrievedString, err := GetStringFromBucket(db, "pingResponses", "ping_first")
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Printf("Retrieved %s", retrievedString)

	// should succeed
	err = DeleteBucketEntry(db, "pingResponses", "ping_first")
	if err != nil {
		t.Errorf(err.Error())
	}

	// should fail
	failGetString, _ := GetStringFromBucket(db, "pingResponses", "ping_last")
	if len(failGetString) > 0 {
		t.Errorf("found %s when expecting to not find anything", failGetString)
	}

	// should fail
	err = DeleteBucketEntry(db, "pingResponses", "ping_last")
	if err == nil {
		t.Logf("deletion of nonexistent key returned nil as expected: %s", err)
	}

}
