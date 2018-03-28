#####################################################
# test.py
#####################################################

from ma import *
from mb import B

def test():
    a = A()
    b = B(5)
    
    assert(a.i == 3)
    assert(a.fnc(2) == 2 * 2 * 3)
    assert(b.fnc(10, 4) == 10 * 4 * 5)
    assert(a.isFirst() == 1)
    assert(a.isSecond == 0)
    assert(b.isFirst() == 0)
    assert(b.isSecond == 1)
    
    assert(isinstance(a, First) == 1)
    assert(isinstance(b, Second) == 1)
    assert(isinstance(a, Parent) == 1)
    assert(isinstance(b, Parent) == 1)

    try:
        a.fnc(7)
    except MyError, v:
        if str(v) != "Error text":
            assert(0)
    else:
        assert(0)

    try:
        a.isSecond = 2
    except AttributeError, v:
        pass
    else:
        assert(0)



if __name__ == "__main__":
    test()
    print "done"


#####################################################


