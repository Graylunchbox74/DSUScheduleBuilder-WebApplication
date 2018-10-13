from flask_wtf import FlaskForm
from wtforms import StringField, PasswordField, SubmitField
from wtforms.validators import DataRequired, Length, EqualTo


class RegistrationForm(FlaskForm):
    email = StringField(
        'Email',
        validators=[DataRequired(), Length(min=1, max=100)]
    )

    first_name = StringField(
        'First name',
        validators=[DataRequired(), Length(min=1, max=25)]
    )

    last_name = StringField(
        'Last name',
        validators=[DataRequired(), Length(min=1, max=25)]
    )

    password = PasswordField(
        'Password',
        validators=[DataRequired()]
    )

    confirm_password = PasswordField(
        'Confirm password',
        validators=[DataRequired(), EqualTo('password')]
    )

    submit = SubmitField('Sign up')
