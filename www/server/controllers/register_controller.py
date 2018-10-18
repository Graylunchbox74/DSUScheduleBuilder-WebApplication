from server import controllers
from server.models.forms.registration_form import RegistrationForm
from server.models.utils import facade_result_codes
from server.facades import user_facade

import flask


class RegisterController(controllers.BaseController):
    decorators = [controllers.logout_required]

    def get(self):
        register_form = RegistrationForm()

        context = {
            "form": register_form
        }

        return flask.render_template('register.html', **context)

    def post(self):
        register_form = RegistrationForm()
        if register_form.validate_on_submit():

            # Tell the DB About this
            (code, _) = user_facade.register_user(
                email=register_form.email.data,
                password=register_form.password.data,
                first_name=register_form.first_name.data,
                last_name=register_form.last_name.data,
            )

            if code == facade_result_codes.SUCCESS:
                flask.flash(f"Successfully created your account. You may now log in.", "success")
                return flask.redirect(flask.url_for('login'))
            else:
                flask.flash(f"An error occured creating your account. Please try again later.", "danger")

        context = {
            "form": register_form
        }

        return flask.render_template('register.html', **context)
