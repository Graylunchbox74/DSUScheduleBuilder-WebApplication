from server import app
from server.models.program_model import ProgramModel
from server.models.utils import facade_result_codes as FRC
import requests


def get_enrolled_programs(token):
    """

    Gets the programs (majors / minors) a user is currently enrolled in.

    Parameters
    ----------
        token : str
            The user's token

    Returns
    ------
        (int, list(ProgramModel))
    """

    try:
        response = requests.get(f"{app.config['API_ENDPOINT']}/user/getUsersPrograms?token={token}")
        json_response = response.json()

        if response.status_code == 401:
            return (FRC.NOT_AUTHENTICATED, [])
        if response.status_code == 400:
            return (FRC.SERVER_ERROR, [])

        programs = []
        for prog in json_response:
            program = ProgramModel.create(
                program_id=prog['ProgramID'],
                catalog_year=prog['CatalogYear'],
                major=prog['Major'],
                program=prog['Program'],
            )

            programs.append(program)

        return (FRC.SUCCESS, programs)
    except:
        return (FRC.CONNECTION_ERROR, [])
