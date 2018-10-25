from server import app
from server.models.user_model import UserModel
from server.models.utils import facade_result_codes as FRC
import requests


def login_admin(email, password):
    """
    Logs admin in using their credentials.

    Parameters
    ----------
        email : str
            the admin's email

        password : str
            the admin's password

    Returns
    -------
        (int, UserModel)
        Tuple of facade result code and admin model.
    """

    data = {
        "email": email,
        "password": password,
    }

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/admin/login", data=data)
        json_response = response.json()

        # Assume not authenticated until proved otherwise
        result_code = FRC.NOT_AUTHENTICATED

        if response.status_code == 200 and json_response['studentID'] != 0:
            result_code = FRC.SUCCESS

        user = None
        if result_code == FRC.SUCCESS:
            user = UserModel.create(
                id=json_response['studentID'],
                email=json_response['email'],
                first_name=json_response['firstName'],
                last_name=json_response['lastName'],
                token=json_response['token'],
            )

        return (result_code, user)
    except:
        return (FRC.CONNECTION_ERROR, None)


def check_admin_token(token):
    """

    Checks a the user token to see if it is still active.

    Parameters
    ----------
        token : str
            The admin token

    Returns
    -------
        bool
           Whether or not its valid
    """

    try:
        response = requests.get(f"{app.config['API_ENDPOINT']}/admin/checkToken?token={token}")

        return response.status_code == 200

    except:
        return False


def logout_admin(token):
    """

    Logs a user out.

    Parameters
    ----------
        token : str
            The admin token

    Returns
    -------
        None
    """

    data = {
        "token": token
    }

    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/admin/logout", data=data)
    except:
        return
