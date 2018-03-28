# test_task_6
Some test task on Golang.

# Solution description.
1 Get html tree for page with all analyzes.
2 Get and parsing analyzes node.
3 For each analyse run new goroutine witch parsing analyse's detail page.
4 Parsed analyse with details putting to chan.
5 Special goroutine take analyse from chan and put to DB.

# Task description.
ТЗ “Инвитро”

Сайта “Инвитро.ру”, вкладка “Анализы и цены” (https://www.invitro.ru/analizes/for-doctors/)

Задание – написать парсер этого раздела.
В разделе 4 уровня вкладок

1 Уровень – Тип исследований – Гематологические исследования, Биохимические исследования и т.д.

2 Уровень – Подтип исследования, например, для типа Гематологические исследования подтипами будут “Клинический анализ крови, Имуногематологические исследование, Коагулологическое исследование”
У некоторых типов подтип отсутствует

3 Уровень – Вид исследования – номер и стоимость сохранять не требуется.
В некоторых ветках уровня 3 таблица содержат подзаголовки. Их можно пропустить.

4 Уровень – Описание исследование. Здесь необходимо получить содержание четырех вкладок – Описание, Подготовка, Показания, Интерпретация результатов (в некоторых ветках имеются только две вкладки)

Результаты необходимо сохронить в базу так, чтобы для каждого исследования можно было установить его тип и подтип.
Процесс сбора информации на каждом уровне должен  быть распараллелен Go рутинами .