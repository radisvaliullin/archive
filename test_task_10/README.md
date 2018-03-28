# test_task_10
Some test task on Python

#### Description
```
"""
В таблице 'orders' находится 15млн записей. Необходимо проверить все заказы в статусе 'hold' пачками по 100шт.
Статус заказа проверяется в функции 'mark_random_orders_accepted', эта функция ставит рандомное кол-во заказов
в статус 'accepted', т.е. '1'. Кол-во, переведенных в статус 'accepted' заказов неизвестно.
Необходимо написать оптимальное решение. Нельзя выгружать все заказы в память(вызов .all() в SQLAlchemy). 
"""

class Order(Base):
    __tablename__ = 'orders'

    id = Column(BigInteger, nullable=False, primary_key=True, autoincrement=True)
    name = Column(Unicode, nullable=False)
    state = Column(Integer, nullable=False, index=True)   # accepted = 1, hold = 0
    
def mark_random_orders_accepted(orders):
  pass
```
