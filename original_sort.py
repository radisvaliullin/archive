# The module realizing specific algorithm to sort array with six random number elements
#  from 1 to 6 example (5,6,4,1,3,2) to (1,2,3,4,5,6).
# For more information look README.
import argparse
from itertools import tee, zip_longest


SORTED_ARRAY = '123456'
ARRAY_LEN = len(SORTED_ARRAY)

TO_MOVE_POSSIBLE_POSITIONS = (0, 1, 2, 3, 4)
SECOND_MOVE_ELEMENT = 1


def get_for_move_possible_positions(curr_to_move_pos):

    for_move_pos = [
        p for p in TO_MOVE_POSSIBLE_POSITIONS
        if (p + SECOND_MOVE_ELEMENT) < curr_to_move_pos or p > (curr_to_move_pos + SECOND_MOVE_ELEMENT)
    ]

    return for_move_pos


def original_sort(need_sort_array):

    all_possible_array_graph = all_possible_array_graph_build(need_sort_array)

    out_path_trace = []

    if SORTED_ARRAY in all_possible_array_graph:

        sort_path = bfs(all_possible_array_graph, need_sort_array, SORTED_ARRAY)

        out_path_trace = get_out_path_trace(sort_path, all_possible_array_graph)

    return out_path_trace


def get_out_path_trace(sort_path, graph):

    to_move_element = 0
    for_move_element = 1
    path_iter, next_path_iter = tee(sort_path, 2)
    next(next_path_iter)

    out_path_trace = []

    for p, n_p in zip_longest(path_iter, next_path_iter, fillvalue=''):

        out_path_trace_item = p
        if n_p:
            out_path_trace_item += ''.join([
                ' ', graph[p][n_p][to_move_element], '<->', graph[p][n_p][for_move_element],
            ])
        else:
            out_path_trace_item += '.'
        out_path_trace.append(out_path_trace_item)

    return out_path_trace


def all_possible_array_graph_build(begin_array):

    graph = {}
    new_nodes = [begin_array]

    while new_nodes:

        current_array = new_nodes.pop(0)

        if current_array not in graph:

            curr_arr_relate_arrays = get_array_related_arrays(current_array)
            graph[current_array] = curr_arr_relate_arrays
            new_nodes.extend(curr_arr_relate_arrays.keys())

    return graph


def get_array_related_arrays(array):

    related_arrays = {}

    for to_move_pos in TO_MOVE_POSSIBLE_POSITIONS:

        for_move_possible_positions = get_for_move_possible_positions(to_move_pos)

        for for_move_pos in for_move_possible_positions:

            related_array, to_move_val, for_move_val = get_new_array_with_swap_pairs(array, to_move_pos, for_move_pos)
            if related_array not in related_arrays:
                related_arrays[related_array] = (to_move_val, for_move_val)

    return related_arrays


def get_new_array_with_swap_pairs(array_to_swap, to_move_pos, for_move_pos):

    to_move_val_0 = array_to_swap[to_move_pos]
    to_move_val_1 = array_to_swap[to_move_pos + SECOND_MOVE_ELEMENT]
    for_move_val_0 = array_to_swap[for_move_pos]
    for_move_val_1 = array_to_swap[for_move_pos + SECOND_MOVE_ELEMENT]

    swap_elements_list = list(array_to_swap)

    swap_elements_list[to_move_pos] = for_move_val_0
    swap_elements_list[to_move_pos + SECOND_MOVE_ELEMENT] = for_move_val_1
    swap_elements_list[for_move_pos] = to_move_val_0
    swap_elements_list[for_move_pos + SECOND_MOVE_ELEMENT] = to_move_val_1

    swaped_array = u''.join(swap_elements_list)
    to_move_val = to_move_val_0 + to_move_val_1
    for_move_val = for_move_val_0 + for_move_val_1
    return swaped_array, to_move_val, for_move_val


def bfs(graph, start, end):
    """
    Breadth-first_search
    """
    # maintain a queue of paths
    queue = []
    # push the first path into the queue
    queue.append([start])
    while queue:
        # get the first path from the queue
        path = queue.pop(0)
        # get the last node from the path
        node = path[-1]
        # path found
        if node == end:
            return path
        # enumerate all adjacent nodes, construct a new path and push it into the queue
        for adjacent in graph.get(node).keys():
            new_path = list(path)
            new_path.append(adjacent)
            queue.append(new_path)


################################################################
# Main
################################################################
if __name__ == '__main__':

    arg_parser = argparse.ArgumentParser(
        '\n Welcome original_sort.py.'
        '\n Use to sort six unique elements (1-6) array, like 123456 or 654321 or 354216, to get 123456.'
        '\n to run: python3 original_sort.py array'
    )
    arg_parser.add_argument('array', type=str, help='Six unique elements (1-6) array, like 123456 or 654321 or 354216.')
    args = arg_parser.parse_args()

    errors = []
    if len(args.array) != len(SORTED_ARRAY):
        errors.append('Incorrect numbers of elements in array, must be six.')
    if len(args.array) == len(SORTED_ARRAY) and len(set(args.array)) != len(SORTED_ARRAY):
        errors.append('Duplicate numbers. Elements in array must be unique.')
    if set(args.array).difference(set(SORTED_ARRAY)):
        errors.append(
            'Array contains incorrect values: {val}. Enable only 1 2 3 4 5 6.'.format(
                val=', '.join([v for v in set(args.array).difference(set(SORTED_ARRAY))])))

    out_trace = errors or original_sort(args.array) or ['Sorry, but enter array haven\'t sorted']

    for o in out_trace:
        print (o)
