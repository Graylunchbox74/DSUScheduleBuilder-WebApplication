from server import app
from server import controllers
from server.models.user_model import UserModel
import flask


class HomeController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        user = UserModel()
        user.from_json(flask.session['user'])

        context = {
            "selected_tab": "home",
            "user": user,
        }

        return flask.render_template('index.html', **context)
