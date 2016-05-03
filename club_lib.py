# -*- coding: utf-8 -*-
#
#
import os
import random


class Club(object):

    _visitors = []

    play_music = None
    play_list = []

    def __init__(self, name):

        self.name = name
        self.play_list = self.get_random_play_list()
        self.play_music = self.play_list[0]
        self.play_music_pos = 0
        self.next_track = self.get_next_track()

        self.add_visitors_randoom(10)

    def add_visitor(
            self, visitor=None, name=None, last_name=None, sex=None,
            dance_skills=None):

        if visitor:
            self._visitors.append(visitor)
        else:
            self._visitors.append(Person(name, last_name, sex, dance_skills))

    def add_visitors_randoom(self, count_visitors=None):
        self._visitors.extend([Person() for i in xrange(count_visitors or 0)])
        self.set_visitors_condition()

    def set_visitors_condition(self):
        for person in self._visitors:
            person.set_person_condition(self.play_music)

    def show_status(self):

        visitors_str = self.get_visitors_list_str()

        status_str = u'''
        {club_name}
        played music: {composition_name}
        played music style: {composition_style}
        next_track: {next_track}

        visitors:
        {visitors_header}
        {visitors}

        '''.format(
            club_name=self.name,
            composition_name=self.play_music.music_compos_name,
            composition_style=(
                self.play_music.music_compos_style.music_style_name),
            next_track=self.next_track.music_compos_name,
            visitors_header=u''.join([
                u'n'.ljust(3, u'.'), u' ',
                u'Name LastName'.ljust(20, u'.'), u' ',
                u'Sex'.ljust(4, u'.'), u' ',
                u'Dance_Skills'.ljust(16, u'.'), u' ',
                u'Head'.ljust(16, u'.'), u' ',
                u'Body'.ljust(16, u'.'), u' ',
                u'Hands'.ljust(16, u'.'), u' ',
                u'Legs'.ljust(16, u'.'), u' ',
            ]),
            visitors=visitors_str,
        )
        return status_str

    def get_next_track(self):
        if not self.play_list:
            next_track = None
        if len(self.play_list) == 1:
            next_track = self.play_list[0]
        elif self.play_music_pos < (len(self.play_list) - 1):
            next_track = self.play_list[self.play_music_pos + 1]
        else:
            next_track = self.play_list[0]
        return next_track

    def get_visitors_list_str(self):
        visitors_str = []
        for n, pers in enumerate(self._visitors):
            v_n = unicode(n + 1)
            v_name_lastname = u' '.join([pers.name, pers.last_name])
            v_sex = pers.sex
            v_dance_skill = u' '.join([
                ds.dance_name for ds in pers.dance_skills])
            v_head = pers.body_condition['head'].cond_name
            v_body = pers.body_condition['body'].cond_name
            v_hands = pers.body_condition['hands'].cond_name
            v_legs = pers.body_condition['legs'].cond_name
            visitor_str = u''.join([
                v_n[:3].ljust(3, u'.'), u' ',
                v_name_lastname[:20].ljust(20, u'.'), u' ',
                v_sex[:4].ljust(4, u'.'), u' ',
                v_dance_skill[:16].ljust(16, u'.'), u' ',
                v_head[:16].ljust(16, u'.'), u' ',
                v_body[:16].ljust(16, u'.'), u' ',
                v_hands[:16].ljust(16, u'.'), u' ',
                v_legs[:16].ljust(16, u'.'),
            ])
            visitors_str.append(visitor_str)
        visitors_list_str = u'\n        '.join(visitors_str)
        return visitors_list_str

    def get_random_play_list(self):
        play_list = [MusicCompos() for i in xrange(42)]
        return play_list

    def change_track(self):
        if not self.play_list:
            self.play_music = None
            self.next_track = None
            self.play_music_pos = None
        if len(self.play_list) == 1:
            self.play_music = self.play_list[0]
            self.next_track = self.play_list[0]
            self.play_music_pos = 0
        elif self.play_music_pos < (len(self.play_list) - 1):
            self.play_music = self.play_list[self.play_music_pos + 1]
            self.play_music_pos += 1
            self.next_track = self.get_next_track()
        else:
            self.play_music_pos = 0
            self.play_music = self.play_list[0]
            self.next_track = self.get_next_track()

        self.set_visitors_condition()

    def change_visitiors(self, count_visitors=None):
        self._visitors = []
        self.add_visitors_randoom(count_visitors)


SEX = {u'F': u'Female', u'M': u'Male'}


