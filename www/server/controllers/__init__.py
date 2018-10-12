from functools import wraps
from flask import session, redirect, url_for

from server import app


def get_global_context_variables():
    return {
        "brand_name": app.config['BRAND_NAME'],
    }


# Route helpers
def login_required(func):
    @wraps(func)
    def decorated_function(*args, **kwargs):
        if session.get('user') is None:
            return redirect(url_for('login'))
        return func(*args, **kwargs)

    return decorated_function
