from server import app
from server import controllers
from server.models.login_form import LoginForm

import requests
import flask


@app.route('/login', methods=['GET', 'POST'])
def login():
    global_context = controllers.get_global_context_variables()

    login_form = LoginForm()
    if login_form.validate_on_submit():
        flask.session['user'] = True
        flask.session['views'] = 0
        return flask.redirect(flask.url_for('index'))

    context = {
        "globals": global_context,
        "form": login_form
    }

    return flask.render_template('login.html', **context)


@app.route('/logout')
def logout():
    del flask.session['user']

    return flask.redirect(flask.url_for('index'))
