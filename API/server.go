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
	db, err = gorm.Open("sqlite3", "db/test.db")
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
	db.AutoMigrate(&AdminSessionToken{})
	db.AutoMigrate(&Semester{})
	db.AutoMigrate(&Location{})
}

func main() {
	defer db.Close()

	router := gin.Default()
	api := router.Group("/api")
	{
		general := api.Group("/general")
		{
			//petty functions for drop downs :/
			general.POST("/getAllCatalogYearsForProgram")
			general.GET("/getAllSemesters", getAllSemesters)

			//searching
			general.POST("/searchForCourse", searchForCourse)
			general.GET("/searchPrograms", searchPrograms)

			//getting program requirement data
			general.GET("/getProgramRequirements", getProgramRequirements)
			general.GET("/getProgramCourseSpecificRequirementField")
			general.GET("/getProgramCourseExclusionsRequirementField")
			general.GET("/getProgramGreaterRequirementField")
		}
		user := api.Group("/user")
		{
			//basic user functions
			user.POST("/login", login)
			user.POST("/logout", logout)
			user.GET("/checkToken", checkToken)
			user.POST("/newUser", newUser)
			user.POST("/deleteUser", deleteUser)

			//current course functions
			user.POST("/enrollInCourse", enrollInCourse)
			user.POST("/dropCourse", dropCourse)
			user.GET("/getEnrolledCourses", getEnrolledCourses)

			//previous course functions
			user.POST("/addPreviousCourse", addPreviousProgram)
			user.POST("/removePreviousCourse", removePreviousProgram)

			//programs/majors/minors
			user.POST("/addStudentProgram", addProgramUser)
			user.POST("/removeStudentProgram", removeProgram)
			user.GET("/getUsersPrograms", getUsersPrograms)
			user.GET("/getRemainingProgramRequirements", getRemainingProgramRequirements)
		}

		adm := api.Group("/admin")
		{
			//general admin functions
			adm.POST("/login", admLogin)
			adm.POST("/logout", admLogout)
			adm.GET("/checkToken", checkTokenAdmin)

			//programs
			adm.POST("/addProgram", addProgram)
			adm.POST("/deleteProgram", deleteProgram)

			//requirements
			adm.POST("/addRequirementToProgram")
			adm.POST("/deleteRequirementFromProgram")

			//requirement fields
			adm.POST("/addGreaterThanRequirementToProgram", addGreaterThanRequirementToProgram)
			adm.POST("/addCourseToProgramRequirement", addCourseToProgramRequirement)
			adm.POST("/addCourseExclusionToProgram", addCourseExclusionToProgram)

			adm.POST("/deleteGreaterThanRequirement")
			adm.POST("/deleteCourseToProgramRequirement")
			adm.POST("/deleteCourseExclusionToProgram")

			//courses
			adm.POST("/addCourse", addCourse)
			adm.POST("/deleteCourse", deleteCourse)
		}
	}

	router.Run(":8080")
}
