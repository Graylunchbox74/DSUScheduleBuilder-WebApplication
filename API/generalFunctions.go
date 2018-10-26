package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProgramRequirements(c *gin.Context) {
	program := Program{}
	programString := c.Request.URL.Query()["programID"][0]
	tmp, _ := strconv.Atoi(programString)
	program.ProgramID = uint64(tmp)

	db.Where(program).Find(&program)
	if program.ProgramID == 0 {
		c.JSON(200, program)
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
