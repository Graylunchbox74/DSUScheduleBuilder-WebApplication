package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Admin struct {
	ID        uint64 `gorm:"primary_key"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Password  string
}

type Student struct {
	StudentID uint64 `gorm:"primary_key"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Password  string
}

type StudentPrograms struct {
	ProgramID uint64 `gorm:"foreignkey:ProgramID;association_foreignkey:ProgramID"`
	StudentID uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
}

type Program struct {
	ProgramID   uint64 `gorm:"primary_key"`
	CatalogYear uint64
	Major       bool
	Program     uint64
}

type StudentProgram struct {
	ProgramID   uint64 // Comes from ProgramID of Program
	StudentID   uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
	CatalogYear uint64
	Major       bool
	Program     uint64
}

type ProgramRequirement struct {
	ProgramID    uint64 `gorm:"foreignkey:ProgramID;association_foreignkey:ProgramID"`
	CourseID     uint64 `gorm:"foreignkey:CourseID;association_foreignkey:CourseID"`
	NumberToTake uint64
}

type StudentProgramRequirement struct {
	ProgramID    uint64 `gorm:"foreignkey:ProgramID;association_foreignkey:ProgramID"`
	CourseID     uint64 `gorm:"foreignkey:CourseID;association_foreignkey:CourseID"`
	StudentID    uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
	NumberToTake uint64
}

type Course struct {
	CourseID    uint64 `gorm:"primary_key"`
	CourseCode  uint64
	CourseName  string
	Credits     uint64
	DaysOfWeek  uint64
	EndTime     uint64
	StartTime   uint64
	EndDate     time.Time
	StartDate   time.Time
	CollegeName string
	Location    string
	Teacher     string // ADD IN TEACHER INFO
}

func encryptPassword(password string) string {
	return password
}

func validateUser(email, password string) bool {
	//encrypt password to compare whith the already encrypted password in the database
	encryptedPassword := encryptPassword(password)

	var student Student
	db.Where(&student, "Email = $1 and Password = $2", email, encryptedPassword).Find(&student)

	//if the database returns an object that matches the user then return a success
	if student != nil {
		return true
	} else {
		return false
	}
}

var db *gorm.DB

func init() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to the database")
	}
}

func main() {
	defer db.Close()
	router := gin.Default()

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/validateUser/:email/:password", func(c *gin.Context) {
				//get the variables from the request
				email := c.PostForm("email")
				password := c.PostForm("password")

				if validateUser(email, password) {
					c.JSON(200, gin.H{"success": 1})
				} else {
					c.JSON(200, gin.H{"success": 0})
				}
			})

			user.POST("/newUser", func(c *gin.Context) {
				var student Student

				student.Email = c.PostForm("email")
				student.Password = c.PostForm("password")
				student.FirstName = c.PostForm("firstName")
				student.LastName = c.PostForm("lastName")

				student.Password = encryptPassword(student.Password)

				db.Where("email = ?", student.Email).Find(&student)
				if rows != nil {
					c.JSON(200, gin.H{"errorMsg": "email already exists"})
				}

				//send email for verification--riley will insert this

				//the email will have its own call that will create the user in the database
				db.Create(student)
				c.JSON(200, gin.H{"errorMsg": ""})
			})

			user.POST("/emailValidation", func(c *gin.Context) {
			})

			user.POST("/deleteUser", func(c *gin.Context) {
				var student Student
				db.Where("email = ?", c.PostForm("email")).Find(&student)
				if student != nil {
					db.Delete(&student)
				}
			})
		}
	}

	r.Run(":80")
}
