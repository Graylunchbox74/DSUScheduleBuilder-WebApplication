package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to the database")
	}
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&SessionToken{})
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&StudentProgram{})
	db.AutoMigrate(&Program{})
	db.AutoMigrate(&StudentToCourse{})
	db.AutoMigrate(&Course{})

}

func main() {
	defer db.Close()

	router := gin.Default()
	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/login", login)
			user.POST("/logout", logout)
			user.POST("/newUser", newUser)
			user.POST("/deleteUser", deleteUser)
			user.POST("/addStudentProgram", addProgramUser)
			user.POST("/removeStudentProgram", removeProgram)
			user.POST("/enrollInCourse", enrollInCourse)
			user.POST("/dropCourse", dropCourse)
			user.POST("/searchForCourse", searchForCourse)
			user.GET("/getEnrolledCourses", getEnrolledCourses)
		}
		adm := api.Group("/admin")
		{
			adm.POST("/addProgram", addProgram)
			adm.POST("/deleteProgram", deleteProgram)
			adm.POST("/addCourse", addCourse)
		}
	}

	router.Run(":8080")
}
