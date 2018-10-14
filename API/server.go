package main

import (
	"fmt"
	"strconv"
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

type StudentToProgram struct {
	ProgramID uint64 `gorm:"foreignkey:ProgramID;association_foreignkey:ProgramID"`
	StudentID uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
}

type Program struct {
	ProgramID   uint64 `gorm:"primary_key"`
	CatalogYear uint64
	Major       bool
	Program     string
}

type StudentProgram struct {
	ProgramID   uint64 // Comes from ProgramID of Program
	StudentID   uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
	CatalogYear uint64
	Major       bool
	Program     string
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

func findUserWithID(id int) Student {
	var student Student
	db.Where("student_id = ?", id).First(&student)
	fmt.Println(student.Email)
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
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&StudentToProgram{})
	db.AutoMigrate(&StudentProgram{})
	db.AutoMigrate(&Program{})

	defer db.Close()
	router := gin.Default()

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/validateUser", func(c *gin.Context) {
				//get the variables from the request
				email := c.PostForm("email")
				password := c.PostForm("password")

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
				db.Delete(&student)
			})

			user.POST("/addStudentProgram", func(c *gin.Context) {
				stringID := c.PostForm("studentID")
				id, _ := strconv.Atoi(stringID)
				var student Student
				student = findUserWithID(id)

				if student.StudentID == 0 {
					c.JSON(200, gin.H{"errorMsg": "student not found"})
					return
				}

				programIDString := c.PostForm("programID")
				var program Program
				db.Where("program_id = ?", programIDString).First(&program)

				if program.ProgramID == 0 {
					c.JSON(200, gin.H{"errorMsg": "program not found"})
					return
				}

				var studentPrograms StudentToProgram
				studentPrograms.StudentID = student.StudentID

				studentPrograms.ProgramID = program.ProgramID

				var testStudentPrograms StudentToProgram
				db.Where("student_id = ? and program_id = ?", studentPrograms.StudentID, studentPrograms.ProgramID).First(&testStudentPrograms)
				if testStudentPrograms.ProgramID != 0 {
					c.JSON(200, gin.H{"errorMsg": "Student already enrolled in this program"})
					return
				}

				db.Create(&studentPrograms)

				var studentProgram StudentProgram
				studentProgram.CatalogYear = program.CatalogYear
				studentProgram.Major = program.Major
				studentProgram.Program = program.Program
				studentProgram.ProgramID = program.ProgramID
				studentProgram.StudentID = student.StudentID

				db.Create(&studentProgram)

				c.JSON(200, gin.H{"errorMsg": ""})
			})
		}
		adm := api.Group("/admin")
		{
			adm.POST("/addProgram", func(c *gin.Context) {
				var program Program
				program.Major = (c.PostForm("major") == "1")
				program.Program = c.PostForm("program")
				tmpCatYear, _ := strconv.Atoi(c.PostForm("catalogYear"))
				program.CatalogYear = uint64(tmpCatYear)

				var testIfExists Program
				db.Where("Major = ? and Program = ? and Catalog_Year = ?", program.Major, program.Program, program.CatalogYear).First(&testIfExists)

				if testIfExists.ProgramID != 0 {
					c.JSON(200, gin.H{"errorMsg": "Program Already Exists"})
					return
				}

				db.Create(&program)

				c.JSON(200, gin.H{"errorMsg": ""})
			})
		}
	}

	router.Run(":8080")
}
