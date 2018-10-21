from server import controllers
from server.models.forms.login_form import LoginForm
from server.models.utils import facade_result_codes
from server.facades import user_facade

import flask


class LoginController(controllers.BaseController):
    decorators = [controllers.logout_required]

    def get(self):
        login_form = LoginForm()

        context = {
            "form": login_form
        }

        return flask.render_template('login.html', **context)

    def post(self):
        login_form = LoginForm()
        if login_form.validate_on_submit():

            # Validate User from DB
            (code, user) = user_facade.validate_user(login_form.email.data, login_form.password.data)

            if code == facade_result_codes.SUCCESS:
                flask.session['user'] = user.to_json()
                flask.session['views'] = 0
                return flask.redirect(flask.url_for('index'))
            elif code == facade_result_codes.NOT_AUTHENTICATED:
                flask.flash(f"Invalid email or password. Please try again.", 'danger')
            else:
                flask.flash(f"An error occurred when logging in. Please try again later.", "danger")

        context = {
            "form": login_form
        }

        return flask.render_template('login.html', **context)


class LogoutController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        if flask.session.get('user') is not None:
            del flask.session['user']
            flask.flash("You were logged out.", "info")

        return flask.redirect(flask.url_for('index'))
