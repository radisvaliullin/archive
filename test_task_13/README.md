# test_task_13
Some test task on Python

## Requirements
python3.4.x or later
BeautifulSoup4
lxml

## Run
python3 delstyle.py test.address

## Task Description
Имеется ссылка на HTML страницу,
Нужно написать приложение на Python, которое получает в качестве параметра
командной строки эту ссылку.
Далее приложение должно получить содержимое HTML документа по ссылке и
во всех тэгах убрать атрибуты style.
Результат вывести на экран.

Например, приложению передается ссылка http://test.address
Приложение получает HTML документ следующего содержания
```
<body>
  <p style="font-size: 12pt">Example</p>
</body>
```

Результат должен быть такой:
```
<body>
  <p>Example</p>
</body>
```