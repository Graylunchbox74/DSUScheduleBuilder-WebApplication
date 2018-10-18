from server import controllers
import flask


class CourseSearchController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        context = {
            "selected_tab": "courses"
        }

        return flask.render_template('courses/search.html', **context)
