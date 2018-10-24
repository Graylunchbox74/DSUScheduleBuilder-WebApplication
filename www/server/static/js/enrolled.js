$(function () {
    let course_facade = modulejs.require("course_facade")

    window.dropCourse = function (course_id) {
        course_facade.dropCourse(course_id).then(_ => {
            window.location.reload()
        })
    }
})