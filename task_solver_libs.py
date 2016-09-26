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


class MapTrackSearcher:

    def __init__(self, map_n, map_m, spos_x, spos_y, epos_x, epos_y):
        """

        :param map_n: map rows;
        :param map_m: map columns;
        :param spos_x: start position rows coordinate;
        :param spos_y: start position columns coordinate;
        :param epos_x: end position rows coordinate;
        :param epos_y: end position columns coordinate,
        """
        self.map_n = map_n
        self.map_m = map_m
        self.spos_x = spos_x
        self.spos_y = spos_y
        self.spos = (self.spos_x, self.spos_y)
        self.epos_x = epos_x
        self.epos_y = epos_y
        self.epos = (self.epos_x, self.epos_y)

        self.track_counts = 0

        self.curr_pos = self.spos
        self.curr_chain = [self.spos]
        self.found_tracks = []
        self.looking_tree = []

    def get_all_posible_tracks_count(self):

        self.looking_tree = [self.get_next_posible_steps(*self.curr_pos)]

        while self.looking_tree:

            if self.looking_tree[-1]:

                self.curr_pos = self.looking_tree[-1].pop()

                if self.curr_pos == self.epos:

                    self.track_counts += 1

                else:

                    self.curr_chain.append(self.curr_pos)
                    self.looking_tree.append(self.get_next_posible_steps(*self.curr_pos))

            else:

                self.looking_tree.pop()
                self.curr_chain.pop()

        all_tracks_count = self.track_counts
        return all_tracks_count

    def get_next_posible_steps(self, pos_x, pos_y):
        allows_points = []
        if (pos_x - 1) >= 0:
            allows_points.append((pos_x - 1, pos_y))
        if (pos_y - 1) >= 0:
            allows_points.append((pos_x, pos_y - 1))
        if (pos_x + 1) <= (self.map_n - 1):
            allows_points.append((pos_x + 1, pos_y))
        if (pos_y + 1) <= (self.map_m - 1):
            allows_points.append((pos_x, pos_y + 1))
        allows_points = [
            pos for pos in allows_points if (pos not in self.curr_chain) or (pos == self.epos)
        ]
        return allows_points


def get_list_of_string_patterns(in_string):

    # Minimum size for pattern chars.
    min_patt_len = 2
    # Minimum pattern repeats.
    min_patt_rep = 2
    # Patterns positins store.
    patterns_pos = {}
    # Patterns list.
    patterns = []

    for str_idx in range(len(in_string)):
        for end_idx in range((str_idx + min_patt_len - 1), len(in_string)):
            pattern = in_string[str_idx: end_idx + 1]
            if pattern in patterns_pos:
                patterns_pos[pattern].append(str_idx)
            else:
                patterns_pos[pattern] = [str_idx]

    patterns = sorted(
        [patt for patt, patt_pos in patterns_pos.items() if len(patt_pos) >= min_patt_rep],
        key=len, reverse=True
    )
    return patterns
