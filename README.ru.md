<p align="center">
  <a href="README.md#apisite">English</a> |
  <span>Pусский</span>
</p>

---

# apisite
> Платформа для веб-приложений, основанных на API

[![GoCard][gc1]][gc2]
 [![GitHub Release][gr1]][gr2]
 [![GitHub code size in bytes][sz]]()
 [![GitHub license][gl1]][gl2]

[gc1]: https://goreportcard.com/badge/apisite/apisite
[gc2]: https://goreportcard.com/report/github.com/apisite/apisite
[gr1]: https://img.shields.io/github/release/apisite/apisite.svg
[gr2]: https://github.com/apisite/apisite/releases
[sz]: https://img.shields.io/github/languages/code-size/apisite/apisite.svg
[gl1]: https://img.shields.io/github/license/apisite/apisite.svg
[gl2]: LICENSE

Проект [apisite](https://github.com/apisite/apisite) является примером реализации архитектуры, в которой бизнес-логика веб-приложения размещается в хранимом коде БД. Такая архитектура имеет следующие особенности:

* внешние (javascript) и внутренние (шаблоны HTML-страниц) потребители могут использовать общий комплект методов API
* наличие в БД метаданных хранимого кода позволяет документировать API на лету без парсинга или кодогенерации
* пакет изменений бизнес-логики может быть протестирован и загружен в приложение в рамках одной транзакции
* для добавления страницы на сайт достаточно скопировать файл с шаблоном и зарегистрировать его в БД (в подсистеме разграничения доступа)
* шлюз между потребителями и БД может быть независим от бизнес-логики и содержать только сервисные функции
* т.к. для потребителей не нужна строгая типизация данных, шлюз может ограничиться поддержкой только базовых типов значений

## Структура проекта

* [procapi](https://github.com/apisite/procapi) - библиотека доступа к хранимым функциям
* [apitpl](https://github.com/apisite/apitpl) - библиотека для шаблонов, использующих API
* [apisite](https://github.com/apisite/apisite) - приложение, объединяющее библиотеки procapi и apitpl
* [pomasql](https://github.com/pomasql/poma) - среда загрузки SQL-кода в БД Postgresql

## Пример приложения

* [app-enfist](https://github.com/apisite/app-enfist) - интерфейс управления настройками приложений проекта [dcape](https://github.com/dopos/dcape)

## Предыдущие версии

* [pgws](https://github.com/LeKovr/pgws) - первичная реализация архитектуры (perl, 2010)
* [dbrpc](https://github.com/LeKovr/dbrpc) - предварительная реализация JSON-RPC интерфейса (go, 2016) 

## License

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2018 Aleksei Kovrizhkin <lekovr+apisite@gmail.com>

