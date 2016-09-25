##
# tesk_task main module
##
import argparse
import timeit

from task_solver_libs import get_next_biggest_with_use_same_digits


start_time = timeit.default_timer()


if __name__ == '__main__':

    arg_parser = argparse.ArgumentParser(description='Some test task solcer.')
    arg_parser.add_argument('tasks_list', type=str, )
    args = arg_parser.parse_args()

    tasks_list_file_path = args.tasks_list
    tasks_file = open(tasks_list_file_path, 'r')

    for line in tasks_file:

        if line.strip().isdecimal():

            next_biggest = get_next_biggest_with_use_same_digits(int(line))
            output_value = str(next_biggest)

        elif len(line.strip().split()) == 6:

            track_counts = get_all_posible_track_to_map()
            output_value = str(next_biggest)

        else:
            output_value = 'Unparsed line.'

        print (output_value + '-----')

    print (timeit.default_timer() - start_time)
