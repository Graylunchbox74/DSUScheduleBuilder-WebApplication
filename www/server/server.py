import flask
from flask import Flask

app = Flask(__name__, instance_relative_config=True)
app.config.from_object('config.default')
app.config.from_envvar('APP_CONFIG')


@app.route("/")
def index():
    return "Home page is the best page"


def start_server():
    app.run(host=app.config['HOST'], port=app.config['PORT'])


if __name__ == '__main__':
    start_server()
