class User():
    def __init__(self, **kwargs):
        """
        Creates user model. Fields required:
            - email
            - first_name
            - last_name
        """
        self.email = kwargs['email']
        self.first_name = kwargs['first_name']
        self.last_name = kwargs['last_name']
