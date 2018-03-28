# -*- coding: utf-8 -*-


from aparser.coinscatalog import coins_catalog_statistic

if __name__ == '__main__':

    coins_unknown, coins_ratio = coins_catalog_statistic()
    print(coins_unknown, coins_ratio)
