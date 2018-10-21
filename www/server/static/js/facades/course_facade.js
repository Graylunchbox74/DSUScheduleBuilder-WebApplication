modulejs.define('course_facade', ['request_manager'], function (request) {

    function addParamIfNotEmpty(params, param_key, param_val) {
        if (param_val) {
            params[param_key] = param_val
        }
    }

    function searchCourses(search_terms) {
        let search_params = {};
        addParamIfNotEmpty(search_params, "college_name", search_terms.college_name)
        addParamIfNotEmpty(search_params, "course_code", search_terms.course_code)
        addParamIfNotEmpty(search_params, "teacher", search_terms.teacher_name)
        addParamIfNotEmpty(search_params, "location", search_terms.location)
        addParamIfNotEmpty(search_params, "semester", search_terms.semester)
        addParamIfNotEmpty(search_params, "course_name", search_terms.course_name)

        console.log("%c Searching with terms:", "color: blue; font-weight: 800")
        console.table(search_params);

        return request
            .post_json("/courses/search", search_params)
    }

    return {
        searchCourses: searchCourses,
    }
})