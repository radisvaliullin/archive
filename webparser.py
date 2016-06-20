# -*- coding: utf-8 -*-
import argparse
import configparser
import json

from webparsertolls.webparsertolls import WebParser


args_parser = argparse.ArgumentParser(description='Simple web arcticle parser.')
args_parser.add_argument('url', nargs='?', default='', help='url to article web page.')
args_parser.add_argument('-p', help='Density coefficient.')
args_parser.add_argument('-o', help='Output file folder.')
args = args_parser.parse_args()


config_parser = configparser.ConfigParser()
config_parser.read('default.config.ini')


if __name__ == '__main__':

    if args.url:
        urls = [args.url]
    else:
        urls = config_parser['URLS']['TEST_URLS'].split()

    web_parser = WebParser()
    web_parser.setup(
        density_coeff=float(config_parser['DEFAULT']['DENSITY']),
        out_dir=config_parser['DEFAULT']['OUT_PATH']
    )

    for url in urls:
        web_parser.set_webpage_url(url)
        web_parser.webpage_parse()
