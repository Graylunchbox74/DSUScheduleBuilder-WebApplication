from server.models.base_model import BaseModel


class UserModel(BaseModel):
    def create(**kwargs):
        """Creates user model. Fields required:
            - id
            - email
            - first_name
            - last_name
        """
        return UserModel.from_json(kwargs)

    def get_full_name(json):
        return f"{json['first_name']} {json['last_name']}"

    def to_str(json):
        return f"UserModel({json['id']}, {json['email']}, {json['first_name']}, {json['last_name']})"

    def from_json(json):
        ret = {}
        try:
            ret['id'] = int(json['id'])
            ret['email'] = str(json['email'])
            ret['first_name'] = str(json['first_name'])
            ret['last_name'] = str(json['last_name'])
            return ret
        except KeyError as e:
            raise RuntimeError('Invalid field in json')
