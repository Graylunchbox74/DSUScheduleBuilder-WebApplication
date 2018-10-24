$(function () {

    let course_facade = modulejs.require("course_facade")

    let $college_name = $("#collegeName")
    let $course_number = $("#courseNumber")
    let $teacher_name = $("#teacherName")
    let $location = $("#location")
    let $semester = $("#semester")
    let $course_name = $("#courseName")

    function getCoursesFromSearchTerms() {
        let college_name = $college_name.val()
        let course_code = $course_number.val()
        let teacher_name = $teacher_name.val()
        let location = $location.val()
        let semester = $semester.val()
        let course_name = $course_name.val()

        return course_facade.searchCourses({
            college_name,
            course_code,
            teacher_name,

            //Leaving out location and semester for now
            // location,
            //semester,

            course_name,
        })
    }

    function renderCourse(course) {
        return `
            <div class="d-flex flex-wrap py-3 my-2 secondary-background">
                <div class="d-flex w-100">
                    <div class="col-md-8"><strong>${course.college_name}-${course.course_code}</strong>: ${course.course_name}</div>
                    <div class="col-md-2 text-right">${course.days_of_week}</div>
                    <div class="col-md-2 text-right">${course.start_time} - ${course.end_time}</div>
                </div>
                <div class="d-flex flex-wrap mt-4 w-100">
                    <div class="col-md-8 col-sm-12">${course.teacher}</div>
                    <div class="col-md-4 col-sm-12 text-right">
                        <button class="btn btn-primary" type="button" data-toggle="collapse" data-target="#course-description-${course.course_id}">More details</button>
                    </div>
                </div>
                <div class="mt-4 w-100 collapse" id="course-description-${course.course_id}" data-parent="#search-container">
                    <div class="card card-body card-dark">
                        <h3>Description:</h3>
                        <p>${course.description}</p>
                        <div class="text-right">
                            <button class="btn btn-primary" onclick="enrollInCourse(${course.course_id})">Enroll</button>
                        </div>
                    </div>
                </div>
            </div>
        `
    }

    window.enrollInCourse = function (course_id) {
        course_facade.enrollInCourse(course_id).then(_ => {
            window.location.reload()
        })
    }

    let $searchContainer = $("#search-container")
    function updateResults(courses) {
        if (courses.length == 0) {
            $searchContainer.html(`
                <div class="p-4 text-center text-muted">No classes found with specified criteria.</div>
            `)
        } else {

            $searchContainer.html(function () {
                let html_string = ""

                for (let course of courses) {
                    html_string += renderCourse(course)
                }

                return html_string
            })
        }
    }


    $("#searchButton").on("click", function () {
        getCoursesFromSearchTerms().then(function (courses) {
            console.log("%c Courses found:", "color: green; font-weight: 800")
            console.table(courses)
            updateResults(courses)
        })
    })
})