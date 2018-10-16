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
        data = {
            "email": login_form.email.data,
            "password": login_form.password.data,
        }

        try:
            response = requests.post(f"{app.config['API_ENDPOINT']}/user/validateUser", data=data)
            response = response.json()
            if response['StudentID'] != 0:
                user = UserModel.create(
                    id=response['StudentID'],
                    email=response['Email'],
                    first_name=response['firstName'],
                    last_name=response['lastName']
                )

                flask.session['user'] = user
                flask.session['views'] = 0
                return flask.redirect(flask.url_for('index'))
            else:
                flask.flash(f"Invalid email or password. Please try again.", 'danger')
        except:
            flask.flash(f"An error occurred when logging in. Please try again later.", "danger")

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
        flask.flash("Successfully logged out.", "success")

    return flask.redirect(flask.url_for('index'))
