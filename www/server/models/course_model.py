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

    def from_json(self, data):
        try:
            self.course_id = int(data['course_id'])
            self.course_code = int(data['course_code'])
            self.course_name = str(data['course_name'])
            self.credits = int(data['credits'])

            self.days_of_week = int(data['days_of_week'])
            self.start_time = int(data['start_time'])
            self.end_time = int(data['end_time'])

            # These should be UTC, not str
            self.start_date = str(data['start_date'])
            self.end_date = int(data['end_date'])

            self.college_name = str(data['college_name'])
            self.location = str(data['location'])
            self.teacher = str(data['teacher'])

        except KeyError:
            raise RuntimeError('Invalid field in json')
