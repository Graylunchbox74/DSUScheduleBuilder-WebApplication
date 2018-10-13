class BaseModel():
    def __init__(self):
        raise RuntimeError(
            'Models should not be instantiated. Use `create` instead.')

    def to_str():
        raise NotImplementedError('`to_str` not implemented for this model.')

    def from_json():
        raise NotImplementedError(
            '`from_json` not implemented for this model.')
