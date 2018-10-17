from functools import wraps
from flask import session, redirect, url_for
from server import app

from server.models.user_model import UserModel

import flask
from flask.views import MethodView


class BaseController(MethodView):
    def handle_not_authorized(self):
        """
        Handles what happens when a user is no longer authorized.
        Be sure to RETURN this function call.
        """

        flask.flash(f"Your session has expired. Please log in again.")

        # Log them out, which will require them to log back in
        return flask.redirect(flask.url_for('logout'))


# Injects all global variables into the scope of every template
@app.context_processor
def inject_globals():
    return dict(globals={
        "brand_name": app.config['BRAND_NAME'],

        # This allows the Model functions to be accessable from the view
        "models": {
            "user": UserModel
        }
    })


# Route helpers
def login_required(func):
    @wraps(func)
    def decorated_function(*args, **kwargs):
        if session.get('user') is None:
            return redirect(url_for('login'))
        return func(*args, **kwargs)

    return decorated_function


def logout_required(func):
    @wraps(func)
    def decorated_function(*args, **kwargs):
        if session.get('user') is not None:
            return redirect(url_for('index'))
        return func(*args, **kwargs)

    return decorated_function
