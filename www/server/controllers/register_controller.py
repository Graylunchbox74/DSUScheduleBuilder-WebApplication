from server import app
from server import controllers
from server.models.forms.registration_form import RegistrationForm
from server.facades import user_facade

import flask


@app.route('/register', methods=['GET', 'POST'])
def register():
    global_context = controllers.get_global_context_variables()

    register_form = RegistrationForm()
    if register_form.validate_on_submit():

        # Tell the DB About this
        (success, _, err) = user_facade.register_user(
            email=register_form.email.data,
            password=register_form.password.data,
            firstName=register_form.first_name.data,
            lastName=register_form.last_name.data,
        )

        if err or not success:
            flask.flash(f"An error occured creating your account. Please try again later.", "danger")
        else:
            flask.flash(f"Successfully created your account. You may now log in.", "success")
            return flask.redirect(flask.url_for('login'))

    context = {
        "globals": global_context,
        "form": register_form
    }

    return flask.render_template('register.html', **context)
