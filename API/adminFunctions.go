package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func addProgram(c *gin.Context) {
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
	courseIDString := c.PostForm("courseID")
	tmp, _ := strconv.Atoi(courseIDString)
	var course Course
	db.Where("course_id = ?", uint64(tmp)).Delete(&course)

	c.JSON(200, gin.H{})
}

func addRequirementToProgram(c *gin.Context) {
	programID := c.PostForm("programID")
	collegeName := c.PostForm("collegeName")
	courseCodeString := c.PostForm("courseCode")

	courseCodeInt, _ := strconv.Atoi(courseCodeString)
	requirementCourse := RequirementCourse{}
	requirementCourse.CourseCode = uint64(courseCodeInt)
	requirementCourse.CollegeName = collegeName

	program := Program{}

	testReqCourse := RequirementCourse{}
	db.Where(requirementCourse).First(&testReqCourse)
	if testReqCourse.CourseCode == 0 {
		db.Create(&requirementCourse)
		db.Where(requirementCourse).First(&testReqCourse)
	}
	requirementCourse = testReqCourse

	db.Where("program_id = ?", programID).First(&program)

}
