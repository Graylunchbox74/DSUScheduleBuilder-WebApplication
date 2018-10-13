from server import app
from server import controllers
from server.models.login_form import LoginForm
from server.models.user_model import UserModel

import requests
import flask


@app.route('/login', methods=['GET', 'POST'])
@controllers.logout_required
def login():
    global_context = controllers.get_global_context_variables()

    login_form = LoginForm()
    if login_form.validate_on_submit():
        # Validate User from DB

        # Load user info from DB

        # For now, assume creds are correct
        user = UserModel.create(
            email=login_form.email.data,
            first_name="UNSET",
            last_name="UNSET",
        )

        flask.session['user'] = user
        flask.session['views'] = 0
        return flask.redirect(flask.url_for('index'))

    context = {
        "globals": global_context,
        "form": login_form
    }

    return flask.render_template('login.html', **context)


@app.route('/logout')
@controllers.login_required
def logout():
    if flask.session.get('user') is not None:
        del flask.session['user']

    return flask.redirect(flask.url_for('index'))
