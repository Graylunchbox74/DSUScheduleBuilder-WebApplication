from server.models.base_model import BaseModel


class UserModel(BaseModel):
    def create(**kwargs):
        """Creates user model. Fields required:
            - id
            - email
            - first_name
            - last_name
            - token
        """
        user = UserModel()
        user.from_json(kwargs)
        return user

    def get_full_name(self):
        return f"{self.first_name} {self.last_name}"

    def to_str(self):
        return f"UserModel({self.id}, {self.email}, {self.first_name}, {self.last_name})"

    def to_json(self):
        return {
            'id': self.id,
            'email': self.email,
            'first_name': self.first_name,
            'last_name': self.last_name,
            'token': self.token,
        }

    def from_json(self, json):
        try:
            self.id = int(json['id'])
            self.email = str(json['email'])
            self.first_name = str(json['first_name'])
            self.last_name = str(json['last_name'])
            self.token = str(json['token'])

        except KeyError as e:
            raise RuntimeError('Invalid field in json')
