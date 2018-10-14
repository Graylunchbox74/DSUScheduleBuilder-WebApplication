import flask
from flask import Flask

# Create app object
app = Flask(
    __name__,
    template_folder='views',
    instance_relative_config=True
)

# Load configs
app.config.from_object('config.default')
app.config.from_envvar('APP_CONFIG')

from server.controllers import home_controller
from server.controllers import login_controller
from server.controllers import register_controller
from server.controllers import schedule_controller
