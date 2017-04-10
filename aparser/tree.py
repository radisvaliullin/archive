# -*- coding: utf-8 -*-


class Branch:
    """ 
    Branch - ветка, описывает элемент дерева в виде последовательности тегов.
    При поиске по дереву тегов html, находит указанную последовательность (branch) и выставляет True в is_branch.
    """

    parent_branch = None

    is_branch = False

    # Описываем структуру (последовательность элементов) в ветке
    # [{'tag': 'div', 'attrs': {'class': 'item-view'}, }, ...]
    branch = []
    branch_elem_idx = 0
    branch_tag_stack = []
    stackable_tag = set()
    branch_elem_stack_pos = []

    def __init__(self, branch_struct, parent_branch=None):

        self.branch = branch_struct
        self.parent_branch = parent_branch
        for branch_elem in branch_struct:
            self.stackable_tag.add(branch_elem.get('tag'))

    def handle_starttag(self, tag, attrs):

        # Добавляем в стэк (первый элемент добавляется отдельно при достижении первого элемента ветки)
        if self.branch_tag_stack and tag in self.stackable_tag:
            self.branch_tag_stack.append(tag)

        if not self.is_branch and (self.parent_branch is None or self.parent_branch.is_branch):

            # Текущий элемент
            branch_elem = self.branch[self.branch_elem_idx]

            current_attrs = dict(attrs)
            expected_attrs = branch_elem.get('attrs')
            if tag == branch_elem.get('tag') and self.is_tag_attrs_true(current_attrs, expected_attrs):

                # Если это первый элемент
                if self.branch_elem_idx == 0:
                    self.branch_tag_stack = [tag]
                    self.branch_elem_stack_pos = [0]
                else:
                    self.branch_elem_stack_pos.append(len(self.branch_tag_stack) - 1)

                if (self.branch_elem_idx + 1) == len(self.branch):
                    self.is_branch = True
                else:
                    self.branch_elem_idx += 1

    def handle_endtag(self, tag):

        if self.branch_tag_stack:
            if tag in self.stackable_tag:
                self.branch_tag_stack.pop()

            last_elem_stack_pos = self.branch_elem_stack_pos[len(self.branch_elem_stack_pos) - 1]

            if (len(self.branch_tag_stack) - 1) < last_elem_stack_pos:
                if self.branch_elem_idx != 0:
                    self.branch_elem_idx -= 1
                self.branch_elem_stack_pos.pop()

            if self.is_branch and ((self.branch_elem_idx + 1) != len(self.branch) or not self.branch_tag_stack):
                self.is_branch = False

    def is_tag_attrs_true(self, curr_attrs, expected_attrs):
        attr_true = True
        for attr, attr_val in expected_attrs.items():
            if curr_attrs.get(attr) != attr_val:
                attr_true = False
                break
        return attr_true
