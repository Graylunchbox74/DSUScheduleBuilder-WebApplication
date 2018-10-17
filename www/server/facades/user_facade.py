from server import app
from server.models.user_model import UserModel
from server.models.utils import facade_result_codes as FRC
import requests


def validate_user(email, password):
    """
    Validates a user's credentials.

    Parameters
    ----------
        email : str
            the user's email

        password : str
            the user's password

    Returns
    -------
        (int, UserModel)
        Tuple of facade result code and user model.
    """

    data = {
        "email": email,
        "password": password,
    }

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/user/validateUser", data=data)
        json_response = response.json()

        # Assume not authenticated until proved otherwise
        result_code = FRC.NOT_AUTHENTICATED

        if response.status_code == 200 and json_response['StudentID'] != 0:
            result_code = FRC.SUCCESS

        user = {}
        if result_code == FRC.SUCCESS:
            user = UserModel.create(
                id=json_response['StudentID'],
                email=json_response['Email'],
                first_name=json_response['firstName'],
                last_name=json_response['lastName']
            )

        return (result_code, user)
    except:
        return (FRC.CONNECTION_ERROR, {})


def register_user(**kwargs):
    """
    Attempts to register a new user.

    Parameters
    ----------
        email : str
            the user's email

        password : str
            the user's password

        firstName : str
            the user's first name

        lastName : str
            the user's last name

    Returns
    -------
        (int, json)
        Tuple of facade result code and the json response.
    """
    data = kwargs

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/user/newUser", data=data)
        json_response = response.json()

        result_code = FRC.SERVER_ERROR

        if response.status_code == 200:
            result_code = FRC.SUCCESS

        return (result_code, json_response)
    except:
        return (FRC.CONNECTION_ERROR, {})