class Person(object):

    name = u''
    last_name = u''
    sex = u''

    body_condition = {
        u'head': None,
        u'body': None,
        u'hands': None,
        u'legs': None,
    }

    dance_skills = []

    def __init__(self, name=None, last_name=None, sex=None, dance_skills=None):
        if (
            name and last_name and sex in SEX and
            isinstance(name, basestring) and isinstance(last_name, basestring)
        ):
            self.name = name
            self.last_name = last_name
            self.sex = sex
            if dance_skills:
                self.dance_skills = dance_skills
        else:
            self.get_random_person()

    def get_random_person(self):

        self.sex = random.choice(SEX.keys())
        self.name = u''.join([
            u'Some', SEX[self.sex], u'Name', os.urandom(4).encode('hex')])
        self.last_name = u''.join([
            u'Some', SEX[self.sex], u'LastName', os.urandom(4).encode('hex')])
        self.dance_skills = list(set([
            random.choice(DANCE_LIST) for i in xrange(random.randint(1, 2))]))

    def set_person_condition(self, music):
        dance_skill = self.get_dance_skill_by_music(music)
        if dance_skill:
            self.body_condition = dance_skill.body_condition
        else:
            self.body_condition = {
                'head': HeadDrunk,
                'body': BodyDrunk,
                'hands': HandsDrunk,
                'legs': LegsDrunk,
            }

    def get_dance_skill_by_music(self, music):
        dance_skills = [
            ds for ds in self.dance_skills
            if music.music_compos_style in ds.dance_right_musics
        ]
        return dance_skills[0] if dance_skills else None


# ======================================
# Music Style
class MusicStyle(object):

    music_style_name = u''


class RnBMusic(MusicStyle):

    music_style_name = u'RnB'


class ElHouseMusic(MusicStyle):

    music_style_name = u'Electrohouse'


class PopMusic(MusicStyle):

    music_style_name = u'Pop'


MusicStyleList = [RnBMusic, ElHouseMusic, PopMusic]


# ======================================
# Music Composition
class MusicCompos(object):

    music_compos_name = u''
    music_compos_style = None

    def __init__(self):
        self.music_compos_style = random.choice(MusicStyleList)
        self.music_compos_name = u''.join([
            u'MusicCompos',
            self.music_compos_style.music_style_name,
            os.urandom(4).encode('hex'),
        ])


# ======================================
# Body Part Condition
class BodyCondition(object):

    cond_name = u''


# Head
class HeadCond(BodyCondition):
    pass


class HeadBackAndForth(HeadCond):

    cond_name = u'Head Back And Forth'


class HeadLow(HeadCond):

    cond_name = u'Head Low'


class HeadFlowingMotion(HeadCond):

    cond_name = u'Head Flowing Motion'


class HeadDrunk(HeadCond):

    cond_name = u'Head Drunk'


# Body
class BodyCond(BodyCondition):
    pass


class BodyBackAndForth(BodyCond):

    cond_name = u'Body Back And Forth'


class BodyFlowingMotion(HeadCond):

    cond_name = u'Body Flowing Motion'


class BodyDrunk(BodyCond):

    cond_name = u'Body Drunk'


# Hands
class HandsCond(BodyCondition):
    pass


class HandsBendElbow(HandsCond):

    cond_name = u'Hands Bent At The Elbow'


class HandsCircleRotating(HandsCond):

    cond_name = u'Hands Circle Rotating'


class HandsFlowingMotion(HeadCond):

    cond_name = u'Hands Flowing Motion'


class HandsDrunk(HandsCond):

    cond_name = u'Hands Drunk'


# Legs
class LegsCond(BodyCondition):
    pass


class LegsCrouch(LegsCond):

    cond_name = u'Legs Crouch'


class LegsMoveRhythm(LegsCond):

    cond_name = u'Legs Move Rhythm'


class LegsFlowingMotion(HeadCond):

    cond_name = u'Legs Flowing Motion'


class LegsDrunk(LegsCond):

    cond_name = u'Legs Drunk'


# ======================================
# Dance
class Dance(object):

    dance_name = u''

    body_condition = {
        'head': None,
        'body': None,
        'hands': None,
        'legs': None,
    }

    dance_right_musics = []

    def __init__(self):

        pass


class HipHopDance(Dance):

    dance_name = u'HipHop'
    body_condition = {
        'head': HeadBackAndForth,
        'body': BodyBackAndForth,
        'hands': HandsBendElbow,
        'legs': LegsCrouch,
    }
    dance_right_musics = [RnBMusic]


class RnBDance(Dance):

    dance_name = u'RnB'
    body_condition = {
        'head': HeadBackAndForth,
        'body': BodyBackAndForth,
        'hands': HandsBendElbow,
        'legs': LegsCrouch,
    }
    dance_right_musics = [RnBMusic]


class ElectroDance(Dance):

    dance_name = u'ElectroDance'
    body_condition = {
        'head': HeadLow,
        'body': BodyBackAndForth,
        'hands': HandsCircleRotating,
        'legs': LegsMoveRhythm,
    }
    dance_right_musics = [ElHouseMusic]


class HouseDance(Dance):

    dance_name = u'HouseDance'
    body_condition = {
        'head': HeadLow,
        'body': BodyBackAndForth,
        'hands': HandsCircleRotating,
        'legs': LegsMoveRhythm,
    }
    dance_right_musics = [ElHouseMusic]


class PopDance(Dance):

    dance_name = u'PopDance'
    body_condition = {
        'head': HeadFlowingMotion,
        'body': BodyFlowingMotion,
        'hands': HandsFlowingMotion,
        'legs': LegsFlowingMotion,
    }
    dance_right_musics = [ElHouseMusic]


DANCE_LIST = [HipHopDance, RnBDance, ElectroDance, HouseDance, PopDance]
