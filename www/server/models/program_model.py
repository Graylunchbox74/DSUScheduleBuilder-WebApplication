from server.models.base_model import BaseModel


class ProgramModel(BaseModel):
    def create(**kwargs):
        """ Creates a program based on the following args:
        - program_id
        - catalog_year
        - major
        - program
        """

        program = ProgramModel()
        program.from_json(kwargs)
        return program

    def to_json(self):
        return {
            "program_id": self.program_id,
            "catalog_year": self.program_id,
            "major": self.major,
            "program": self.program,
        }

    def from_json(self, data):
        try:
            self.program_id = int(data['program_id'])
            self.catalog_year = int(data['catalog_year'])
            self.major = int(data['major'])
            self.program = str(data['program'])

        except KeyError:
            raise RuntimeError('Invalid field in json')
