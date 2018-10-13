package main

import (
	"fmt"
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
	StudentID uint64 `gorm:"primary_key" , json:"studentID"`
	Email     string `gorm:"unique" , json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
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

func validateUser(email, password string) Student {
	//encrypt password to compare whith the already encrypted password in the database
	encryptedPassword := encryptPassword(password)

	var student Student
	//if the database returns an object that matches the user then return a success
	db.Where("Email = ? and Password = ?", email, encryptedPassword).First(&student)
	fmt.Println(email)
	return student
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect to the database")
	}
	db.AutoMigrate(&Student{})

	defer db.Close()
	router := gin.Default()

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/validateUser/:email/:password", func(c *gin.Context) {
				//get the variables from the request
				email := c.Params.ByName("email")
				password := c.Params.ByName("password")
				//				fmt.Println(email + " " + password)

				c.JSON(200, validateUser(email, password))

			})

			user.POST("/newUser", func(c *gin.Context) {
				var student Student

				student.Email = c.PostForm("email")
				student.Password = c.PostForm("password")
				student.FirstName = c.PostForm("firstName")
				student.LastName = c.PostForm("lastName")

				student.Password = encryptPassword(student.Password)

				//send email for verification--riley will insert this

				//the email will have its own call that will create the user in the database
				err := db.Create(&student).Error
				if err != nil {
					fmt.Println(err.Error())
					c.JSON(200, gin.H{"errorMsg": err})
				} else {
					db.Update()
					c.JSON(200, gin.H{"errorMsg": ""})
				}
			})

			user.POST("/emailValidation", func(c *gin.Context) {
			})

			user.POST("/deleteUser", func(c *gin.Context) {
				var student Student
				db.Where("email = ?", c.PostForm("email")).Find(&student)
				if db.Debug().First(&student, "email = ?", c.PostForm("email")).RecordNotFound() {
					db.Delete(&student)
				}
			})
		}
	}

	router.Run(":8080")
}
