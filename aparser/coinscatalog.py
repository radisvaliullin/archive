# -*- coding: utf-8 -*-
import random
import re
import urllib.request
from html.parser import HTMLParser
from urllib.parse import urljoin

import time

from aparser import tree


def coins_catalog_statistic():
    """ Парсит информацию по каталогу монет """

    coins_unknown_year = 0
    coins_early_2000 = 0
    coins_later_2000 = 0
    coins_ratio = 0

    parser = CatalogHTMLParser()
    cl, ce, cu, next_page = parser.parse_coins_years_info()
    coins_unknown_year += cu
    coins_early_2000 += ce
    coins_later_2000 += cl

    # Лимит на количество обрабатываемых страниц
    catalog_page_limit = 4
    page_cnt = 1
    while next_page:
        if page_cnt == catalog_page_limit:
            break
        # print(next_page)

        parser = CatalogHTMLParser(next_page)
        cl, ce, cu, next_page = parser.parse_coins_years_info()
        coins_unknown_year += cu
        coins_early_2000 += ce
        coins_later_2000 += cl
        page_cnt += 1


    if coins_later_2000:
        coins_ratio = coins_early_2000/coins_later_2000

    return coins_unknown_year, coins_ratio


class CatalogHTMLParser(HTMLParser):
    """ 
    Парсит страницу с каталогом монет. Возвращает информацию по годам и следующею страницу для парсинга. 
    """

    domain = 'https://www.avito.ru'
    catalog_url = 'moskva/kollektsionirovanie/monety'
    next_catalog_page_url = ''
    next_page_trigger = False

    coins_unknown_year = 0
    coins_early_2000 = 0
    coins_later_2000 = 0

    # Задаем структуру веток дерева каталога
    catalog_before_branch_struct = [
        {'tag': 'div', 'attrs': {'class': 'js-catalog_before-ads'}, },
    ]
    catalog_after_branch_struct = [
        {'tag': 'div', 'attrs': {'class': 'js-catalog_after-ads'}, },
    ]
    catalog_item_branch_struct = [
        # {'tag': 'div', 'attrs': {'class': 'item item_table clearfix js-catalog-item-enum item-highlight'}, },
        # {'tag': 'div', 'attrs': {}, },
        {'tag': 'div', 'attrs': {'class': 'description'}, },
        {'tag': 'h3', 'attrs': {'class': 'title item-description-title'}, },
        {'tag': 'a', 'attrs': {'class': 'item-description-title-link'}, },
    ]

    # Структура веток дерева пагинации
    catalog_pages_branch_struct = [
        {'tag': 'div', 'attrs': {'class': 'pagination-pages clearfix'}, },
    ]
    catalog_curr_page_branch_struct = [
        {'tag': 'span', 'attrs': {'class': 'pagination-page pagination-page_current'}, },
    ]
    catalog_page_branch_struct = [
        {'tag': 'a', 'attrs': {'class': 'pagination-page'}, },
    ]

    def __init__(self, catalog_usl=''):
        super(CatalogHTMLParser, self).__init__()

        # Инициализируем парсер веток
        self.catalog_before = tree.Branch(self.catalog_before_branch_struct)
        self.catalog_before_item = tree.Branch(self.catalog_item_branch_struct, self.catalog_before)
        self.catalog_after = tree.Branch(self.catalog_after_branch_struct)
        self.catalog_after_item = tree.Branch(self.catalog_item_branch_struct, self.catalog_after)
        self.catalog_pages = tree.Branch(self.catalog_pages_branch_struct)
        self.catalog_curr_page = tree.Branch(self.catalog_curr_page_branch_struct, self.catalog_pages)
        self.catalog_page = tree.Branch(self.catalog_page_branch_struct, self.catalog_pages)

        if catalog_usl:
            self.catalog_url = catalog_usl

    def parse_coins_years_info(self):

        # Добавляем задержку что бы избежать подозрение на робата
        # TODO: (не всегда срабатывает, нужна более лучшая стратегия).
        time.sleep(0.5 + random.random())

        # запрашиваем первую страницу
        open_url = urljoin(self.domain, self.catalog_url)
        resp = urllib.request.urlopen(open_url)

        # определяем кодировку
        charset = resp.info().get_content_charset()

        # декодируем согласно кодировке
        html = resp.read().decode(charset)

        # тестовый принт html
        # print(html)

        # Запускаем разбор страницы
        self.feed(html)

        return self.coins_later_2000, self.coins_early_2000, self.coins_unknown_year, self.next_catalog_page_url

    def handle_starttag(self, tag, attrs):

        self.catalog_before.handle_starttag(tag, attrs)
        self.catalog_before_item.handle_starttag(tag, attrs)
        self.catalog_after.handle_starttag(tag, attrs)
        self.catalog_after_item.handle_starttag(tag, attrs)
        self.catalog_pages.handle_starttag(tag, attrs)
        self.catalog_curr_page.handle_starttag(tag, attrs)
        self.catalog_page.handle_starttag(tag, attrs)

        # Если это элемент каталога то вычисляем url и парсим страницу элемента
        if self.catalog_before_item.is_branch: # or self.catalog_after_item.is_branch:
            for attr, attr_val in attrs:
                if attr == 'href':
                    # for test
                    # print('URL ', attr_val)

                    item_parser = CatalogItemParser(attr_val)
                    coins_year = item_parser.parse_coins_year()
                    if coins_year:
                        if coins_year >= 2000:
                            self.coins_later_2000 += 1
                        else:
                            self.coins_early_2000 += 1
                    else:
                        self.coins_unknown_year += 1

                    # for test
                    # print('Coins', self.coins_later_2000, self.coins_early_2000, self.coins_unknown_year)

        # Если элемент пагинации текущая страница то выставляем тригер следующей страницы
        if self.catalog_curr_page.is_branch:
            self.next_page_trigger = True

        # Определяем следующею страницу
        if self.catalog_page.is_branch and self.next_page_trigger:
            self.next_page_trigger = False
            for attr, attr_val in attrs:
                if attr == 'href':
                    self.next_catalog_page_url = attr_val

    def handle_endtag(self, tag):

        self.catalog_before.handle_endtag(tag)
        self.catalog_before_item.handle_endtag(tag)
        self.catalog_after.handle_endtag(tag)
        self.catalog_after_item.handle_endtag(tag)
        self.catalog_pages.handle_endtag(tag)
        self.catalog_curr_page.handle_endtag(tag)
        self.catalog_page.handle_endtag(tag)

    def handle_data(self, data):
        pass


