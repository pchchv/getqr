import abc


class Image:
    def __init__(self, border, width, box_size):
        self.border = border
        self.width = width
        self.box_size = box_size
        self.pixel_size = (self.width + self.border * 2) * self.box_size
        self._image = self.new_image(**kwargs)

    @abc.abstractmethod
    def new_image(self, **kwargs):
        """
        Creates an image class. Subclasses must return the created class.
        """

    def get_image(self, **kwargs):
        """
        Returns the image class for later processing.
        """
        return self._image
