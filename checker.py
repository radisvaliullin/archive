import unittest
from random import randint

from sqlalchemy import create_engine, Unicode, Column, Integer
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

engine = create_engine('sqlite:///:memory:', echo=False)
Session = sessionmaker(bind=engine)

Base = declarative_base()


class Order(Base):
    __tablename__ = 'orders'

    id = Column(Integer, nullable=False, primary_key=True)
    name = Column(Unicode, nullable=False)
    state = Column(Integer, nullable=False, index=True)  # accepted = 1, hold = 0


def mark_random_orders_accepted(orders):

    for order in orders:
        state = randint(0,1)
        order.state = state


def check_hold_orders_in_order_table():

    session = Session()

    orders = session.query(Order).filter_by(state=0)[:100]

    while orders:

        mark_random_orders_accepted(orders)
        session.commit()

        orders = session.query(Order).filter_by(state=0)[:100]


# Тest checker
if __name__ == '__main__':

    Base.metadata.create_all(engine)

    session = Session()
    for i in range(1000):
        state = randint(0,1)
        o = Order(name=str(i), state=state)
        session.add(o)
        session.commit()

    oq = session.query(Order).filter_by(state=0)
    print("orders with hold state ", oq.count())

    check_hold_orders_in_order_table()

    oq = session.query(Order).filter_by(state=0)
    print("orders with hold state ", oq.count())