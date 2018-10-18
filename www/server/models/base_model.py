class BaseModel():
    @staticmethod
    def create(**kwargs):
        raise NotImplementedError('`create` method not defined.')

    def __init__(self):
        pass

    def to_str():
        raise NotImplementedError('`to_str` not implemented for this model.')

    def to_json():
        raise NotImplementedError('`to_json` not implemented for this model.')

    def from_json():
        raise NotImplementedError('`from_json` not implemented for this model.')
