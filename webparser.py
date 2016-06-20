# -*- coding: utf-8 -*-
import argparse

from webparsertolls.webparsertolls import WebParser

input_args_parser = argparse.ArgumentParser(description='Simple web parser.')
input_args_parser.add_argument('url')
args = input_args_parser.parse_args()

web_parser = WebParser(args.url)
web_parser.webpage_parse()
web_parser.create_out_file()


sites = [
    'https://lenta.ru/news/2016/06/20/baribal_attacks/',
    'https://lenta.ru/news/2016/06/20/batman/',
    'http://www.gazeta.ru/politics/news/2016/06/20/n_8783951.shtml',

    'http://www.f1news.ru/news/f1-113012.html',
    'http://www.f1news.ru/Championship/2016/europe/race.shtml',

    # 'https://meduza.io/feature/2016/06/20/izgnanie-islamskogo-gosudarstva',
    # 'https://meduza.io/news/2016/06/20/oon-soobschila-o-rekordnom-chisle-bezhentsev-v-2015-godu',

    'https://news.yandex.ru/yandsearch?cl4url=www.kommersant.ru%2Fdoc%2F3018006&lr=43&lang=ru&rubric=index',

    'http://matchtv.ru/news/uefa-ne-soglasoval-zapros-rossii-na-traurnye-povyazki-vo-vremya-matcha-protiv-uelsa/',

]

for page in sites:

    web_parser = WebParser(page)
    web_parser.webpage_parse()
    web_parser.create_out_file()
