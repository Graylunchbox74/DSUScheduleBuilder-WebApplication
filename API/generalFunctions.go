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
