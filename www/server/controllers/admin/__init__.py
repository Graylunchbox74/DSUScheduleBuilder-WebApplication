from functools import wraps
from flask import session, redirect, url_for
from server import app

from server.models.user_model import UserModel
from server.facades import admin_facade

import flask
from flask.views import MethodView


class BaseController(MethodView):
    def handle_not_authorized(self):
        """
        Handles what happens when an admin is no longer authorized.
        Be sure to RETURN this function call.
        """

        flask.flash(f"Your session has expired. Please log in again.", "danger")

        # Log them out, which will require them to log back in
        return flask.redirect(flask.url_for('logout'))


# Route helpers
def admin_login_required(func):
    @wraps(func)
    def decorated_function(*args, **kwargs):
        logged_in = True

        if session.get('admin') is None:
            logged_in = False
        else:
            admin = UserModel.create(**session.get('admin'))
            if not admin_facade.check_admin_token(admin.token):
                logged_in = False

        if not logged_in:
            session['admin'] = None

            flask.flash(f"Your session has expired. Please log in again.", "danger")
            return redirect(url_for('admin_login'))
        return func(*args, **kwargs)

    return decorated_function


def admin_logout_required(func):
    @wraps(func)
    def decorated_function(*args, **kwargs):
        if session.get('admin') is not None:
            return redirect(url_for('admin_index'))
        return func(*args, **kwargs)

    return decorated_function
