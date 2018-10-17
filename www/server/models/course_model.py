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
        return CourseModel.from_json(kwargs)

    def to_str(data):
        return f"CourseModel({data['course_id']}, {data['course_code']}, {data['course_name']})"

    def from_json(data):
        ret = {}
        try:
            ret['course_id'] = int(data['course_id'])
            ret['course_code'] = int(data['course_code'])
            ret['course_name'] = str(data['course_name'])
            ret['credits'] = int(data['credits'])

            ret['days_of_week'] = int(data['days_of_week'])
            ret['start_time'] = int(data['start_time'])
            ret['end_time'] = int(data['end_time'])

            # These should be UTC, not str
            ret['start_date'] = str(data['start_date'])
            ret['end_date'] = int(data['end_date'])

            ret['college_name'] = str(data['college_name'])
            ret['location'] = str(data['location'])
            ret['teacher'] = str(data['teacher'])

            return ret
        except KeyError:
            raise RuntimeError('Invalid field in json')
