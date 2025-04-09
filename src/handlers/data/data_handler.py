class DataHandler:
    def write(self, measurement: str, fields: dict, tags: dict = None):
        """
        Write data to the data store.

        Args:
            measurement (str): The measurement name.
            fields (dict): The fields to write.
            tags (dict, optional): The tags to write. Defaults to None.
        """
        raise NotImplementedError("Subclasses should implement this method.")
