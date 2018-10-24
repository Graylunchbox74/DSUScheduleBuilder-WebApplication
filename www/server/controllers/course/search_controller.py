from server import controllers
from server.models.user_model import UserModel
from server.facades import course_facade
from server.models.utils import facade_result_codes as FRC
import flask


class CourseSearchController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        context = {
            "selected_tab": "courses"
        }

        return flask.render_template('courses/search.html', **context)

    def post(self):
        user = UserModel.create(**flask.session['user'])

        search_data = flask.request.get_json()

        (status, courses) = course_facade.search_for_courses(user.token, search_data)

        if status == FRC.NOT_AUTHENTICATED:
            return self.handle_not_authorized()
        if status != FRC.SUCCESS:
            return flask.jsonify({})

        course_list = list(map(lambda c: c.to_json(), courses))

        return flask.jsonify(course_list)


class CourseEnrollingController(controllers.BaseController):
    decorators = [controllers.login_required]

    def post(self):
        user = UserModel.create(**flask.session['user'])

        course_id = flask.request.get_json().get('course_id')

        if course_id == None:
            return flask.jsonify({})

        (status, conflicts) = course_facade.enroll_in_course(user.token, course_id)

        if status == FRC.NOT_AUTHENTICATED:
            return self.handle_not_authorized()
        if status != FRC.SUCCESS:
            if conflicts: 
                flask.flash(f"The course you are trying to register for conflicts with a course you are already in.", "warning")
            else:
                flask.flash(f"There was an error registering for the course. Please try again later.", "danger")
        else:
            flask.flash(f"Successfully registered.", "success")

        return flask.jsonify({})
        