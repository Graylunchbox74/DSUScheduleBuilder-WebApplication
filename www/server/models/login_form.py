from flask_wtf import FlaskForm
from wtforms import StringField, PasswordField
from wtforms.validators import DataRequired, Length


class LoginForm(FlaskForm):
    email = StringField(
        'Email',
        validators=[DataRequired(), Length(min=1, max=100)]
    )

    password = PasswordField(
        'Password',
        validators=[DataRequired()]
    )
