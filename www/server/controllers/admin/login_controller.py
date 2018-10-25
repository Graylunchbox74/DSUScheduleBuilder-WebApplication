from server import app
from server.controllers import admin
from server.models.forms.login_form import LoginForm
from server.models.user_model import UserModel

from server.models.utils import facade_result_codes

from server.facades import admin_facade

import flask


class AdminLoginController(admin.BaseController):
    decorator = [admin.admin_login_required]

    def get(self):
        login_form = LoginForm()

        context = {
            "form": login_form
        }

        return flask.render_template('admin/login.html', **context)

    def post(self):
        login_form = LoginForm()
        if login_form.validate_on_submit():

            # Validate User from DB
            (code, user) = admin_facade.login_admin(login_form.email.data, login_form.password.data)

            if code == facade_result_codes.SUCCESS:
                flask.session['admin'] = user.to_json()
                return flask.redirect(flask.url_for('admin_index'))
            elif code == facade_result_codes.NOT_AUTHENTICATED:
                flask.flash(f"Invalid email or password. Please try again.", 'danger')
            else:
                flask.flash(f"An error occurred when logging in. Please try again later.", "danger")

        context = {
            "form": login_form
        }

        return flask.render_template('admin/login.html', **context)


class AdminLogoutController(admin.BaseController):
    decorator = [admin.admin_login_required]

    def get(self):
        if flask.session.get('admin') is not None:
            admin = UserModel.create(**flask.session['admin'])

            admin_facade.logout_admin(admin.token)

            del flask.session['admin']
            flask.flash("You were logged out.", "info")

        return flask.redirect(flask.url_for('admin_login'))
