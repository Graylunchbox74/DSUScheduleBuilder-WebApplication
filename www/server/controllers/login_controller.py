from server import app
from server import controllers
import flask


@app.route('/login')
def login():
    global_context = controllers.get_global_context_variables()

    context = {
        "globals": global_context
    }

    return flask.render_template('login.html', **context)
