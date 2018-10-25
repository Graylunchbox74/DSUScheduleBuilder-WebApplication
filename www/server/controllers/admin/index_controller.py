from server import app
from server.controllers import admin
import flask


class AdminIndexController(admin.BaseController):
    decorators = [admin.admin_login_required]

    def get(self):
        context = {
            "selected_tab": "home",
        }

        return flask.render_template('admin/index.html', **context)