class CatalogItemParser(HTMLParser):
    """
    Парсер страницы для элемента каталога
    """

    domain_url = 'https://www.avito.ru'
    catalog_item_url = ''

    coins_year = 0

    item_branch_struct = [
        {'tag': 'div', 'attrs': {'class': 'item-view'}, },
        {'tag': 'div', 'attrs': {'class': 'item-view-content'}, },
        {'tag': 'div', 'attrs': {'class': 'item-view-left'}, },
    ]

    item_title_branch_struct = [
        {'tag': 'div', 'attrs': {'class': 'item-view-title-info js-item-view-title-info'}, },
        {'tag': 'div', 'attrs': {'class': 'title-info title-info_mode-with-favorite'}, },
        {'tag': 'div', 'attrs': {'class': 'title-info-main'}, },
        {'tag': 'h1', 'attrs': {'class': 'title-info-title'}, },
        {'tag': 'span', 'attrs': {'class': 'title-info-title-text'}, },
    ]

    item_description_branch_struct = [
        {'tag': 'div', 'attrs': {'class': 'item-view-main js-item-view-main'}, },
        {'tag': 'div', 'attrs': {'class': 'item-view-block'}, },
        {'tag': 'div', 'attrs': {'class': 'item-description'}, },
        {'tag': 'div', 'attrs': {'class': 'item-description-text'}, },
        {'tag': 'p', 'attrs': {}, },
    ]

    item_title_and_description = ''

    def __init__(self, item_url):
        super(CatalogItemParser, self).__init__()

        self.catalog_item_url = item_url
        self.item_branch = tree.Branch(self.item_branch_struct)
        self.item_title_branch = tree.Branch(self.item_title_branch_struct, self.item_branch)
        self.item_description_branch = tree.Branch(self.item_description_branch_struct, self.item_branch)

    def parse_coins_year(self):

        # Добавляем задержку что бы избежать подозрение на робата.
        time.sleep(0.5 + random.random())

        open_url = urljoin(self.domain_url, self.catalog_item_url)

        resp = urllib.request.urlopen(open_url)

        charset = resp.info().get_content_charset()

        html = resp.read().decode(charset)

        self.feed(html)

        # for test
        # print(self.item_title_and_description)

        # Находим дату монеты в описании элемента католога монет
        self.find_coins_year()

        return self.coins_year

    def handle_starttag(self, tag, attrs):

        self.item_branch.handle_starttag(tag, attrs)
        self.item_title_branch.handle_starttag(tag, attrs)
        self.item_description_branch.handle_starttag(tag, attrs)

    def handle_endtag(self, tag):

        self.item_branch.handle_endtag(tag)
        self.item_title_branch.handle_endtag(tag)
        self.item_description_branch.handle_endtag(tag)

    def handle_data(self, data):
        # Если заголовок или описание то добавляем в общею переменную для последующего анализа по годам
        if self.item_title_branch.is_branch or self.item_description_branch.is_branch:
            self.item_title_and_description = ' '.join([self.item_title_and_description, data])

    def find_coins_year(self):
        """
        Простой поиск даты через регулярное выражение 
        """
        date_exp = re.compile('(\d\d\d\d)')
        date_match = date_exp.search(self.item_title_and_description)
        if date_match:
            self.coins_year = int(date_match.group())
