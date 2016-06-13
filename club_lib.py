# -*- coding: utf-8 -*-
#
#
import random
import uuid
from abc import ABCMeta, abstractmethod

import itertools


class Observer(metaclass=ABCMeta):

    @abstractmethod
    def handler_observerable_notice(self):
        pass


class Observable(metaclass=ABCMeta):

    def __init__(self):

        self._observers = []

    def register_observer(self, observer: Observer):
        self._observers.append(observer)

    def unregister_observer(self, observer: Observer):
        self._observers.remove(observer)

    def notify_observer(self, *args, **kwargs):
        for observer in self._observers:
            observer.handler_observerable_notice(*args, **kwargs)


class Club(Observable):

    def __init__(self, name):
        super(Club, self).__init__()

        self.name = name
        self.music_box = MusicBox().random()

    def add_visitor(self, visitor=None, name=None, last_name=None, sex=None, dance_skill=None):

        if visitor is None:
            visitor = Visitor(name=name, last_name=last_name, sex=sex, dance_skill=dance_skill)

        self.register_observer(visitor)

    def remove_visitor(self, visitor):

        self.unregister_observer(visitor)

    def music_track_change(self):

        self.notify_observer(played_music=self.music_box.play_music)

    def add_visitors_randoom(self, count_visitors=None):

        for i in range(count_visitors or 0):
            self.register_observer(Visitor(random=True))

    def get_visitors_count_str(self):

        return str(len(self._observers))

    def show_status(self):

        visitors_str = self.get_visitors_list_str()

        status_str = '''
        {club_name}
        played music: {composition_name}
        played music style: {composition_style}
        next_track: {next_track}

        visitors count: {visitors_count}
        visitors:
        {visitors_header}
        {visitors}

        '''.format(
            club_name=self.name,
            composition_name=self.music_box.play_music.music_compos_name,
            composition_style=self.music_box.play_music.music_compos_style,
            next_track=self.music_box.next_track.music_compos_name,
            visitors_count=self.get_visitors_count_str(),
            visitors_header=''.join([
                'n'.ljust(3, '.'), ' ',
                'Name LastName'.ljust(20, '.'), ' ',
                'Sex'.ljust(4, '.'), ' ',
                'Dance_Skills'.ljust(16, '.'), ' ',
                'Head'.ljust(16, '.'), ' ',
                'Trunk'.ljust(16, '.'), ' ',
                'Hands'.ljust(16, '.'), ' ',
                'Legs'.ljust(16, '.'), ' ',
            ]),
            visitors=visitors_str,
        )
        return status_str

    def get_visitors_list_str(self):
        visitors_str = []
        for n, pers in enumerate(self._observers):
            v_n = str(n + 1)
            v_name_lastname = ' '.join([pers.name, pers.last_name])
            v_sex = pers.sex
            v_dance_skill = ' '.join([ds.dance_name for ds in pers.dance_skills])
            v_head = pers.head
            v_trunk = pers.trunk
            v_hands = pers.hands
            v_legs = pers.legs
            visitor_str = ''.join([
                v_n[:3].ljust(3, '.'), ' ',
                v_name_lastname[:20].ljust(20, '.'), ' ',
                v_sex[:4].ljust(4, '.'), ' ',
                v_dance_skill[:16].ljust(16, '.'), ' ',
                (v_head or 'None')[:16].ljust(16, '.'), ' ',
                (v_trunk or 'None')[:16].ljust(16, '.'), ' ',
                (v_hands or 'None')[:16].ljust(16, '.'), ' ',
                (v_legs or 'None')[:16].ljust(16, '.'),
            ])
            visitors_str.append(visitor_str)
        visitors_list_str = '\n        '.join(visitors_str)
        return visitors_list_str

    def change_track(self):
        self.music_box.change_track()
        self.music_track_change()

    def change_visitiors(self, count_visitors=None):
        self._visitors = []
        self.add_visitors_randoom(count_visitors)


class Visitor(Observer):

    name = ''
    last_name = ''
    sex = ''

    head = None
    trunk = None
    hands = None
    legs = None

    dance_skills = []

    def __init__(self, name=None, last_name=None, sex=None, dance_skill=None, random=False):
        if not random:
            self.name = name
            self.last_name = last_name
            self.sex = sex
            if dance_skill:
                self.dance_skills = dance_skill if isinstance(dance_skill, list) else [dance_skill]
        else:
            self.get_random_person()

    def handler_observerable_notice(self, played_music=None):

        if played_music:
            self.set_visitor_condition(played_music)

    def set_visitor_condition(self, music):
        dance_skill_by_music = self.get_dance_skill_by_music(music)
        if dance_skill_by_music:
            self.head = dance_skill_by_music.head
            self.trunk = dance_skill_by_music.trunk
            self.hands = dance_skill_by_music.hands
            self.legs = dance_skill_by_music.legs
        else:
            self.head = HEAD_DRUNK
            self.trunk = TRUNK_DRUNK
            self.hands = HANDS_DRUNK
            self.legs = LEGS_DRUNK

    def get_dance_skill_by_music(self, music):
        dance_skills = [
            ds for ds in self.dance_skills
            if music.music_compos_style in ds.dance_right_musics
        ]
        return dance_skills[0] if dance_skills else None

    def get_random_person(self):
        """
        Create random Visitor.
        """
        self.sex = random.choice(list(SEX.keys()))
        self.name = ''.join([
            'Some', SEX[self.sex], 'Name', uuid.uuid4().hex[:4]
        ])
        self.last_name = ''.join([
            'Some', SEX[self.sex], 'LastName', uuid.uuid4().hex[:4]
        ])
        self.dance_skills = list(set([
            random.choice(DANCE_LIST) for i in range(random.randint(1, 2))]))


