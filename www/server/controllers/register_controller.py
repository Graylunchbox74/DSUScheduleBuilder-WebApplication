from server import app
from server import controllers
from server.models.registration_form import RegistrationForm
import flask


@app.route('/register', methods=['GET', 'POST'])
def register():
    global_context = controllers.get_global_context_variables()

    register_form = RegistrationForm()
    if register_form.validate_on_submit():
        # Tell the DB About this

        return flask.redirect(flask.url_for('login'))

    context = {
        "globals": global_context,
        "form": register_form
    }

    return flask.render_template('register.html', **context)
