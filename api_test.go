package main

import (
	"goapiproject/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestPingRoute(t *testing.T) {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/online", handlers.Online())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/online", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"I'm alive!\"}", w.Body.String())
}

func TestDatabaseAPIGetAll(t *testing.T) {
	db, err := storm.Open("getall_test.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close() // close database connection to release file lock

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/achievement/all", handlers.GetAllAchievements(db))

	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	if err = db.Save(&achievement1); err != nil {
		t.Errorf(err.Error())
	}

	// get all
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/achievement/all", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"achievements\":[{\"ID\":1,\"Name\":\"Explore\",\"Description\":\"Explored many parts of the world\"}]}", w.Body.String())

}

func TestDatabaseAPIGetOne(t *testing.T) {
	db, err := storm.Open("getone_test.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close() // close database connection to release file lock

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/achievement/one", handlers.GetOneAchievement(db))

	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	if err = db.Save(&achievement1); err != nil {
		t.Errorf(err.Error())
	}

	// get one by id
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/achievement/one?id=1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"achievement\":{\"ID\":1,\"Name\":\"Explore\",\"Description\":\"Explored many parts of the world\"}}", w.Body.String())

	// get one by name
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/achievement/one?name=Explore", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"achievement\":{\"ID\":1,\"Name\":\"Explore\",\"Description\":\"Explored many parts of the world\"}}", w.Body.String())
}

func TestDatabaseAPICreate(t *testing.T) {
	db, err := storm.Open("create_test.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close() // close database connection to release file lock

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.POST("/achievement/create", handlers.CreateANewAchievement(db))

	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	if err = db.Save(&achievement1); err != nil {
		t.Errorf(err.Error())
	}

	// create new
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/achievement/create?name=Hello&description=Waved to a pedestrian while mounted in a vehicle", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"achievement\":{\"ID\":2,\"Name\":\"Hello\",\"Description\":\"Waved to a pedestrian while mounted in a vehicle\"},\"message\":\"successfully saved new achievement 'Hello' as ID: 2\"}", w.Body.String())
}

func TestDatabaseAPIUpdate(t *testing.T) {
	db, err := storm.Open("update_test.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close() // close database connection to release file lock

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.PUT("/achievement/update", handlers.UpdateExistingAchievement(db))

	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	if err = db.Save(&achievement1); err != nil {
		t.Errorf(err.Error())
	}

	// update
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/achievement/update?id=1&name=Updated&description=Updated description", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"achievement\":{\"ID\":1,\"Name\":\"Updated\",\"Description\":\"Updated description\"},\"message\":\"successfully updated achievement to ID: '1 Name: 'Updated' Description: Updated description\"}", w.Body.String())
}

func TestDatabaseAPIDelete(t *testing.T) {
	db, err := storm.Open("delete_test.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close() // close database connection to release file lock

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.DELETE("/achievement/delete", handlers.DeleteExistingAchievement(db))

	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	if err = db.Save(&achievement1); err != nil {
		t.Errorf(err.Error())
	}

	// delete
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/achievement/delete?id=1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"successfully deleted achievement with ID: '1\"}", w.Body.String())
}

func TestDatabaseAPIReset(t *testing.T) {
	db, err := storm.Open("reset_test.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close() // close database connection to release file lock

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.PUT("/achievement/reset/really", handlers.ResetDatabase(db))
	router.GET("/achievement/all", handlers.GetAllAchievements(db))

	achievement1 := Achievement{
		Name:        "Explore",
		Description: "Explored many parts of the world",
	}
	if err = db.Save(&achievement1); err != nil {
		t.Errorf(err.Error())
	}

	// reset
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/achievement/reset/really", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"Successfully reset the database to its initial state\"}", w.Body.String())

	// get all
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/achievement/all", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"achievements\":[{\"ID\":1,\"Name\":\"Explore\",\"Description\":\"Explored many parts of the world\"},{\"ID\":2,\"Name\":\"Produce\",\"Description\":\"Produced more than 200 items across all scenarios\"},{\"ID\":3,\"Name\":\"Showtime\",\"Description\":\"Put on a show with the combine harvester\"}]}", w.Body.String())
}
