# libs module


def get_next_biggest_with_use_same_digits(number):
    # Search algorithm:
    # Convert number to digits array.
    # Find rightmost position where current value less than next value.
    # From next right digits need find less digit from digits bigger than value in position found before.
    # Replace two value.
    # Sort right part by ascending.

    # Converting the number to the digits list.
    number_digits = [int(num_char) for num_char in str(number)]
    number_len = len(number_digits)

    # Replace positions variables.
    r_pos_1 = 0
    r_pos_1_val = 0
    r_pos_2 = 0
    r_pos_2_val = 0

    # Search position to replace.
    for pos, num in enumerate(number_digits):

        # Position where current value less than next.
        if pos < (number_len - 1) and num < number_digits[pos + 1]:
            r_pos_1 = pos
            r_pos_1_val = num
        # Find second replace value.
        if r_pos_1_val < num:
            r_pos_2 = pos
            r_pos_2_val = num

    # Bigger value not exist.
    if not r_pos_1 and not r_pos_1_val:
        next_biggest = 0
    # Get bigger value.
    else:
        # Replace two value.
        number_digits[r_pos_1] = r_pos_2_val
        number_digits[r_pos_2] = r_pos_1_val
        # Get number.
        next_biggest = int(''.join([
            str(num) for num in (
                # First part digits.
                number_digits[:r_pos_1+1]
                # Sort right part digits.
                + sorted(number_digits[r_pos_1+1:])
            )
        ]))

    return next_biggest


# def get_all_posible_track_to_map(map_n, map_m, str_pos_x, str_pos_y, end_pos_x, end_pos_y):
#     track_counts = 0
#
#     finded_tracks = []
#     looked_matrix = {}
#
#     looked_matrix[(str_pos_x, str_pos_y)] = get_posible_steps(map_str_pos_x, str_pos_y)
#
#     return track_counts


class MapTrackSearcher():

    def __init__(self, map_n, map_m, str_pos_x, str_pos_y, end_pos_x, end_pos, y):
        """

        :param map_n:
        :param map_m:
        :param str_pos_x:
        :param str_pos_y:
        :param end_pos_x:
        :param end_pos:
        :param y:
        """
        pass
