# -*- coding: utf-8 -*-
#
# Script run club and display club and visitor status.
#
import argparse

import club_lib as ccl


CLUB_WELCOME_MSG = (
    '''
    Welcome to the Cotton Club Info Terminal.
    Enter commands or -h.
    ''')
CLUB_CLOSE_MSG = 'Bye!'


print (CLUB_WELCOME_MSG)

# Create club object (Open Club)
cotton_club = ccl.Club('Cotton Club')
# show club and visitor status
print(cotton_club.show_status())

input_command = ''

while True:
    try:
        input_command = input('enter command: ')
    except KeyboardInterrupt:
        break
    except EOFError:
        break

    if input_command in ('help', '--help', '-h'):
        print('''
        Commands List:
            -s or --status     Show information about visitors and play music.
            -c or --change     Change played music to next track.
            -r or --random      Random change visitors.
            -a or --add        Add New Visitor.
            -h or --help       Help.
        ''')

    elif input_command in ('-s', '--status'):
        print (cotton_club.show_status())

    elif input_command in ('-c', '--change'):
        cotton_club.change_track()
        print (cotton_club.show_status())

    elif input_command in ('-r', '--random'):
        cotton_club.change_visitiors(10)
        print (cotton_club.show_status())

    elif input_command == 'exit':
        break

print(CLUB_CLOSE_MSG)
