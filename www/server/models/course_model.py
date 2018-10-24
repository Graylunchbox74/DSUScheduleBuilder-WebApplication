from server.models.base_model import BaseModel


class CourseModel(BaseModel):
    def create(**kwargs):
        """Creates a course json object from the following fields:
            - course_id
            - course_code
            - course_name
            - credits
            - days_of_week
            - start_time
            - end_time
            - start_date
            - end_date
            - college_name
            - location
            - teacher
        """
        course = CourseModel()
        course.from_json(kwargs)
        return course

    def to_str(data):
        return f"CourseModel({data['course_id']}, {data['course_code']}, {data['course_name']})"

    def to_json(self):
        return {
            "course_id": self.course_id,
            "course_code": self.course_code,
            "course_name": self.course_name,
            "credits": self.credits,
            "days_of_week": self.days_of_week,
            "start_time": self.start_time,
            "end_time": self.end_time,
            "start_date": self.start_date,
            "end_date": self.end_date,
            "college_name": self.college_name,
            "location": self.location,
            "teacher": self.teacher,
        }

    def convert_dow(self, dow_int):
        s = ""
        if dow_int % 10 == 1:
            s += "Mon, "
        if int(dow_int / 10) % 10 == 1:
            s += "Tues, "
        if int(dow_int / 100) % 10 == 1:
            s += "Wed, "
        if int(dow_int / 1000) % 10 == 1:
            s += "Thur, "
        if int(dow_int / 10000) % 10 == 1:
            s += "Fri, "

        return s[:-2]

    def convert_time(self, t):
        hour = int(t / 100)
        am_or_pm = "AM"
        if hour >= 12:
            am_or_pm = "PM"
        if hour >= 13:
            hour -= 12
        minute = t % 100
        return f"{str(hour)}:{str(minute).zfill(2)}{am_or_pm}"


    def from_json(self, data):
        try:
            self.course_id = int(data['course_id'])
            self.course_code = int(data['course_code'])
            self.course_name = str(data['course_name'])
            self.credits = int(data['credits'])

            self.days_of_week = self.convert_dow(int(data['days_of_week']))
            self.start_time = self.convert_time(data['start_time'])
            self.end_time = self.convert_time(data['end_time'])

            # These should be UTC, not str
            self.start_date = str(data['start_date'])
            self.end_date = str(data['end_date'])

            self.college_name = str(data['college_name'])
            self.location = str(data['location'])
            self.teacher = str(data['teacher'])

        except KeyError:
            raise RuntimeError('Invalid field in json')
