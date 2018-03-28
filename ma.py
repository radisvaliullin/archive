# ma


class MyError(Exception):

    pass


class Parent(object):

    def __init__(self, arg_1 = 3):

        self.i = arg_1

    def fnc(self, arg_1, arg_2 = 2):

        if arg_1 == 7:
            raise MyError("Error text")
        return arg_1 * arg_2 * self.i

    def isFirst(self):
        return isinstance(self, First)

    @property
    def isSecond(self):
        return isinstance(self, Second)


class First(Parent):

    pass


class Second(Parent):

    pass


class A(First):

    pass
