# -*- coding: utf-8 -*-
#
#
import club_lib as ccl


print(
    u'''
    Welcome to the Cotton Club Info Terminal.
    Enter commands or help.
    ''')

cotton_club = ccl.Club(u'Cotton Club')

print(cotton_club.show_status())

r_in = None

while r_in != u'exit':
    try:
        r_in = raw_input(u'...')
    except KeyboardInterrupt:
        break
    except EOFError:
        break

    if r_in in (u'help', u'--help', u'-h'):
        print(u'''
        Commands List:
            -s or --status or --club_visitors_status    Its show information about visitors and play music.
            -c or --change or --change_track            Its change played music to next track.
            -r or --rnd_v or --random_visitors          Its random change visitors.
        ''')

    elif r_in in (u'-s', u'--status', u'--club_visitors_status'):
        print (cotton_club.show_status())

    elif r_in in (u'-c', u'--change', u'--change_track'):
        cotton_club.change_track()
        print (cotton_club.show_status())

    elif r_in in (u'-r', u'--rnd_v', u'--random_visitors'):
        cotton_club.change_visitiors(10)
        print (cotton_club.show_status())

print(u'Bye.')
