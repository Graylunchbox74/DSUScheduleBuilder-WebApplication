from server import app
from server import controllers
import flask


@app.route('/schedule')
@controllers.login_required
def schedule():
    global_context = controllers.get_global_context_variables()

    courses = [
        {
            "first_row": "Test Class",
            "second_row": "8:00AM - 8:50AM",
            "start_time": 8,
            "end_time": 8 + (50.0 / 60.0),
            "day_of_week": 0,
            "color": "#b34",
        },
        {
            "first_row": "Test Class",
            "second_row": "8:00AM - 8:50AM",
            "start_time": 8,
            "end_time": 8 + (50.0 / 60.0),
            "day_of_week": 2,
            "color": "#b34",
        },
        {
            "first_row": "Test Class",
            "second_row": "8:00AM - 8:50AM",
            "start_time": 8,
            "end_time": 8 + (50.0 / 60.0),
            "day_of_week": 4,
            "color": "#b34",
        },
        {
            "first_row": "Second Class",
            "second_row": "1:00PM - 2:15PM",
            "start_time": 13,
            "end_time": 14 + (15.0 / 60.0),
            "day_of_week": 1,
            "color": "#34b",
        },
        {
            "first_row": "Second Class",
            "second_row": "1:00PM - 2:15PM",
            "start_time": 13,
            "end_time": 14 + (15.0 / 60.0),
            "day_of_week": 3,
            "color": "#34b",
        },
    ]

    context = {
        "globals": global_context,
        "selected_tab": "schedule",

        "courses": courses
    }

    return flask.render_template('schedule.html', **context)
