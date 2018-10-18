package main

import (
	"fmt"
	"math/rand"
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

type SessionToken struct {
	StudentID uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
	Token     string `gorm:"unique"`
}

type Student struct {
	StudentID uint64 `gorm:"primary_key" , json:"studentID"`
	Email     string `gorm:"unique" , json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type StudentToCourse struct {
	CourseID  uint64 `gorm:"foreignkey:courseID;association_foreignkey:CourseID"`
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

type returnStudent struct {
	StudentID uint64 `json:"studentID"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
}

//fundamental functions
func encryptPassword(password string) string {
	return password
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func createUniqueKey(id uint64) string {
	keyToTry := randStringRunes(30)
	var studentToken SessionToken
	db.Where(SessionToken{StudentID: 0, Token: keyToTry}).First(&studentToken)
	for studentToken.StudentID != 0 {
		keyToTry := randStringRunes(30)
		db.Where(SessionToken{StudentID: 0, Token: keyToTry}).First(&studentToken)
	}

	studentToken.StudentID = id
	studentToken.Token = keyToTry
	db.Create(&studentToken)
	return studentToken.Token
}

//functions specific for requests
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

func findStudentGivenToken(token string) Student {
	var student Student
	var sessiontoken SessionToken
	db.Where(SessionToken{StudentID: 0, Token: token}).First(&sessiontoken)

	if sessiontoken.StudentID == 0 {
		return student
	}

	db.Where(Student{StudentID: sessiontoken.StudentID}).First(&student)
	return student
}

var db *gorm.DB

func main() {
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

	defer db.Close()
	router := gin.Default()

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{

			user.POST("/login", func(c *gin.Context) {
				//get the variables from the request
				email := c.PostForm("email")
				password := c.PostForm("password")
				studentInformation := validateUser(email, password)
				var studentReturn returnStudent

				studentReturn.Email = studentInformation.Email
				studentReturn.FirstName = studentInformation.FirstName
				studentReturn.LastName = studentInformation.LastName
				studentReturn.StudentID = studentInformation.StudentID
				studentReturn.Token = createUniqueKey(studentInformation.StudentID)

				c.JSON(200, studentReturn)
			})

			user.POST("/logout", func(c *gin.Context) {
				token := c.PostForm("token")
				db.Where(SessionToken{StudentID: 0, Token: token}).Delete(&SessionToken{})
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
				token := c.PostForm("token")
				var student Student
				student = findStudentGivenToken(token)
				if student.StudentID == 0 {
					c.JSON(401, gin.H{"errorMsg": "student not found"})
					return
				}

				db.Delete(&student)
			})

			user.POST("/addStudentProgram", func(c *gin.Context) {
				token := c.PostForm("token")
				var student Student
				student = findStudentGivenToken(token)
				if student.StudentID == 0 {
					c.JSON(401, gin.H{"errorMsg": "student not found"})
					return
				}

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

				var studentPrograms StudentProgram
				studentPrograms.StudentID = student.StudentID

				studentPrograms.ProgramID = program.ProgramID

				var testStudentPrograms StudentProgram
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

			user.POST("/removeStudentProgram", func(c *gin.Context) {
				token := c.PostForm("token")
				var student Student
				student = findStudentGivenToken(token)
				if student.StudentID == 0 {
					c.JSON(401, gin.H{"errorMsg": "student not found"})
					return
				}

				var program Program
				stringProgramID := c.PostForm("programID")
				normalIntID, _ := strconv.Atoi(stringProgramID)
				db.Where("program_id = ?", normalIntID).First(&program)
				if program.ProgramID == 0 {
					c.JSON(200, gin.H{"errMsg": "Program not found"})
					return
				}

				var studentToProgram StudentProgram
				db.Where("program_id = ? and  student_id = ?", program.ProgramID, student.StudentID).First(&studentToProgram)
				if studentToProgram.StudentID == 0 {
					c.JSON(200, gin.H{"errorMsg": "Student is not currently in this program"})
					return
				}
				db.Delete(&studentToProgram)

				var studentProgram StudentProgram
				db.Where("program_id = ? and student_id = ?", program.ProgramID, student.StudentID).First(&studentProgram)
				db.Delete(studentProgram)
				c.JSON(200, gin.H{"errorMsg": ""})
			})

			user.GET("/getEnrolledCourses/:studentID", func(c *gin.Context) {
				var courses []Course
				var studentToCourses []StudentToCourse

				token := c.Params.ByName("token")
				var student Student
				student = findStudentGivenToken(token)
				if student.StudentID == 0 {
					c.JSON(401, gin.H{"errorMsg": "student not found"})
					return
				}
				db.Where("student_id = ?", student.StudentID).Find(&studentToCourses)
				if len(studentToCourses) == 0 {
					c.JSON(200, courses)
				}

				var singleCourse Course
				for _, current := range studentToCourses {
					singleCourse = Course{}
					db.Where("course_id = ?", current.CourseID).Find(&singleCourse)
					if singleCourse.CourseID != 0 {
						courses = append(courses, singleCourse)
					}
				}

				c.JSON(200, courses)

				// c.JSON(200, courses)
			})

			user.POST("/enrollInCourse", func(c *gin.Context) {
				token := c.PostForm("token")
				var student Student
				student = findStudentGivenToken(token)
				if student.StudentID == 0 {
					c.JSON(401, gin.H{"errorMsg": "student not found"})
					return
				}
				courseID, _ := strconv.Atoi(c.PostForm("courseID"))

				//check that student is not already enrolled in this course
				var studentToCourse StudentToCourse
				db.Where("student_id = ? and course_id = ?", student.StudentID, courseID).First(&studentToCourse)
				if studentToCourse.CourseID != 0 {
					c.JSON(200, gin.H{"errorMsg": "Student already enrolled in this course"})
					return
				}

				//get the course to make sure it exists
				var course Course
				db.Where("course_id = ?", courseID).First(&course)
				if course.CourseID == 0 {
					c.JSON(200, gin.H{"errorMsg": "Course does not exist with this id"})
					return
				}
				//now we need to make sure that the class does not conflict with other classes
				var studentToCourses []StudentToCourse
				var courseToCompare Course
				db.Where("student_id = ?", student.StudentID).Find(&studentToCourses)
				for _, currentCourse := range studentToCourses {
					db.Where("course_id = ?", currentCourse.CourseID).First(&courseToCompare)

					//also need to check for dates on this

					//check if times interfear at all
					if (courseToCompare.StartTime <= course.StartTime && courseToCompare.EndTime >= course.StartTime) ||
						(courseToCompare.StartTime <= course.EndTime && courseToCompare.EndTime >= course.EndTime) ||
						(courseToCompare.StartTime <= course.StartTime && courseToCompare.EndTime >= course.EndTime) ||
						(courseToCompare.StartTime >= course.StartTime && courseToCompare.EndTime <= course.EndTime) {
						//check if they happen on the same day
						daysForRegisteringCourse := strconv.Itoa(int(course.DaysOfWeek))
						daysForCurrentCourse := strconv.Itoa(int(courseToCompare.DaysOfWeek))
						for i := 0; i < 5; i++ {
							if daysForCurrentCourse[i] == daysForRegisteringCourse[i] && daysForCurrentCourse[i] == '1' {
								c.JSON(200, gin.H{"errorMsg": "Course conflicts with " + courseToCompare.CourseName})
								return
							}
						}
					}
				}

				//do some major stuff here

				//actually add the course
				var newStudentToCourse StudentToCourse
				newStudentToCourse.CourseID = uint64(courseID)
				newStudentToCourse.StudentID = uint64(student.StudentID)

				db.Create(&newStudentToCourse)
				c.JSON(200, gin.H{"errorMsg": ""})
			})

			user.POST("/dropCourse", func(c *gin.Context) {
				token := c.PostForm("token")
				var student Student
				student = findStudentGivenToken(token)
				if student.StudentID == 0 {
					c.JSON(401, gin.H{"errorMsg": "student not found"})
					return
				}
				courseID, _ := strconv.Atoi(c.PostForm("courseID"))
				db.Where("student_id = ? and course_id = ?", student.StudentID, courseID).Delete(&StudentToCourse{})
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

			adm.POST("/deleteProgram", func(c *gin.Context) {
				var program Program
				stringID := c.PostForm("programID")
				tmpID, _ := strconv.Atoi(stringID)
				program.ProgramID = uint64(tmpID)
				db.Where("program_id = ?", program.ProgramID).First(&program)
				if program.ProgramID == 0 {
					c.JSON(200, gin.H{"errorMsg": "Program not found"})
					return
				}
				db.Delete(&program)
				c.JSON(200, gin.H{"errorMsg": ""})
			})

			adm.POST("/addCourse", func(c *gin.Context) {
				var course Course
				course.CollegeName = c.PostForm("collegeName")

				tmpString := c.PostForm("courseCode")
				tmp, _ := strconv.Atoi(tmpString)
				course.CourseCode = uint64(tmp)

				course.CourseName = c.PostForm("courseName")

				tmpString = c.PostForm("credits")
				tmp, _ = strconv.Atoi(tmpString)
				course.Credits = uint64(tmp)

				tmpString = c.PostForm("daysOfWeek")
				tmp, _ = strconv.Atoi(tmpString)
				course.DaysOfWeek = uint64(tmp)

				tmpString = c.PostForm("endTime")
				tmp, _ = strconv.Atoi(tmpString)
				course.EndTime = uint64(tmp)

				course.Location = c.PostForm("location")

				tmpString = c.PostForm("credits")
				tmp, _ = strconv.Atoi(tmpString)
				course.StartTime = uint64(tmp)

				course.Teacher = c.PostForm("teacher")

				db.Create(&course)
			})
		}
	}

	router.Run(":8080")
}
