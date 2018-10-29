from server import controllers
from server.models.user_model import UserModel
from server.models.utils import facade_result_codes as FRC
from server.facades import course_facade

import flask
import re


class ScheduleController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        courses = []

        user = UserModel.create(**flask.session['user'])
        
        (status, cs) = course_facade.get_enrolled_courses(user.token)

        if status == FRC.NOT_AUTHENTICATED:
            return self.handle_not_authorized()

        if status != FRC.SUCCESS:
            flask.flash(f"There was an error loading your enrolled courses.", "danger")

        for c in cs:
            days = []
            if "Mon" in c.days_of_week:
                days.append(0)
            if "Tues" in c.days_of_week:
                days.append(1)
            if "Wed" in c.days_of_week:
                days.append(2)
            if "Thur" in c.days_of_week:
                days.append(3)
            if "Fri" in c.days_of_week:
                days.append(4)


            start_m = re.match("^([0-9]+)\:([0-9]+)", c.start_time)
            start_off = int(start_m.group(1)) + int(start_m.group(2)) / 60.0

            end_m = re.match("^([0-9]+)\:([0-9]+)", c.end_time)
            end_off = int(end_m.group(1)) + int(end_m.group(2)) / 60.0

            for d in days:
                courses.append({
                    "first_row": c.course_name,
                    "second_row": c.start_time + " - " + c.end_time,
                    "start_time": start_off,
                    "end_time": end_off,
                    "end_off": end_off,
                    "day_of_week": d,
                    "color": "white",
                })

        context = {
            "selected_tab": "schedule",

            "courses": courses
        }

        return flask.render_template('schedule.html', **context)
