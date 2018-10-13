from functools import wraps
from flask import session, redirect, url_for
from server import app

from server.models.user_model import UserModel


def get_global_context_variables():
    return {
        "brand_name": app.config['BRAND_NAME'],

        # This allows the Model functions to be accessable from the view
        "models": {
            "user": UserModel
        }
    }


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
