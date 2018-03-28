# -*- coding: utf-8 -*-
#
# Script run club and display club and visitor status.
# python3 club_run.py
# Created club instance, add visitors, send notify, update_visitor status.
# control command: -s - show status; -m - change music track, -r - random change visitors.
#
#
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
# Add random visitors
cotton_club.add_visitors_randoom(10)
# Notify about played music
cotton_club.played_music_notify()
# show club and visitor status
print(cotton_club.show_status())

# input command to update information or status.
input_command = ''

while True:
    try:
        input_command = input('enter command: ')
    except KeyboardInterrupt:
        break
    except EOFError:
        break

    if input_command == '-h':
        print('''
        Commands List:
            -s        Show information about visitors and play music.
            -m        Change played music to next track.
            -r        Random change visitors.
            -e        Exit terminal.
        ''')

    elif input_command == '-s':
        print (cotton_club.show_status())

    elif input_command == '-m':
        cotton_club.change_track()
        print (cotton_club.show_status())

    elif input_command == '-r':
        cotton_club.change_visitiors(10)
        print (cotton_club.show_status())

    elif input_command == '-e':
        break

print(CLUB_CLOSE_MSG)
