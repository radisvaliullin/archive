# -*- coding: utf-8 -*-
import argparse

from webparsertolls.webparsertolls import WebParser

input_args_parser = argparse.ArgumentParser(description=u'Simple web parser.')
input_args_parser.add_argument(u'url')
args = input_args_parser.parse_args()

web_parser = WebParser(args.url)
web_parser.webpage_parse()
web_parser.create_out_file()


# webparser.simple_parse()
# webparser.create_out_file()

# sites = [
#     u'lenta.ru', u'gazeta.ru', u'meduza.io', u'news.yandex.ru', u'f1news.ru',
# ]
#
# for page in sites:
#
#     parser = WebParser(page)
#
#     parser.simple_parse_v2()
#     parser.create_out_file(u'.'.join([
#         u''.join([u'tree__', page.replace(u'.', u'_')]),
#         u'txt'
#     ]))
#
#     parser.out_text = u''
#     parser.simple_parse_v4()
#     parser.create_out_file(u'.'.join([page.replace(u'.', u'_'), u'txt']))
