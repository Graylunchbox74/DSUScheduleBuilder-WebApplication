package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var mg mailgun.Mailgun

func init() {
	var emailData struct {
		Private, Public string
	}

	emailDataFile, err := os.Open("emailData.json")

	if err != nil {
		panic("Could not get email data.")
	}

	err = json.NewDecoder(emailDataFile).Decode(&emailData)

	if err != nil {
		panic("Could not parse the email data from the json file.")
	}

	mg = mailgun.NewMailgun("mail.therileyjohnson.com", emailData.Private, emailData.Public)
}

func postAddData(c *gin.Context) {
	// hypothetically add data here
}

func getUUID() string {
	var err error
	var uid uuid.UUID
	for uid, err = uuid.NewV4(); err != nil; {
		uid, err = uuid.NewV4()
	}
	return uid.String()
}

func getProgramRequirements(c *gin.Context) {
	program := Program{}
	programString := c.Request.URL.Query()["programID"][0]
	tmp, _ := strconv.Atoi(programString)
	program.ProgramID = uint64(tmp)

	//make sure the program exists
	db.Where(program).Find(&program)
	if program.ProgramID == 0 {
		c.JSON(200, gin.H{})
	}

	programReq := []ProgramRequirement{}

	db.Where("program_id = ?", program.ProgramID).Find(&programReq)

	c.JSON(200, programReq)
}

func searchPrograms(c *gin.Context) {
	program := Program{}
	program.Program = c.Request.URL.Query()["programName"][0]
	strTmp := c.Request.URL.Query()["catalogYear"][0]
	tmp, _ := strconv.Atoi(strTmp)
	program.CatalogYear = uint64(tmp)

	returnPrograms := []Program{}
	db.Where("program like ?", "%"+program.Program+"%").Find(&returnPrograms)

	c.JSON(200, returnPrograms)
}

func searchForCourse(c *gin.Context) {
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

func getAllSemesters(c *gin.Context) {
	semesters := []Semester{}
	db.Find(&semesters)
	c.JSON(200, semesters)
}

func getAllLocations(c *gin.Context) {
	locations := []Location{}
	db.Find(&locations)
	c.JSON(200, locations)
}

func getProgramCourseExclusionsRequirementField(c *gin.Context) {
	requirementIDString := c.Request.URL.Query()["requirementID"][0]
	requirementID, _ := strconv.Atoi(requirementIDString)
	exclusionCourses := []RequirementToExcludeThisCourse{}
	db.Where("program_requirement_id = ?", uint64(requirementID)).Find(&exclusionCourses)

	courses := []RequirementCourse{}
	courseTmp := RequirementCourse{}
	for _, currentRecord := range exclusionCourses {
		courseTmp = RequirementCourse{}
		db.Where("requirement_course_id = ?", currentRecord.RequirementCourseID).First(&courseTmp)
		if courseTmp.RequirementCourseID != 0 {
			courses = append(courses, courseTmp)
		}
	}

	c.JSON(200, courses)
}

func getProgramCourseSpecificRequirementField(c *gin.Context) {
	requirementIDString := c.Request.URL.Query()["requirementID"][0]
	requirementID, _ := strconv.Atoi(requirementIDString)
	courseRequirementList := []RequirementToRequirementCourse{}
	db.Where("program_requirement_id = ?", uint64(requirementID)).Find(&courseRequirementList)

	courses := []RequirementCourse{}
	courseTmp := RequirementCourse{}
	for _, currentRecord := range courseRequirementList {
		courseTmp = RequirementCourse{}
		db.Where("requirement_course_id = ?", currentRecord.RequirementCourseID).First(&courseTmp)
		if courseTmp.RequirementCourseID != 0 {
			courses = append(courses, courseTmp)
		}
	}

	c.JSON(200, courses)
}

func getProgramGreaterRequirementField(c *gin.Context) {
	requirementIDString := c.Request.URL.Query()["requirementID"][0]
	requirementID, _ := strconv.Atoi(requirementIDString)
	courseRequirementList := []RequirementToRequirementGreaterThan{}
	db.Where("program_requirement_id = ?", uint64(requirementID)).Find(&courseRequirementList)
	c.JSON(200, courseRequirementList)
}

func postCreateAccount(context *gin.Context) {
	defer context.Request.Body.Close()

	pendingUser := PendingAccount{}

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&pendingUser); err != nil {
		errorResponse(
			context,
			"Error recieving registration information, please try again.",
			err.Error(),
		)
		return
	}

	if pendingUser.FirstName == "" || pendingUser.LastName == "" || pendingUser.Password == "" {
		errorResponse(
			context,
			"None of the fields can be blank.",
			"Blank field.",
		)
		return
	}

	matchEmail := false
	emailDomains := []string{"trojans.dsu.edu", "pluto.dsu.edu", "dsu.edu"}
	pendingUser.Email = strings.ToLower(pendingUser.Email)

	for _, emailRegex := range emailDomains {
		r := regexp.MustCompile(fmt.Sprintf(`^[A-Za-z0-9][A-Za-z0-9_\+\.]*@%s$`, emailRegex))
		if r.Match([]byte(pendingUser.Email)) {
			matchEmail = true
		}
	}

	if !matchEmail {
		errorResponse(
			context,
			"Email is invalid, must be a valid email with either a trojans.dsu.edu, pluto.dsu.edu, or dsu.edu domain.",
			"Email no match.",
		)
		return
	}

	if len(pendingUser.Email) > 150 {
		context.JSON(
			400,
			gin.H{
				"error": true,
				"msg":   "Email is too long (over 150 characters).",
			},
		)
		return
	}

	pendingUser.UUID = getUUID()
	// pendingUser.Password = hashPassword(pendingUser.Password)

	pendingAccount, err := CreatePendingAccount(db, &pendingUser)

	if err != nil {
		errorResponse(
			context,
			"Error creating account, try again.",
			err.Error(),
		)

		return
	}

	// Reference https://github.com/the-rileyj/DSU_Chess/blob/master/chess.go
	_, _, err = mg.Send(mailgun.NewMessage("robot@mail.therileyjohnson.com", "Registration", fmt.Sprintf("Click http://localhost:8079/api/user/confirm/%s to confirm your email!", pendingUser.UUID), pendingUser.Email))

	context.JSON(
		200,
		gin.H{
			"data":  pendingAccount,
			"error": false,
		},
	)
}

func hashPassword(password string) string {
	var err error
	var hashedPassword []byte

	for hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14); err != nil; {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	}

	return string(hashedPassword)
}

func getConfirmAccount(context *gin.Context) {
	uuid := context.Param("uuid")

	pendingUser := PendingAccount{}

	if err := db.Where(PendingAccount{UUID: uuid}).First(&pendingUser).Error; err != nil {
		errorResponse(
			context,
			"Sorry, an error occurred when confirming your account, please check that you have the correct link and try again.",
			err.Error(),
		)
	}

	newUser := Student{
		Email:     pendingUser.Email,
		FirstName: pendingUser.FirstName,
		LastName:  pendingUser.LastName,
		Password:  pendingUser.Password,
	}

	account, err := CreateAccount(db, &newUser)

	if err != nil {
		errorResponse(
			context,
			"Sorry, an error occurred when confirming your account, please check that you have the correct link and try again.",
			err.Error(),
		)
		return
	}

	db.Delete(pendingUser)

	context.JSON(
		200,
		gin.H{
			"data":  account,
			"error": false,
		},
	)
}

func errorResponse(context *gin.Context, msg, debug string) {
	context.JSON(
		400,
		gin.H{
			"error": true,
			"msg":   msg,
			"debug": debug,
		},
	)
}
