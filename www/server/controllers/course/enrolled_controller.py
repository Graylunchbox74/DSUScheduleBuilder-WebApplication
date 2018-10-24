from server import controllers
from server.facades import course_facade
from server.models.user_model import UserModel
from server.models.utils import facade_result_codes as FRC
import flask


class CoursesEnrolledController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        user = UserModel.create(**flask.session['user'])

        context = {
            "selected_tab": "courses",
            "courses": [],
        }

        (status, course_list) = course_facade.get_enrolled_courses(user.token)

        if status == FRC.NOT_AUTHENTICATED:
            return self.handle_not_authorized()

        if status != FRC.SUCCESS:
            flask.flash(f"Error loading your currently enrolled courses.", "danger")
        else:
            context['courses'] = course_list

        return flask.render_template('courses/enrolled.html', **context)


class CourseDroppingController(controllers.BaseController):
    decorators = [controllers.login_required]

    def post(self):
        user = UserModel.create(**flask.session['user'])

        course_id = flask.request.get_json()['course_id']

        auth = course_facade.drop_course(user.token, course_id)

        if not auth:
            return self.handle_not_authorized()
        else:
            flask.flash(f"Successfully dropped course.", "success")
            return flask.jsonify({})

        return
