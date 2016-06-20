# -*- coding: utf-8 -*-
import codecs
import os
import textwrap
from urllib import request, parse

from lxml import html


class WebParser(object):

    def __init__(self):

        self.webpage_url = ''

        self.out_text = ''
        # out_texts = [(text, parts_count), ...]
        self.out_texts = []
        self.out_text_max_density = 0
        self.out_text_density_coeff = 0.33

        self.out_files_dir = 'out_files'

        self.exclude_tag = ['script', 'style']
        self.tags_stack = []

    def setup(self, url='', density_coeff=0, out_dir=''):
        if url:self.set_webpage_url(url)
        if density_coeff: self.set_density_coeff(density_coeff)
        if out_dir: self.set_out_dir(out_dir)

    def set_webpage_url(self, url):
        if url.startswith('http://') or url.startswith('https://'):
            self.webpage_url = url
        else:
            self.webpage_url = ''.join(['http://', url])

    def set_density_coeff(self, density_coeff):
        self.out_text_density_coeff = density_coeff

    def set_out_dir(self, out_dir):
        self.out_files_dir = out_dir

    def webpage_parse(self):

        self.webpage_article_parse()
        self.create_out_file()

    def webpage_article_parse(self):

        tree = self.get_webpage_html_tree()
        root_element = tree.getroot()
        self.tags_stack = []
        self.webpage_article_text_grab(root_element)

    def get_webpage_html_tree(self):

        response = request.urlopen(self.webpage_url)
        charset = response.headers.get_param('charset')
        setup_parser = html.HTMLParser(encoding=charset)
        tree = html.parse(response, parser=setup_parser)

        return tree

    def webpage_article_text_grab(self, element):

        self.article_text_recursive_grabber(element)
        self.choose_best_text()

    def article_text_recursive_grabber(self, element):

        self.tags_stack.append(str(element.tag))

        el_level_text = ''
        el_level_text_part_count = 0

        element_text = self.get_element_text(element)
        el_level_text += ''.join([self.format_text_by_80_chars_in_line(element_text), '\n', ]) if element_text else ''
        el_level_text_part_count += 1 if el_level_text else 0

        for sub_el in element:

            sub_element_text = self.get_sub_element_text(sub_el)
            if sub_element_text:
                el_level_text += '\n' if el_level_text else ''
                el_level_text += ''.join([self.format_text_by_80_chars_in_line(sub_element_text), '\n', ])
                el_level_text_part_count += 1
            else:
                el_level_text += ''

        if el_level_text:
            self.out_texts.append((el_level_text, el_level_text_part_count))

            if len(el_level_text)/el_level_text_part_count > self.out_text_max_density:
                self.out_text_max_density = len(el_level_text)/el_level_text_part_count

        for sub_el in element:
            self.article_text_recursive_grabber(sub_el)

        self.tags_stack.pop()

    def get_element_text(self, element):

        element_text = ''
        if self.is_articlte_page_tag_ok(element.tag):
            element_text += getattr(element, 'text', '') or ''

        return element_text.strip()

    def get_sub_element_text(self, sub_element):

        sub_element_text = ''

        if self.is_articlte_page_tag_ok(sub_element.tag):

            sub_element_text += getattr(sub_element, 'text', '') or ''

            for sub_sub_el in sub_element:
                if sub_sub_el.tag == 'a' and sub_sub_el.text:
                    sub_element_text += getattr(sub_sub_el, 'text', '') or ''
                    sub_element_text += ' [{link}]'.format(link=sub_sub_el.attrib.get('href', ''))
                    sub_element_text += getattr(sub_sub_el, 'tail', '') or ''

        return sub_element_text.strip()

    def is_articlte_page_tag_ok(self, el_tag):
        res = (
            isinstance(el_tag, str) and
            el_tag not in self.exclude_tag and
            self.tags_stack[:2] == ['html', 'body']
        )
        return res

    def choose_best_text(self):
        self.out_text = ''
        for txt, cnt in self.out_texts:
            if (
                txt and
                len(txt)/cnt > self.out_text_max_density * self.out_text_density_coeff and
                len(txt) > len(self.out_text)
            ):
                self.out_text = txt

    @staticmethod
    def format_text_by_80_chars_in_line(text):
        new_text = textwrap.fill(text.strip())
        return new_text

    def create_out_file(self):
        out_file_name = self.get_out_file_name()
        with codecs.open(out_file_name, 'w', encoding='utf-8') as f:
            f.write(self.out_text)

    def get_out_file_name(self):
        url_parse_res = parse.urlparse(self.webpage_url)
        unormalize_file_dir_path = url_parse_res.hostname + url_parse_res.path
        normalize_file_dir_path = request.pathname2url(unormalize_file_dir_path)
        out_file_dir_path = os.path.join(self.out_files_dir, normalize_file_dir_path)
        if not os.path.isdir(out_file_dir_path):
            os.makedirs(out_file_dir_path)
        file_name = os.path.join(out_file_dir_path, 'parse_result.txt')
        return file_name
