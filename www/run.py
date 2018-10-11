from server import app


def start_server():
    app.run(host=app.config['HOST'], port=app.config['PORT'])


if __name__ == '__main__':
    start_server()
