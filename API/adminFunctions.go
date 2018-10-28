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

	password = encryptPassword(password)

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

	newRequirment := ProgramRequirement{}

	newRequirment.RequirementName = c.PostForm("requirementName")
	numberToTakeString := c.PostForm("numberToTake")
	numberToTake, _ := strconv.Atoi(numberToTakeString)

	//make sure the requirement at least requires one class
	if numberToTake < 1 {
		c.JSON(200, gin.H{})
		return
	}

	programIDString := c.PostForm("programID")
	programID, _ := strconv.Atoi(programIDString)
	program := Program{}
	//make sure the program exists with this id
	db.Where("program_id = ?", uint64(programID)).First(&program)
	if program.ProgramID == 0 {
		c.JSON(200, gin.H{})
		return
	}

	newRequirment.NumberToTake = uint64(numberToTake)
	newRequirment.ProgramID = program.ProgramID

	//make sure this requirement doesn't already exist
	newRequirmentTest := ProgramRequirement{}
	db.Where(newRequirment).First(&newRequirmentTest)
	if newRequirmentTest.ProgramRequirementID != 0 {
		c.JSON(200, gin.H{"programAlreadyExists": "true"})
		return
	}

	db.Create(&newRequirment)
	db.Update()
	c.JSON(200, gin.H{})
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
	if programRequirement.ProgramRequirementID == 0 {
		c.JSON(200, gin.H{})
		return
	}

	courseCodeInt, _ := strconv.Atoi(courseCodeString)
	requirementCourse := RequirementCourse{}
	requirementCourse.CourseCode = uint64(courseCodeInt)
	requirementCourse.CollegeName = collegeName

	testReqCourse := RequirementCourse{}
	db.Where(requirementCourse).First(&testReqCourse)
	if testReqCourse.CourseCode == 0 {
		db.Create(&requirementCourse)
		db.Update()
		db.Where(requirementCourse).First(&testReqCourse)
	}
	requirementCourse = testReqCourse

	//make sure this requirement does not already exist
	requirementClasses := RequirementToRequirementCourse{}
	db.Where("program_requirement_id = ? and requirement_course_id = ?", programRequirement.ProgramRequirementID, requirementCourse.RequirementCourseID).Find(&requirementClasses)
	if requirementClasses.ProgramRequirementID != 0 {
		c.JSON(200, gin.H{})
		return
	}

	//make sure this requirement does not exist in the exlude catagory
	exludeCourse := RequirementToExcludeThisCourse{}
	db.Where("program_requirement_id = ? and requirement_course_id = ?", programRequirement.ProgramRequirementID, requirementCourse.RequirementCourseID).Find(&exludeCourse)
	if exludeCourse.ProgramRequirementID != 0 {
		c.JSON(200, gin.H{})
		return
	}

	requirementClasses.ProgramRequirementID = programRequirement.ProgramRequirementID
	requirementClasses.RequirementCourseID = requirementCourse.RequirementCourseID

	db.Create(&requirementClasses)

	c.JSON(200, gin.H{})
}

func addGreaterThanRequirementToProgram(c *gin.Context) {
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

	greaterThanRequirement := RequirementToRequirementGreaterThan{}
	greaterThanRequirement.CollegeName = c.PostForm("collegeName")
	courseCodeString := c.PostForm("courseCode")
	courseCode, _ := strconv.Atoi(courseCodeString)
	greaterThanRequirement.CourseCodeMinimum = uint64(courseCode)
	requirementIDString := c.PostForm("requirementID")
	requirementID, _ := strconv.Atoi(requirementIDString)
	greaterThanRequirement.ProgramRequirementID = uint64(requirementID)

	requirement := ProgramRequirement{}
	db.Where("program_requirement_id = ?", greaterThanRequirement.ProgramRequirementID).First(&requirement)
	//make sure this program exists
	if requirement.ProgramRequirementID == 0 {
		c.JSON(200, gin.H{})
		return
	}
	//make sure that this requirement does not already exist
	greaterThanRequirementTest := RequirementToRequirementGreaterThan{}
	db.Where(greaterThanRequirement).First(&greaterThanRequirementTest)
	if greaterThanRequirementTest.ProgramRequirementID != 0 {
		c.JSON(200, gin.H{})
		return
	}

	//now we add the requirement
	db.Create(&greaterThanRequirement)
	c.JSON(200, gin.H{})
}

func addCourseExclusionToProgram(c *gin.Context) {
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

	exclusionCourse := RequirementToExcludeThisCourse{}

	course := RequirementCourse{}
	course.CollegeName = c.PostForm("collegeName")
	courseCodeString := c.PostForm("courseCode")
	courseCode, _ := strconv.Atoi(courseCodeString)
	course.CourseCode = uint64(courseCode)

	//see if this course exists, if not: add it
	courseTest := RequirementCourse{}
	db.Where(course).First(&courseTest)
	if courseTest.RequirementCourseID == 0 {
		db.Create(&course)
		db.Update()
	}
	db.Where(course).First(&course)

	requirement := ProgramRequirement{}
	requirementIDString := c.PostForm("requirementID")
	requirementID, _ := strconv.Atoi(requirementIDString)

	//make sure the requirement exists
	db.Where("program_requirement_id = ?", uint64(requirementID)).First(&requirement)
	if requirement.ProgramRequirementID == 0 {
		c.JSON(200, gin.H{})
		return
	}

	exclusionCourse.RequirementCourseID = course.RequirementCourseID
	exclusionCourse.ProgramRequirementID = requirement.ProgramRequirementID

	//make sure this exclusion does not already exist
	exclusionCourseTest := RequirementToExcludeThisCourse{}
	db.Where(exclusionCourse).First(&exclusionCourseTest)
	if exclusionCourseTest.ProgramRequirementID != 0 {
		c.JSON(200, gin.H{})
		return
	}

	//make sure that this course is not enlisted as a specific course to meet the requirement
	requirementCourse := RequirementCourse{}
	db.Where(exclusionCourse).First(&requirementCourse)
	if requirementCourse.RequirementCourseID != 0 {
		c.JSON(200, gin.H{})
		return
	}

	db.Create(&exclusionCourse)
	db.Update()

	c.JSON(200, gin.H{})
}

func deleteGreaterThanRequirement(c *gin.Context) {
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
	courseCodeString := c.PostForm("courseCodeMinimum")
	requirementID, _ := strconv.Atoi(requirementIDString)
	requirementCourseCode, _ := strconv.Atoi(courseCodeString)

	greaterThanRequirement := RequirementToRequirementGreaterThan{}
	greaterThanRequirement.ProgramRequirementID = uint64(requirementID)
	greaterThanRequirement.CourseCodeMinimum = uint64(requirementCourseCode)
	greaterThanRequirement.CollegeName = c.PostForm("collegeName")

	db.Delete(&greaterThanRequirement)
	c.JSON(200, gin.H{})
}
