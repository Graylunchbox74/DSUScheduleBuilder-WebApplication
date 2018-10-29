from server.models.utils import facade_result_codes as FRC
from server import app

import requests

def add_program(token, program_name, is_major, catalog_year):
    """
    
    Adds a program to the list of registered programs.

    Parameters
    ---------
        token : str
            The admin's token

        program_name : str
            self explanitory

        major : boolean
            Whether its a major or not

        catalog_year : str
            The catatlog year (2018, 2019, ETC)

    Returns
    -------
        int

        Facade result code
    """

    data = {
        "token": token,
        "program": program_name,
        "major": "1" if is_major else "0",
        "catalogYear": catalog_year,
    }
    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/admin/addProgram", data=data)

        if response.status_code == 401:
            return FRC.NOT_AUTHENTICATED
        if response.status_code == 400:
            return FRC.SERVER_ERROR

        return FRC.SUCCESS
    except:
        return FRC.NOT_AUTHENTICATED

def delete_program(token, program_id):
    """
    
    Adds a program to the list of registered programs.

    Parameters
    ---------
        token : str
            The admin's token

        program_id : str
            The program's id

    Returns
    -------
        int

        Facade result code
    """
    data = {
        "token": token,
        "programID": program_id
    }
    try:
        response = requests.post(f"{app.config['API_ENDPOINT']}/admin/deleteProgram", data=data)

        if response.status_code == 401:
            return FRC.NOT_AUTHENTICATED
        if response.status_code == 400:
            return FRC.SERVER_ERROR

        return FRC.SUCCESS
    except:
        return FRC.NOT_AUTHENTICATED
