##
# tesk_task main module
##
import argparse
import timeit

from task_solver_libs import get_next_biggest_with_use_same_digits, MapTrackSearcher

start_time = timeit.default_timer()


if __name__ == '__main__':

    arg_parser = argparse.ArgumentParser(description='Some test task solcer.')
    arg_parser.add_argument('tasks_list', type=str, )
    args = arg_parser.parse_args()

    tasks_list_file_path = args.tasks_list
    tasks_file = open(tasks_list_file_path, 'r')

    for line in tasks_file:

        line_vals = line.strip().split()

        if len(line_vals) == 1 and line_vals[0].isdecimal():

            next_biggest = get_next_biggest_with_use_same_digits(int(line))
            output_value = str(next_biggest)

        elif len(line_vals) == 6:

            line_vals = [int(val) for val in line_vals]
            mts = MapTrackSearcher(line_vals[0], line_vals[1], line_vals[2], line_vals[3], line_vals[4], line_vals[5])
            track_counts = mts.get_all_posible_tracks_count()
            output_value = str(track_counts)

        else:
            output_value = 'Unparsed line.'

        print(output_value + '-----')

    print(timeit.default_timer() - start_time)