class MusicBox():

    def __init__(self):

        self.play_list = []
        self.play_pos = 0
        self.play_music = None
        self.next_track = None

    def random(self):

        self.play_list = self.get_random_play_list()
        self.play_pos = 0
        self.play_music = self.play_list[self.play_pos]
        self.next_track = self.get_next_track()
        return self

    def get_random_play_list(self):

        play_list = [MusicCompos().random() for i in range(42)]
        return play_list

    def get_next_track(self):

        next_track = self.play_list[self.play_pos + 1] if self.play_pos < len(self.play_list) - 1 else 0
        return next_track

    def change_track(self):
        if self.play_list:
            self.play_pos = self.play_pos + 1 if self.play_pos < len(self.play_list) - 1 else 0
            self.play_music = self.play_list[self.play_pos]
            next_track = self.get_next_track()


#=======================================
# SEX
FEMALE = 'F'
MALE = 'M'
SEX = {FEMALE: 'Female', MALE: 'Male'}


# ======================================
# Music Style
RNB_MUSIC_STYLE = 'RnB'
ELHOUSE_MUSIC_STYLE = 'Electrohouse'
POP_MUSIC_STYLE = 'Pop'

MUSIC_STYLE_LIST = [RNB_MUSIC_STYLE, ELHOUSE_MUSIC_STYLE, POP_MUSIC_STYLE]


# ======================================
# Music Composition
class MusicCompos():

    def __init__(self):
        self.music_compos_name = ''
        self.music_compos_style = None

    def random(self):
        self.music_compos_style = random.choice(MUSIC_STYLE_LIST)
        self.music_compos_name = ''.join([
            'MusicCompos',
            self.music_compos_style,
            uuid.uuid4().hex[:4]
        ])
        return self


# ======================================
# Body Part Condition
# HEAD
HEAD_BACK_AND_FORTH = 'Head Back And Forth'
HEAD_LOW = 'Head Low'
HEAD_FLOWING_MOTION = 'Head Flowing Motion'
HEAD_DRUNK = 'Head Drunk'

# TRUNK
TRUNK_BACK_AND_FORTH = 'Trunk Back And Forth'
TRUNK_FLOWING_MOTION = 'Trunk Flowing Motion'
TRUNK_DRUNK = 'Trunk Drunk'

# HANDS
HANDS_BEND_ELBOW = 'Hands Bent At The Elbow'
HANDS_CIRCLE_ROTATING = 'Hands Circle Rotating'
HANDS_FLOWING_MOTION = 'Hands Flowing Motion'
HANDS_DRUNK = 'Hands Drunk'

# LEGS
LEGS_CROUCH = 'Legs Crouch'
LEGS_MOVE_RHYTHM = 'Legs Move Rhythm'
LEGS_FLOWING_MOTION = 'Legs Flowing Motion'
LEGS_DRUNK = 'Legs Drunk'


# ======================================
# Dance
class Dance():

    def __init__(self, name, head, trunk, hands, legs, right_music):

        self.dance_name = name

        self.head = head
        self.trunk = trunk
        self.hands = hands
        self.legs = legs

        self.dance_right_musics = [right_music] if isinstance(right_music, list) else [right_music]


HIP_HOP_DANCE = Dance(
    'HipHop',
    HEAD_BACK_AND_FORTH,
    TRUNK_BACK_AND_FORTH,
    HANDS_BEND_ELBOW,
    LEGS_CROUCH,
    RNB_MUSIC_STYLE,
)

RNB_DANCE = Dance(
    'RnB',
    HEAD_BACK_AND_FORTH,
    TRUNK_BACK_AND_FORTH,
    HANDS_BEND_ELBOW,
    LEGS_CROUCH,
    RNB_MUSIC_STYLE,
)

ELECTRO_DANCE = Dance(
    'ElectroDance',
    HEAD_LOW,
    TRUNK_BACK_AND_FORTH,
    HANDS_CIRCLE_ROTATING,
    LEGS_MOVE_RHYTHM,
    ELHOUSE_MUSIC_STYLE,
)

HOUSE_DANCE = Dance(
    'HouseDance',
    HEAD_LOW,
    TRUNK_BACK_AND_FORTH,
    HANDS_CIRCLE_ROTATING,
    LEGS_MOVE_RHYTHM,
    ELHOUSE_MUSIC_STYLE,
)

POP_DANCE = Dance(
    'PopDance',
    HEAD_FLOWING_MOTION,
    TRUNK_BACK_AND_FORTH,
    HANDS_FLOWING_MOTION,
    LEGS_FLOWING_MOTION,
    POP_MUSIC_STYLE,
)

DANCE_LIST = [HIP_HOP_DANCE, RNB_DANCE, ELECTRO_DANCE, HOUSE_DANCE, POP_DANCE]
