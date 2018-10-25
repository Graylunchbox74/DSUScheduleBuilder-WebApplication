package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func createUniqueKeyAdmin(id uint64) string {
	keyToTry := randStringRunes(50)
	var adminToken AdminSessionToken
	db.Where(AdminSessionToken{AdminID: 0, Token: keyToTry}).First(&adminToken)
	for adminToken.AdminID != 0 {
		keyToTry := randStringRunes(50)
		db.Where(AdminSessionToken{AdminID: 0, Token: keyToTry}).First(&adminToken)
	}

	adminToken.AdminID = id
	adminToken.Token = keyToTry
	adminToken.TimeUpdated = time.Now()
	db.Create(&adminToken)
	return adminToken.Token
}

func findAdminGivenToken(token string) (Admin, bool) {
	var admin Admin
	var sessiontoken AdminSessionToken
	db.Where(AdminSessionToken{AdminID: 0, Token: token}).First(&sessiontoken)

	if sessiontoken.AdminID == 0 {
		return admin, false
	}

	println(sessiontoken.TimeUpdated.String())
	println(time.Now().Add(time.Minute).String())

	//session expires after a day
	if time.Now().After(sessiontoken.TimeUpdated.AddDate(0, 0, 1)) {
		return admin, true
	}

	sessiontoken.TimeUpdated = time.Now()

	db.Where(Admin{ID: sessiontoken.AdminID}).First(&admin)
	return admin, false
}

func checkTokenAdmin(c *gin.Context) {
	defer func() {
		if recover() != nil {
			c.JSON(401, gin.H{"errorMsg": "token not found"})
		}
	}()

	token := c.Request.URL.Query()["token"][0]

	admin, expired := findAdminGivenToken(token)
	if admin.ID == 0 || expired {
		c.JSON(401, gin.H{"errorMsg": "token not valid"})
	} else {
		c.JSON(200, gin.H{"errorMsg": ""})
	}
}

func admLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	admin := Admin{}
	db.Where("email = ? and password = ?", email, password).First(&admin)
	if admin.ID == 0 {
		c.JSON(400, gin.H{})
		return
	}
	token := createUniqueKeyAdmin(admin.ID)

	adminReturn := returnStudent{}
	adminReturn.Email = admin.Email
	adminReturn.FirstName = admin.FirstName
	adminReturn.LastName = admin.LastName
	adminReturn.StudentID = admin.ID
	adminReturn.Token = token

	c.JSON(200, adminReturn)
}

func admLogout(c *gin.Context) {
	token := c.PostForm("token")
	db.Where(AdminSessionToken{AdminID: 0, Token: token}).Delete(&AdminSessionToken{})
	c.JSON(200, gin.H{})
}

func addProgram(c *gin.Context) {
	token := c.PostForm("token")
	var admin Admin
	admin, isExpired := findAdminGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if admin.ID == 0 {
		c.JSON(401, gin.H{"errorMsg": "admin not found"})
		return
	}

	var program Program
	program.Major = (c.PostForm("major") == "1")
	program.Program = c.PostForm("program")
	tmpCatYear, _ := strconv.Atoi(c.PostForm("catalogYear"))
	program.CatalogYear = uint64(tmpCatYear)

	var testIfExists Program
	db.Where("Major = ? and Program = ? and Catalog_Year = ?", program.Major, program.Program, program.CatalogYear).First(&testIfExists)

	if testIfExists.ProgramID != 0 {
		c.JSON(400, gin.H{"errorMsg": "Program Already Exists"})
		return
	}

	db.Create(&program)

	c.JSON(200, gin.H{"errorMsg": ""})
}

func deleteProgram(c *gin.Context) {
	token := c.PostForm("token")
	var admin Admin
	admin, isExpired := findAdminGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if admin.ID == 0 {
		c.JSON(401, gin.H{"errorMsg": "admin not found"})
		return
	}
	var program Program
	stringID := c.PostForm("programID")
	tmpID, _ := strconv.Atoi(stringID)
	program.ProgramID = uint64(tmpID)
	db.Where("program_id = ?", program.ProgramID).First(&program)
	if program.ProgramID == 0 {
		c.JSON(400, gin.H{"errorMsg": "Program not found"})
		return
	}
	db.Delete(&program)
	c.JSON(200, gin.H{"errorMsg": ""})
}

func addCourse(c *gin.Context) {
	token := c.PostForm("token")
	var admin Admin
	admin, isExpired := findAdminGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if admin.ID == 0 {
		c.JSON(401, gin.H{"errorMsg": "admin not found"})
		return
	}
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

	tmpString = c.PostForm("startTime")
	tmp, _ = strconv.Atoi(tmpString)
	course.StartTime = uint64(tmp)

	course.Teacher = c.PostForm("teacher")

	course.Semester = c.PostForm("semester")

	db.Create(&course)

	c.JSON(200, gin.H{})
}

func deleteCourse(c *gin.Context) {
	token := c.PostForm("token")
	var admin Admin
	admin, isExpired := findAdminGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if admin.ID == 0 {
		c.JSON(401, gin.H{"errorMsg": "admin not found"})
		return
	}

	courseIDString := c.PostForm("courseID")
	tmp, _ := strconv.Atoi(courseIDString)
	var course Course
	db.Where("course_id = ?", uint64(tmp)).Delete(&course)

	c.JSON(200, gin.H{})
}

func addRequirementToProgram(c *gin.Context) {
	token := c.PostForm("token")
	var admin Admin
	admin, isExpired := findAdminGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if admin.ID == 0 {
		c.JSON(401, gin.H{"errorMsg": "admin not found"})
		return
	}
}

func addCourseToProgramRequirement(c *gin.Context) {
	token := c.PostForm("token")
	var admin Admin
	admin, isExpired := findAdminGivenToken(token)
	if isExpired {
		c.JSON(401, gin.H{"errorMsg": "token expired"})
		return
	}
	if admin.ID == 0 {
		c.JSON(401, gin.H{"errorMsg": "admin not found"})
		return
	}
	requirementIDString := c.PostForm("requirementID")
	collegeName := c.PostForm("collegeName")
	courseCodeString := c.PostForm("courseCode")

	tmp, _ := strconv.Atoi(requirementIDString)

	programRequirement := ProgramRequirement{}
	db.Where("program_requirement_id = ?", uint64(tmp)).First(&programRequirement)

	courseCodeInt, _ := strconv.Atoi(courseCodeString)
	requirementCourse := RequirementCourse{}
	requirementCourse.CourseCode = uint64(courseCodeInt)
	requirementCourse.CollegeName = collegeName

	testReqCourse := RequirementCourse{}
	db.Where(requirementCourse).First(&testReqCourse)
	if testReqCourse.CourseCode == 0 {
		db.Create(&requirementCourse)
		db.Where(requirementCourse).First(&testReqCourse)
	}
	requirementCourse = testReqCourse

	// classesAlreadyRequired := []RequirementCourse{}
	// requirementClasses := []RequirementToRequirementCourse{}
	db.Where("program_requirement_id = ?", programRequirement.ProgramID)

	c.JSON(200, gin.H{})
}
