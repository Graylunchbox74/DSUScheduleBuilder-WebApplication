from server import app
from server.models.utils import facade_result_codes as FRC
import requests


def enroll_in_course(token, course_id):
    """

    Enrolls a user in the course with id `course_id`.

    Parameters
    ----------
        token : str
            The session token for the user

        course_id : int
            The id of the course being enrolled in

    Returns
    -------
        int
        Success code from facade_result_codes
    """

    data = {
        "key": token,
        "courseID": course_id,
    }

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/user/enrollInCourse", data=data)
        json_response = response.json()

        if response.status_code == 401:
            return FRC.NOT_AUTHENTICATED

        if response.status_code == 400:
            return FRC.SERVER_ERROR

        return FRC.SUCCESS

    except:
        return FRC.CONNECTION_ERROR
