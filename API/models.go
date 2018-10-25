package main

import "time"

type Admin struct {
	ID        uint64 `gorm:"primary_key"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Password  string
}

type AdminSessionToken struct {
	AdminID     uint64
	Token       string `gorm:"unique:"`
	TimeUpdated time.time
}

type SessionToken struct {
	StudentID   uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
	Token       string `gorm:"unique"`
	TimeUpdated time.Time
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
	RequirementName      string
	ProgramID            uint64 `gorm:"foreignkey:ProgramID;association_foreignkey:ProgramID"`
	ProgramRequirementID uint64 `gorm:"primary_key"`
	NumberToTake         uint64
}

type StudentProgramRequirement struct {
	RequirementName     string
	ProgramID           uint64 `gorm:"foreignkey:ProgramID;association_foreignkey:ProgramID"`
	StudentID           uint64 `gorm:"foreignkey:StudentID;association_foreignkey:StudentID"`
	NumberToTake        uint64
	RequirementCourseID uint64
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
	Semester    string
}

type PreviouslyEnrolled struct {
	StudentID   uint64 `json:"studentID"`
	Credits     uint64 `json:"credits"`
	CollegeName string `json:"collegeName"`
	CourseCode  uint64 `json:"courseCode"`
}

type ProgramToRequirementCourse struct {
	ProgramID           uint64
	RequirementCourseID uint64
}

type RequirementCourse struct {
	RequirementCourseID uint64 `gorm:"primary_key"`
	CollegeName         string
	CourseCode          uint64
}

type returnStudent struct {
	StudentID uint64 `json:"studentID"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
}

type returnRequirement struct {
	isCompleted     bool
	requirementName string
	courses         []string
}
