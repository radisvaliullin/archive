# -*- coding: utf-8 -*-
import codecs
import os
import textwrap
import urllib
import urllib2
from urlparse import urlparse

from lxml import etree, html


class WebParser(object):

    def __init__(self, url):
        if url.startswith(u'http://') or url.startswith(u'https://'):
            self.webpage_url = url
        else:
            self.webpage_url = u''.join([u'http://', url])

        self.out_text = u''

        self.special_tag = []
        self.enable_tag = []
        self.exclude_tag = [u'script']
        self.tags_stack = []

    def webpage_parse(self):
        url_parse_res = urlparse(self.webpage_url)
        if url_parse_res.path:
            self.webpage_article_parse()
        else:
            self.webpage_main_parse()

    def webpage_main_parse(self):
        tree = self.get_webpage_html_tree()
        root_element = tree.getroot()
        self.tags_stack = []
        self.recursive_text_graber(root_element)

    def webpage_article_parse(self):
        tree = self.get_webpage_html_tree()
        root_element = tree.getroot()
        self.tags_stack = []
        self.webpage_article_text_grab(root_element)

    def get_webpage_html_source(self):
        response = urllib2.urlopen(self.webpage_url)
        charset = response.headers.getparam('charset')
        webpage_html_source = unicode(response.read(), encoding=charset)
        return webpage_html_source

    def get_webpage_html_tree(self):
        response = urllib2.urlopen(self.webpage_url)
        charset = response.headers.getparam('charset')
        setup_parser = html.HTMLParser(encoding=charset)
        tree = html.parse(response, parser=setup_parser)
        return tree

    def recursive_text_graber(self, element):

        self.tags_stack.append(unicode(element.tag))

        if self.is_tag_ok(element.tag):

            if self.is_el_text_ok(element):

                el_text = unicode(getattr(element, u'text', u'') or u'')
                self.out_text += u'\n'
                self.out_text += u''.join([
                    self.format_text_by_80_chars_in_line(el_text), u'\n',
                ])

        for sub_el in element:
            self.recursive_text_graber(sub_el)

        self.tags_stack.pop()

    def webpage_article_text_grab(self, element):

        self.tags_stack.append(unicode(element.tag))

        if self.is_art_pg_tag_ok(element.tag):

            if self.is_el_text_ok(element):

                el_text = unicode(getattr(element, u'text', u'') or u'')
                self.out_text += u'\n'
                self.out_text += u''.join([
                    self.format_text_by_80_chars_in_line(el_text), u'\n',
                ])

        for sub_el in element:
            self.recursive_text_graber(sub_el)

        self.tags_stack.pop()

    def is_tag_ok(self, el_tag):
        res = (
            isinstance(el_tag, basestring) and
            el_tag not in self.exclude_tag and
            (
                self.tags_stack[:2] != [u'html', u'head'] or
                (
                    self.tags_stack[:2] == [u'html', u'head'] and
                    el_tag == u'title'
                )
            )
        )
        return res

    def is_art_pg_tag_ok(self, el_tag):
        res = (
            isinstance(el_tag, basestring) and
            el_tag not in self.exclude_tag and
            self.tags_stack[:2] == [u'html', u'body']
        )
        return res

    def is_el_text_ok(self, element):
        el_text = unicode(getattr(element, u'text', u'') or u'')
        res = (
            len(el_text.strip().split()) >= 4 or
            self.is_sub_elements_text_len_ok(element)
        )
        return res

    def is_sub_elements_text_len_ok(self, element):
        el_text = unicode(getattr(element, u'text', u'') or u'')
        sub_texts = u''
        sub_texts_count = 0
        for el in element:
            if self.is_tag_ok(el.tag):
                sub_el_text = unicode(getattr(el, u'text', u'') or u'')
                sub_el_text = sub_el_text.strip()
                sub_texts += u' '
                sub_texts += sub_el_text
                if sub_el_text:
                    sub_texts_count += 1
        res = (
            len(el_text.strip().split()) >= 1 and
            len(sub_texts.strip().split()) > 5 and
            sub_texts_count >= 1
        )
        return res

    @staticmethod
    def format_text_by_80_chars_in_line(text):
        new_text = textwrap.fill(text.strip())
        return new_text

    def create_out_file(self):
        out_file_name = self.get_out_file_name()
        with codecs.open(out_file_name, u'w', encoding=u'utf-8') as f:
            f.write(self.out_text)

    def get_out_file_name(self):
        url_parse_res = urlparse(self.webpage_url)
        unormalize_file_dir_path = url_parse_res.hostname + url_parse_res.path
        normalize_file_dir_path = urllib.pathname2url(unormalize_file_dir_path)
        out_file_dir_path = os.path.join(u'out_files', normalize_file_dir_path)
        if not os.path.isdir(out_file_dir_path):
            os.makedirs(out_file_dir_path)
        file_name = os.path.join(out_file_dir_path, u'parse_result.txt')
        return file_name
