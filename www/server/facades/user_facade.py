from server import app
from server.models.user_model import UserModel
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
        (bool, UserModel, bool)
        Tuple of is validated, user model, and if an error occurred
    """

    data = {
        "email": email,
        "password": password,
    }

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/user/validateUser", data=data)
        json_response = response.json()

        is_validated = response.status_code == 200 and json_response['StudentID'] != 0

        user = {}
        if is_validated:
            user = UserModel.create(
                id=json_response['StudentID'],
                email=json_response['Email'],
                first_name=json_response['firstName'],
                last_name=json_response['lastName']
            )

        return (is_validated, user, False)
    except:
        return (False, {}, True)


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
        (bool, json, bool)
        Tuple of successful, the json response and if an error occurred
    """
    data = kwargs

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/user/newUser", data=data)
        json_response = response.json()

        success = response.status_code == 200 and json_response.get('errorMsg', "") == ""
        return (success, json_response, False)
    except:
        return (False, {}, True)
