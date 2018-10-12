from server import app
from server import controllers
from server.models.login_form import LoginForm
import flask


@app.route('/login', methods=['GET', 'POST'])
def login():
    global_context = controllers.get_global_context_variables()

    login_form = LoginForm()
    if login_form.validate_on_submit():
        flask.session['user'] = True
        return flask.redirect(flask.url_for('index'))

    context = {
        "globals": global_context,
        "form": login_form
    }

    return flask.render_template('login.html', **context)
