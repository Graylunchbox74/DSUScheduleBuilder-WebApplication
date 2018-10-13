from flask.json import JSONEncoder


class UserModel():
    def create(**kwargs):
        """Creates user model. Fields required:
            - email
            - first_name
            - last_name
        """
        return UserModel.from_json(kwargs)

    def get_full_name(json):
        return f"{json['first_name']} {json['last_name']}"

    def __repr__(json):
        return f"UserModel({json['email']}, {json['first_name']}, {json['last_name']})"

    def from_json(json):
        ret = {}
        try:
            ret['email'] = json['email']
            ret['first_name'] = json['first_name']
            ret['last_name'] = json['last_name']
            return ret
        except KeyError as e:
            pass
