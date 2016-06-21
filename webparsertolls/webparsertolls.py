# -*- coding: utf-8 -*-
import codecs
import os
import textwrap
from urllib import request, parse

from lxml import html


class WebParser(object):

    def __init__(self):

        self.webpage_url = ''

        self.element_text_search_levels = 2
        self.out_text = ''
        # out_texts = [(text, parts_count), ...]
        self.out_texts = []
        self.out_texts_max_density = 0
        self.out_text_density_coeff = 0.33

        self.out_files_dir = 'out_files'

        self.exclude_tag = ['script', 'style']
        self.tags_stack = []

    def setup(self, url='', density_coeff=0, out_dir='', text_search_levels=0):
        if url:self.set_webpage_url(url)
        if density_coeff: self.set_density_coeff(density_coeff)
        if out_dir: self.set_out_dir(out_dir)
        if text_search_levels: self.set_text_search_levels(text_search_levels)

    def set_webpage_url(self, url):
        if url.startswith('http://') or url.startswith('https://'):
            self.webpage_url = url
        else:
            self.webpage_url = ''.join(['http://', url])

    def set_density_coeff(self, density_coeff):
        self.out_text_density_coeff = density_coeff

    def set_out_dir(self, out_dir):
        self.out_files_dir = out_dir

    def set_text_search_levels(self, text_search_levels):
        self.element_text_search_levels = text_search_levels

    def webpage_parse(self):

        self.webpage_article_parse()
        self.create_out_file()

    def webpage_article_parse(self):

        tree = self.get_webpage_html_tree()
        root_element = tree.getroot()
        self.out_text = ''
        self.out_texts = []
        self.out_texts_max_density = 0
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
        self.choose_best_out_text()
        self.formating_out_text()

    def article_text_recursive_grabber(self, element):

        self.tags_stack.append(str(element.tag))

        element_level_text = self.get_element_level_text(element)

        if element_level_text:

            el_text_parts_count = len(element_level_text.split('\n'))
            el_text_density = len(element_level_text)/el_text_parts_count
            self.out_texts.append((
                element_level_text,
                el_text_parts_count,
                el_text_density,
            ))

            if el_text_density > self.out_texts_max_density:
                self.out_texts_max_density = el_text_density

        for sub_el in element:
            self.article_text_recursive_grabber(sub_el)

        self.tags_stack.pop()

    def get_element_level_text(self, element, level=0):

        self.tags_stack.append(str(element.tag))

        element_level_text = self.get_element_text(element)

        if level <= self.element_text_search_levels:

            for sub_el in element:

                sub_element_text = self.get_element_level_text(sub_el, level=level+1)

                element_level_text += sub_element_text

        self.tags_stack.pop()

        return element_level_text

    def get_element_text(self, element):

        element_text = ''
        text_tag_raw = getattr(element, 'text', '') or ''
        tail_tag_raw = getattr(element, 'tail', '') or ''
        text_tag = text_tag_raw.strip()
        tail_tag = tail_tag_raw.strip()

        if self.is_articlte_page_tag_ok(element.tag):

            if tail_tag:
                element_text += ' ' + text_tag if text_tag else ''
                if hasattr(element, 'attrib') and element.attrib.get('href', ''):
                    element_text += ' [{link}]'.format(link=element.attrib.get('href', ''))
                element_text += ' ' + tail_tag
            else:
                element_text += '\n' + text_tag if text_tag else ''

        return element_text

    def is_articlte_page_tag_ok(self, el_tag):
        res = (
            isinstance(el_tag, str) and
            el_tag not in self.exclude_tag and
            self.tags_stack[:2] == ['html', 'body']
        )
        return res

    def choose_best_out_text(self):
        self.out_text = ''
        for txt, cnt, den in self.out_texts:
            if (
                txt and
                # den > self.out_texts_max_density * self.out_text_density_coeff and
                len(txt) > len(self.out_text)
            ):
                self.out_text = txt

    def formating_out_text(self):
        self.out_text = '\n\n'.join([
            self.format_text_by_80_chars_in_line(text_part)
            for text_part in self.out_text.split('\n')
        ])
        self.out_text += '\n'

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
