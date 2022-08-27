package handlers

import (
	"fmt"
	"goapiproject/helpers"
	"log"
	"strconv"

	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Achievement struct {
	ID          int `storm:"id,increment"`
	Name        string
	Description string
}

// GetAllAchievements retrieves all entries in the bucket "Achievement"
// the function returns a function that satisfies Gin's router methods
// which demand a gin.HandlerFunc-compatible return value
func GetAllAchievements(db *storm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var achievements []Achievement

		err := db.All(&achievements)
		if err != nil {
			handleError(c, err, 500, "Failed to retrieve achievements")
		}

		c.JSON(200, gin.H{
			"achievements": achievements,
		})
	}

	return gin.HandlerFunc(fn)
}

// GetOneAchievement retrieves one entry in the bucket "Achievement"
func GetOneAchievement(db *storm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Query("id")
		name := c.Query("name")
		var achievement Achievement

		if id != "" && id != "0" { // if we have some identifiable information
			id, _ := strconv.Atoi(id)

			err := db.One("ID", id, &achievement)
			if err != nil || achievement.ID == 0 {
				handleError(c, err, 400, "the ID you provided didn't exist in the database")

			} else {
				c.JSON(200, gin.H{
					"achievement": achievement,
				})
			}

		} else if name != "" { // if we have some identifiable information
			err := db.One("Name", name, &achievement)
			if err != nil {
				handleError(c, err, 400)
			}

			c.JSON(200, gin.H{
				"achievement": achievement,
			})

		} else {
			handleError(c, fmt.Errorf("please provide an achievement ID >=1, or a title-cased name"), 400)
		}
	}

	return gin.HandlerFunc(fn)
}

// CreateANewAchievement creates a new entry in the bucket "Achievement"
func CreateANewAchievement(db *storm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		name := c.Query("name")
		description := c.Query("description")
		newAchievement := Achievement{
			Name:        cases.Title(language.English).String(name),
			Description: description,
		}

		var existingAchievement Achievement
		err := db.One("Name", newAchievement.Name, &existingAchievement)
		if err != nil { // err will exist if no clashing entry exists
			log.Printf("Found no achivement under \"%v\" in the database, creation can safely proceed: %v", newAchievement.Name, err.Error())
		}

		if newAchievement.Name == existingAchievement.Name { // if the error is nil we got a match, and have to say no to the requestee
			handleError(c, err, 409, "An entry using the provided name already exists, select another name")

		} else if err = db.Save(&newAchievement); err != nil { // if the above check doesn't run we can go ahead and attempt to save
			handleError(c, err, 500, "Failed to save new achievement")
		} else { // we're good to go as the new entry got saved, and we'll provide it back to the requestee
			var responseAchievement Achievement
			if err = db.One("Name", name, &responseAchievement); err != nil {
				handleError(c, err, 500, "Failed to fetch newly created achievement")
			}
			c.JSON(200, gin.H{
				"message":     "successfully saved new achievement '" + name + "' as ID: " + strconv.Itoa(responseAchievement.ID),
				"achievement": responseAchievement,
			})
		}
	}

	return gin.HandlerFunc(fn)
}

// UpdateExistingAchievement updates an entry in the bucket "Achievement"
func UpdateExistingAchievement(db *storm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Query("id")
		newName := c.Query("name")
		newDescription := c.Query("description")
		var achievementToUpdate Achievement
		var responseAchievement Achievement

		// ensure that we can work with what we got
		if id == "" || id == "0" {
			c.JSON(400, gin.H{
				"message": "Please provide an ID >=1 along with the new name or description for the achievement.",
			})

		} else { // find existing entry by id
			id, _ := strconv.Atoi(id)
			if err := db.One("ID", id, &achievementToUpdate); err != nil {
				handleError(c, err, 500)
			}

			// change it
			if newName != "" {
				achievementToUpdate.Name = newName
			}
			if newDescription != "" {
				achievementToUpdate.Description = newDescription
			}

			// save it
			err := db.Save(&achievementToUpdate)
			if err != nil {
				handleError(c, err, 500)
			}

			// retrieve the now updated achievement
			err = db.One("ID", id, &responseAchievement)
			if err != nil {
				handleError(c, err, 500)
			}
			c.JSON(200, gin.H{ // and return it to the user
				"message":     "successfully updated achievement to ID: '" + strconv.Itoa(responseAchievement.ID) + " Name: '" + responseAchievement.Name + "' Description: " + responseAchievement.Description,
				"achievement": responseAchievement,
			})
		}
	}

	return gin.HandlerFunc(fn)
}

// DeleteExistingAchievement deletes an existing entry in the bucket "Achievement"
func DeleteExistingAchievement(db *storm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Query("id")
		intID, _ := strconv.Atoi(id)
		var achievementToDelete = Achievement{ID: intID}

		if err := db.DeleteStruct(&achievementToDelete); err != nil {
			handleError(c, err, 400)
		}

		c.JSON(200, gin.H{
			"message": "successfully deleted achievement with ID: '" + id,
		})
	}

	return gin.HandlerFunc(fn)
}

// ResetDatabase drops the database, recreates it, and puts some initial data into it
// This puts the database back in the same state that it was when the API server was started
func ResetDatabase(db *storm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		err := helpers.ResetAndFillDB(db)
		if err != nil {
			handleError(c, err, 500)
		}

		c.JSON(200, gin.H{
			"message": "Successfully reset the database to its initial state",
		})
	}

	return gin.HandlerFunc(fn)
}

// Online returns a reply to any GET request as a way to verify that the server is up and reachable
func Online() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I'm alive!",
		})
	}
	return gin.HandlerFunc(fn)
}

// handleError is a VERY simple approach to hiding some of the repetitive error handling logic
func handleError(c *gin.Context, err error, statusCode int, additionalInfo ...string) bool {
	if err != nil {
		c.JSON(statusCode, gin.H{
			"error":           err.Error(),
			"additional_info": additionalInfo,
		})
		return true
	}
	return false
}
