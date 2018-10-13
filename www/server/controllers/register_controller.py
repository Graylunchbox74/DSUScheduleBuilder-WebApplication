from server import app
from server import controllers
from server.models.registration_form import RegistrationForm

import requests
import flask


@app.route('/register', methods=['GET', 'POST'])
def register():
    global_context = controllers.get_global_context_variables()

    register_form = RegistrationForm()
    if register_form.validate_on_submit():
        # Tell the DB About this
        data = {
            'email': register_form.email.data,
            'password': register_form.password.data,
            'firstName': register_form.first_name.data,
            'lastName': register_form.last_name.data,
        }

        response = requests.post(f"{app.config['API_ENDPOINT']}/user/newUser", data=data)
        response = response.json()
        if response['errorMsg'] == "":
            return flask.redirect(flask.url_for('login'))
        else:
            flask.flash(f"An error occured creating your account. Please try again later.")

    context = {
        "globals": global_context,
        "form": register_form
    }

    return flask.render_template('register.html', **context)
