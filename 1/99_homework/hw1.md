Домашнее задание №1

Срок выполнения - 28 октября 2017 (РК1)

----

bugs.

Несколько багов и неочевидных моментов го. Про всех из них можно прочитать в https://habrahabr.ru/company/mailru/blog/314804/. Для разминки.

Запускать тесты через `go test -v` находясь в папке `bugs`.

----

Утилита tree.

Выводит дерево каталогов и файлов (если указана опция -f).

Необходимо реализовать функцию `dirTree` внутри `main.go`. Начать можно с с `https://golang.org/pkg/os/#Open` и дальше смотреть какие методы есть у результата.

Обращаю ваше внимание на строку `├───main.go (vary)`. Поскольку размер меняется меняется - придётся вкрутить небольшой костыль для этого файла на 1-м уровне.

Запускать тесты через `go test -v` находясь в папке `tree`. После запуска вы должны увидеть такой результат:
```
$ go test -v
=== RUN   TestTreeFull
--- PASS: TestTreeFull (0.00s)
=== RUN   TestTreeDir
--- PASS: TestTreeDir (0.00s)
PASS
ok      gitlab.com/rvasily/golang-2017-2/1/99_homework/tree     0.127s
```

```
go run main.go . -f
├───main.go (vary)
├───main_test.go (1318b)
└───testdata
	├───project
	│	├───file.txt (19b)
	│	└───gopher.png (70372b)
	├───static
	│	├───css
	│	│	└───body.css (28b)
	│	├───html
	│	│	└───index.html (57b)
	│	└───js
	│		└───site.js (10b)
	├───zline
	│	└───empty.txt (empty)
	└───zzfile.txt (empty)
```

```
go run main.go .
└───testdata
	├───project
	├───static
	│	├───css
	│	├───html
	│	└───js
	└───zline
```

Рекомендуется к прочтению https://habrahabr.ru/post/306914/