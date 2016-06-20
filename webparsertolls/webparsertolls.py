# -*- coding: utf-8 -*-
import codecs
import os
import textwrap
import urllib

from lxml import html


class WebParser(object):

    def __init__(self, url):

        if url.startswith('http://') or url.startswith('https://'):
            self.webpage_url = url
        else:
            self.webpage_url = ''.join(['http://', url])

        self.out_text = ''

        self.exclude_tag = ['script', 'style']
        self.tags_stack = []

    def webpage_parse(self):

        self.webpage_article_parse()

    def webpage_article_parse(self):

        tree = self.get_webpage_html_tree()
        root_element = tree.getroot()
        self.tags_stack = []
        self.webpage_article_text_grab(root_element)

    def get_webpage_html_tree(self):

        response = urllib.request.urlopen(self.webpage_url)
        charset = response.headers.get_param('charset')
        setup_parser = html.HTMLParser(encoding=charset)
        tree = html.parse(response, parser=setup_parser)

        return tree

    def webpage_article_text_grab(self, element):

        self.article_text_recursive_grabber(element)

    def article_text_recursive_grabber(self, element):

        self.tags_stack.append(str(element.tag))

        el_level_text = ''

        element_text = self.get_element_text(element)
        el_level_text += ''.join([self.format_text_by_80_chars_in_line(element_text), '\n', ]) if element_text else ''

        for sub_el in element:

            sub_element_text = self.get_sub_element_text(sub_el)
            if sub_element_text:
                el_level_text += '\n' if el_level_text else ''
                el_level_text += ''.join([self.format_text_by_80_chars_in_line(sub_element_text), '\n', ])
            else:
                el_level_text += ''

        if len(el_level_text) > len(self.out_text):
            self.out_text = el_level_text

        for sub_el in element:
            self.article_text_recursive_grabber(sub_el)

        self.tags_stack.pop()

    def get_element_text(self, element):

        element_text = ''
        if self.is_articlte_page_tag_ok(element.tag):
            element_text += getattr(element, 'text', '') or ''

        return element_text

    def get_sub_element_text(self, sub_element):

        sub_element_text = ''

        if self.is_articlte_page_tag_ok(sub_element.tag):

            sub_element_text += getattr(sub_element, 'text', '') or ''

            for sub_sub_el in sub_element:
                if sub_sub_el.tag == 'a' and sub_sub_el.text:
                    sub_element_text += getattr(sub_sub_el, 'text', '') or ''
                    sub_element_text += ' [{link}]'.format(link=sub_sub_el.attrib.get('href', ''))
                    sub_element_text += getattr(sub_sub_el, 'tail', '') or ''

        return sub_element_text

    def is_articlte_page_tag_ok(self, el_tag):
        res = (
            isinstance(el_tag, str) and
            el_tag not in self.exclude_tag and
            self.tags_stack[:2] == ['html', 'body']
        )
        return res

    @staticmethod
    def format_text_by_80_chars_in_line(text):
        new_text = textwrap.fill(text.strip())
        return new_text

    def create_out_file(self):
        out_file_name = self.get_out_file_name()
        with codecs.open(out_file_name, 'w', encoding='utf-8') as f:
            f.write(self.out_text)

    def get_out_file_name(self):
        url_parse_res = urllib.parse.urlparse(self.webpage_url)
        unormalize_file_dir_path = url_parse_res.hostname + url_parse_res.path
        normalize_file_dir_path = urllib.request.pathname2url(unormalize_file_dir_path)
        out_file_dir_path = os.path.join('out_files', normalize_file_dir_path)
        if not os.path.isdir(out_file_dir_path):
            os.makedirs(out_file_dir_path)
        file_name = os.path.join(out_file_dir_path, 'parse_result.txt')
        return file_name
