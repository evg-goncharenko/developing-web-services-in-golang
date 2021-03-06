## Обработка запросов и написание тестов
Проект представляет собой реализацию отправки запросов, получание ответов, работа с параметрами, хедерами, а так же написание тестов.

У нас есть какой-то поисковый сервис:
- `SearchClient` - структура с методом `FindUsers`, который отправляет запрос во внешнюю систему и возвращает результат, немного преобразуя его.
- `SearchServer` - своего рода внешняя система. Непосредственно занимается поиском данных в файле `dataset.xml`. В продакшене бы запускалась в виде отдельного веб-сервиса.

Структура:
- Функция `SearchServer()` в файле `server.go`, который запускается через тестовый сервер.
- Покрытие тестами метод FindUsers и SearchServer в `client_test.go` (покрытие 100%). 
- Генерация html-отчета с покрытием в  `cover.html`.
- Данные для работы лежат в файле `dataset.xml`
- Параметр `query` ищет по полям `Name` и `About`
- Параметр `order_field` работает по полям `Id`, `Age`, `Name`, если пустой - то возвращает по `Name`, если что-то другое - SearchServer ругается ошибкой. `Name` - это first_name + last_name из XML.
- Если `query` пустой, то делаем только сортировку, т.е. возвращаем все записи.
- В XML 2 поля с именем, наше поле Name это first_name + last_name из XML
- http://www.golangprograms.com/files-directories-examples.html - в помощь для работы с файлами
- Запуск тестов: `go test -cover`.
- Построение покрытия: `go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html`. 
