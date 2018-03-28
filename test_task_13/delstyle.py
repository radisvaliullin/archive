#!/urs/bin/env python
# delstyle - is a simple utility,
#  helps you to get html page without style args.
import argparse
import urllib.request

from bs4 import BeautifulSoup


# if execute (not import)
if __name__ == '__main__':

    # parse args
    parser = argparse.ArgumentParser(
        description='Print html page without style')
    parser.add_argument('url', type=str, help='print page url')
    args = parser.parse_args()
    print(args.url)

    # if need add http prefix
    url = args.url if args.url.startswith('http') else 'http://' + args.url

    # open web page
    try:
        res = urllib.request.urlopen(url)

        # decode and parse html with lxml lib
        soup = BeautifulSoup(res, 'lxml')

        # delete style tag
        for tag in soup.recursiveChildGenerator():
            if hasattr(tag, 'attrs'):
                tag.attrs = {
                    key: value for key, value in tag.attrs.items()
                    if key != 'style'
                }

        # print result
        print(soup)

    except Exception:
        print("can't request url")
