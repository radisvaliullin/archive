# some classes test module.
import uuid


class Asset:
    """
    Asset class - for some asset object.
    """

    def __init__(self, name: str, price: float):
        """
        Asset object initialisation.
        :param name: stores the name of the asset.
        :param price: keeps asset price.
        """

        # Validation of name parameter. Name can't be more than 100 character.
        assert name and isinstance(name, str) and len(name) <= 100, 'Wrong Asset \'name\' parameter.'
        # The price parameter validation.
        assert isinstance(price, float) and price > 0.0, 'Wrong Asset \'price\' parameter.'

        # Set random id property.
        self._id = uuid.uuid4().int
        # Set name property.
        self._name = name
        # Set price property.
        self._price = price

    @property
    def id(self):
        """
        Id property getter. Use to get id property.
        """
        return self._id

    @property
    def name(self):
        """
        Name property getter. Use to get name property.
        """
        return self._name

    @property
    def price(self):
        """
        Price property getter. Use to get price property.
        """
        return self._price

    def __repr__(self):
        """
        Reload __repr__ method for display more actual information.
        """
        return 'Asset: id - {id}; name - {name}; price - {price}.'.format(id=self.id, name=self.name, price=self.price)


class Company:
    """
    Company class - used to keep information about some company.
    """

    def __init__(self, name: str, employees_number: int = 1, balance: float=0.0, assets: list = None):
        """
        Company object initialisation.
        :param name: keeps company name.
        :param employees_number: keeps company employees number.
        :param balance: keeps company balance.
        :param assets: keeps company assets.
        """

        # Company parameters validation.
        # Name can't be more than 100 character.
        assert name and isinstance(name, str) and len(name) <= 100, 'Wrong Asset \'name\' parameter.'
        assert isinstance(employees_number, int), 'Wrong Asset \'employees_number\' parameter.'
        assert isinstance(balance, float) and balance >= 0.0, 'Wrong Asset \'balance\' parameter.'
        # In assets list must stored only Asset objects.
        assert (
            assets and isinstance(assets, list) and all(isinstance(a, Asset) for a in assets) or
            assets == [] or assets is None
        ), 'Wrong Asset \'assets\' parameter.'

        # Set property values.
        self._id = uuid.uuid4().int
        self._name = name
        self._employees_number = employees_number
        self._balance = balance
        # If assets none need set empty list.
        self._assets = assets or []

    @property
    def id(self):
        """
        Id property getter. Use to get id property.
        """
        return self._id

    @property
    def name(self):
        """
        Name property getter. Use to get name property.
        """
        return self._name

    @property
    def employees_number(self):
        """
        Employees_number property getter. Use to get employees number value.
        """
        return self._employees_number

    @property
    def balance(self):
        """
        Balance property getter. Use to get balance value.
        """
        return self._balance

    @property
    def assets(self):
        """
        Assets property getter. Use to get assets list.
        """
        return self._assets

    def hire_employee(self):
        """
        To increase company's employees number.
        """
        self._employees_number += 1

    def fire_employee(self):
        """
        To decrease company's employees number.
        """
        # Validation decrease operation.
        assert (self._employees_number - 1) != 0, 'Company can\'t have 0 employees.'
        self._employees_number -= 1

    def buy_asset(self, asset: Asset):
        """
        Asset buy operation.
        :param asset: Asset object.
        """
        # Validation of buy operation.
        assert (self._balance - asset.price) >= 0.0, 'Company can\'t have negative balance.'
        self._assets.append(asset)
        self._balance -= asset.price

    def sell_asset(self, asset: Asset):
        """
        Asset sell operation.
        :param asset: Asset object.
        """
        # Validation of sell operation.
        assert list(filter(lambda a: a.id == asset.id, self._assets)), 'Company can\'t sell assets it doesn\'t have.'
        self._assets.remove(asset)
        self._balance += asset.price

    def __repr__(self):
        """
        Reload __repr__ method for display more actual information.
        """
        return 'Company: id - {id}; name - {name}, employees number - {employees_number}.'.format(
            id=self.id, name=self.name, employees_number=self.employees_number)


# Testing.
if __name__ == '__main__':

    # Create asset objects
    a1 = Asset('Some building 1', 1000000.0)
    print (a1)
    print ('asset id: ' + str(a1.id))
    print ('asset name: ' + a1.name)
    print ('asset price: ' + str(a1.price))
    a2 = Asset('Some car 1', 32000.0)
    print (a2)

    # Create some company.
    print ()
    c1 = Company('Some company 1', 2, 2000000.0)
    print (c1)
    print (c1.balance)
    print (c1.assets)

    # Buy assets to company.
    print ('\nBuy assets.')
    c1.buy_asset(a1)
    print (c1.balance)
    print (c1.assets)
    c1.buy_asset(a2)
    print (c1.balance)
    print (c1.assets)

    # Sell asset.
    print ('\nSell asset')
    c1.sell_asset(a2)
    print (c1.balance)
    print (c1.assets)

    # Add and remove employee.
    print ('\nFire and hire employee.')
    c1.hire_employee()
    print(c1)
    c1.fire_employee()
    print(c1)

    print ('\nTry fire all employee.')
    try:
        c1.fire_employee()
        print(c1)
        c1.fire_employee()
    except Exception as exp:
        print (exp)
        print(c1)
