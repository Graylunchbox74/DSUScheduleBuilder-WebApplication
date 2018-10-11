from server import app


@app.route("/")
def index():
    return "Home page is the best page"
