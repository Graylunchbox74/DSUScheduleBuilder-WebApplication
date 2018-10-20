from server import app
from server.models.course_model import CourseModel
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
        "token": token,
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


def get_enrolled_courses(token):
    """

    Gets the courses that a user is currently enrolled in.

    Parameters
    ----------
        token : str
        The session token for the user

    Returns
    -------
        (int, list(CourseModel))
        Tuple of facade result code and potentially empty list of courses.
    """

    try:
        response = requests.get(f"{app.config['API_ENDPOINT']}/user/getEnrolledCourses/{token}")
        json_response = response.json()

        if response.status_code == 401:
            return (FRC.NOT_AUTHORIZED, [])
        elif response.status_code == 400:
            return (FRC.SERVER_ERROR, [])

        courses = []

        for course in json_response:
            c = CourseModel.create(
                course_id=course['CourseID'],
                course_code=course['CourseCode'],
                course_name=course['CouseName'],
                credits=course['Credits'],
                days_of_week=course['DaysOfWeek'],
                start_time=course['StartTime'],
                end_time=course['EndTime'],
                start_date=course['StartDate'],
                end_date=course['EndDate'],
                college_name=course['CollegeName'],
                location=course['Location'],
                teacher=course['Teacher'],
            )
            courses.append(c)

        return (FRC.SUCCESS, courses)
    except:
        return (FRC.CONNECTION_ERROR, [])


def search_for_courses(token, kwargs):
    """

    Searches for courses given a set of parameters

    Parameters
    ----------
        token: str
            The token given when the user the logs in

        college_name: str?
        course_code: str?
        teacher: str?
        course_name: str?
        semester: str?
        location: str?

    Returns
    -------
        (int, list(CourseModel))
    """

    data = {
        "token": token
    }
    if kwargs.get('college_name') is not None:
        data["collegeName"] = kwargs.get('college_name')
    if kwargs.get('course_code') is not None:
        data["courseCode"] = kwargs.get('course_code')
    if kwargs.get('teacher') is not None:
        data["teacherName"] = kwargs.get('teacher')
    if kwargs.get('course_name') is not None:
        data["courseName"] = kwargs.get('course_name')
    if kwargs.get('location') is not None:
        data["location"] = kwargs.get('location')
    if kwargs.get('semester') is not None:
        data["semester"] = kwargs.get('semester')

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/user/searchForCourse", data=data)
        json_response = response.json()

        if response.status_code == 401:
            return (FRC.NOT_AUTHORIZED, [])
        elif response.status_code == 400:
            return (FRC.SERVER_ERROR, [])

        courses = []

        for course in json_response:
            c = CourseModel.create(
                course_id=course['CourseID'],
                course_code=course['CourseCode'],
                course_name=course['CourseName'],
                credits=course['Credits'],
                days_of_week=course['DaysOfWeek'],
                start_time=course['StartTime'],
                end_time=course['EndTime'],
                start_date=course['StartDate'],
                end_date=course['EndDate'],
                college_name=course['CollegeName'],
                location=course['Location'],
                teacher=course['Teacher'],
            )
            courses.append(c)

        return (FRC.SUCCESS, courses)
    except:
        return (FRC.CONNECTION_ERROR, [])
