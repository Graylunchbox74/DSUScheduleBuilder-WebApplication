from server import app
from server import controllers
import flask


@app.route('/')
@controllers.login_required
def index():
    views = flask.session.get('views', 0)
    views += 1
    flask.session['views'] = views

    global_context = controllers.get_global_context_variables()

    context = {
        "globals": global_context,
        "views": views
    }

    return flask.render_template('index.html', **context)
