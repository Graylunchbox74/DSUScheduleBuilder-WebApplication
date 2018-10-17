from server import app
from server import controllers
import flask


class HomeController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        context = {
            "selected_tab": "home",
        }

        return flask.render_template('index.html', **context)
