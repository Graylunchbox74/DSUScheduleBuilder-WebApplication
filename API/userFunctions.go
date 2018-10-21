package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
	keyToTry := randStringRunes(50)
	var studentToken SessionToken
	db.Where(SessionToken{StudentID: 0, Token: keyToTry}).First(&studentToken)
	for studentToken.StudentID != 0 {
		keyToTry := randStringRunes(50)
		db.Where(SessionToken{StudentID: 0, Token: keyToTry}).First(&studentToken)
	}

	studentToken.StudentID = id
	studentToken.Token = keyToTry
	studentToken.TimeUpdated = time.Now()
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

func findStudentGivenToken(token string) (Student, bool) {
	var student Student
	var sessiontoken SessionToken
	db.Where(SessionToken{StudentID: 0, Token: token}).First(&sessiontoken)

	if sessiontoken.StudentID == 0 {
		return student, false
	}

	println(sessiontoken.TimeUpdated.String())
	println(time.Now().Add(time.Minute).String())

	//session expires after a day
	if time.Now().After(sessiontoken.TimeUpdated.AddDate(0, 0, 1)) {
		return student, true
	}

	sessiontoken.TimeUpdated = time.Now()

	db.Where(Student{StudentID: sessiontoken.StudentID}).First(&student)
	return student, false
}

func login(c *gin.Context) {
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
}

func logout(c *gin.Context) {
	token := c.PostForm("token")
	db.Where(SessionToken{StudentID: 0, Token: token}).Delete(&SessionToken{})
	c.JSON(200, gin.H{})
}

func newUser(c *gin.Context) {
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
		c.JSON(400, gin.H{"errorMsg": err})
	} else {
		db.Update()
		c.JSON(200, gin.H{"errorMsg": ""})
	}
}

func deleteUser(c *gin.Context) {
	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}
	db.Delete(&student)
}

func addProgramUser(c *gin.Context) {
	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}

	programIDString := c.PostForm("programID")
	var program Program
	db.Where("program_id = ?", programIDString).First(&program)

	if program.ProgramID == 0 {
		c.JSON(400, gin.H{"errorMsg": "program not found"})
		return
	}

	var studentPrograms StudentProgram
	studentPrograms.StudentID = student.StudentID

	studentPrograms.ProgramID = program.ProgramID

	var testStudentPrograms StudentProgram
	db.Where("student_id = ? and program_id = ?", studentPrograms.StudentID, studentPrograms.ProgramID).First(&testStudentPrograms)
	if testStudentPrograms.ProgramID != 0 {
		c.JSON(400, gin.H{"errorMsg": "Student already enrolled in this program"})
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
}

func removeProgram(c *gin.Context) {
	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
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
		c.JSON(400, gin.H{"errorMsg": "Student is not currently in this program"})
		return
	}
	db.Delete(&studentToProgram)

	var studentProgram StudentProgram
	db.Where("program_id = ? and student_id = ?", program.ProgramID, student.StudentID).First(&studentProgram)
	db.Delete(studentProgram)
	c.JSON(200, gin.H{"errorMsg": ""})
}

func getEnrolledCourses(c *gin.Context) {
	courses := []Course{}
	var studentToCourses []StudentToCourse

	token := c.Request.URL.Query()["token"][0]

	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}
	db.Where("student_id = ?", student.StudentID).Find(&studentToCourses)

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
}

func enrollInCourse(c *gin.Context) {
	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}
	courseID, _ := strconv.Atoi(c.PostForm("courseID"))

	//check that student is not already enrolled in this course
	var studentToCourse StudentToCourse
	db.Where("student_id = ? and course_id = ?", student.StudentID, courseID).First(&studentToCourse)
	if studentToCourse.CourseID != 0 {
		c.JSON(400, gin.H{"errorMsg": "Student already enrolled in this course"})
		return
	}

	//get the course to make sure it exists
	var course Course
	db.Where("course_id = ?", courseID).First(&course)
	if course.CourseID == 0 {
		c.JSON(400, gin.H{"errorMsg": "Course does not exist with this id"})
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
					c.JSON(400, gin.H{"conflicts": true})
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
}

func dropCourse(c *gin.Context) {
	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}
	courseID, _ := strconv.Atoi(c.PostForm("courseID"))
	db.Where("student_id = ? and course_id = ?", student.StudentID, courseID).Delete(&StudentToCourse{})
}

func searchForCourse(c *gin.Context) {

	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}

	var course Course

	course.CollegeName = c.PostForm("collegeName")
	courseCode := c.PostForm("courseCode")
	CourseCodeint, _ := strconv.Atoi(courseCode)
	course.CourseCode = uint64(CourseCodeint)
	course.Teacher = c.PostForm("teacherName")
	course.CourseName = c.PostForm("courseName")
	course.Semester = c.PostForm("semester")
	course.Location = c.PostForm("location")

	returnCourses := []Course{}

	db.Where(course).Find(&returnCourses)

	c.JSON(200, returnCourses)
}

func getProgramRequirements(c *gin.Context) {
	token := c.PostForm("token")
	var student Student
	student, isExpired := findStudentGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if student.StudentID == 0 {
		c.JSON(401, gin.H{"errorMsg": "student not found"})
		return
	}
}
