# -*- coding: utf-8 -*-
import argparse

from webparsertolls.webparsertolls import WebParser

input_args_parser = argparse.ArgumentParser(description='Simple web parser.')
input_args_parser.add_argument('url')
args = input_args_parser.parse_args()

web_parser = WebParser(args.url)
web_parser.webpage_parse()
web_parser.create_out_file()


# webparser.simple_parse()
# webparser.create_out_file()

# sites = [
#     'lenta.ru', 'gazeta.ru', 'meduza.io', 'news.yandex.ru', 'f1news.ru',
# ]
#
# for page in sites:
#
#     parser = WebParser(page)
#
#     parser.simple_parse_v2()
#     parser.create_out_file('.'.join([
#         u''.join(['tree__', page.replace('.', '_')]),
#         'txt'
#     ]))
#
#     parser.out_text = u''
#     parser.simple_parse_v4()
#     parser.create_out_file('.'.join([page.replace('.', '_'), 'txt']))
