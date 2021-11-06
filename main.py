def _check_box_size(size):
    if int(size) <= 0:
        raise ValueError(f"Invalid box size. Must be larger than 0")


def _check_border(size):
    if int(size) <= 0:
        raise ValueError(f"Invalid border value. Must be larger than 0")


class QRCode:
    def __init__(self, box_size=10,
                 border=2):
        _check_box_size(box_size)
        _check_border(border)
