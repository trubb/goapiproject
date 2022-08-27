package helpers

import (
	"errors"
	"net"

	"github.com/asdine/storm/v3"
	bolt "go.etcd.io/bbolt"
)

type Achievement struct {
	ID          int `storm:"id,increment"`
	Name        string
	Description string
}

// GetServersOwnIP is a hacky solution to the problem of retrieving the IP of the equipment that is running this program
func GetServersOwnIP() (net.IP, error) {
	// Attempt to connect to <some> address
	// It doesn't need to exist so we're using an IP from TEST-NET-1 (192.0.2.0/24)
	connection, err := net.Dial("udp", "192.0.2.1:80")
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	// return the IP of whatever we're running this program on
	return connection.LocalAddr().(*net.UDPAddr).IP, nil
}

// ResetAndFillDB is a helper function that clears the DB, and then fills it with some initial data
func ResetAndFillDB(db *storm.DB) (err error) {

	// drop the one existing bucket
	if err = db.Drop(&Achievement{}); err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
		return err
	}

	// fill
	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	achievement2 := Achievement{
		Name:        "Produce",
		Description: "Produced more than 200 items across all scenarios",
	}
	achievement3 := Achievement{
		Name:        "Showtime",
		Description: "Put on a show with the combine harvester",
	}

	if err = db.Save(&achievement1); err != nil {
		return err
	}
	if err = db.Save(&achievement2); err != nil {
		return err
	}
	if err = db.Save(&achievement3); err != nil {
		return err
	}

	return nil
}
