///////////////////////////////////////////////////////////////////////
----------------first commit (router and form)----------------
///////////////////////////////////////////////////////////////////////

# First commit (router and form)

Запустил сервер в main.go с помощью созданного роутера
сервер слушает входящие запросы на порту 8080

роутер принимает запросы
роутер - диспечер: он получает от браузера URL - вызывает нужную функцию

роутер состоит из маршрутов:
-функция homeHandler 

    Обрабатывает запрос GET /

    Функция получает:
    - w http.ResponseWriter — объект для отправки ответа браузеру;
    - r *http.Request — контейнер с данными входящего запроса.

    homeHandler читает HTML-файл формы через template.ParseFiles.
    Если при чтении файла возникла ошибка, сервер возвращает ошибку 500.
    Если ошибки нет, функция tmpl.Execute отправляет HTML-страницу браузеру.

    В HTML-форме указан метод POST и адрес action="/form".
    Поэтому при нажатии кнопки отправки браузер делает POST-запрос на /form.

-функция formHandler 

    Обрабатывает отправку формы на адрес POST /form.

    проверяет что запрос пришел методом POST
    Если пользователь просто открыл /form в браузере через адресную строку, это будет GET-запрос, и сервер вернёт ошибку 405 Method Not Allowed.


    Если запрос POST, функция достаёт из контейнера запроса значения полей:
        - name
        - email

    Затем записывает их в отдельные переменные, формирует текстовый ответ и отправляет его пользователю.

-form.html 
    обычная форма
    
    особенностью является то что мы указываем метод POST (method="POST")
    ВАЖНО: указываем у полей ввода name (name="name", name="email"), т.к. именно по этим именнам мы обращаемся к ним с помощью метода FormValue
    именно по значениям атрибута name backend получает данные через:

    Браузер формирует POST-запрос примерно в таком виде:

    name=Валентин&email=test@mail.com

    После чего backend может получить эти значения из объекта Request.


ВАЖНО:
функции handlers написаны с большой буквы, например HomeHandler, FormHandler.
Это нужно потому что они находятся в другом package и router должен иметь к ним доступ.

///////////////////////////////////////////////////////////////////////
----------------second commit(validation, json and api)----------------
///////////////////////////////////////////////////////////////////////

***создал domain/application.go***

domain - это слой, где хранятся основные сущности проекта
Пока создал структуру Application:
    - Name
    - Email

Раньше данные были отдельными переменными:
    name
    email

Теперь данные собираются в один объект application

***создал validation/application_validator.go***

validation отвечает только за проверку данных.
ValidateApplication принимает application и возвращает error, eсли данные не соответствуют требованиям, иначе nil (ошибок нет)


в formHandler теперь логика такая:
1. проверяем что запрос POST
2. собираем данные из формы в domain.Application
3. передаем application в validation
4. если ошибка возвращаем 400 Bad Request
5. если все нормально возвращаем успешный ответ


***новый API endpoint:***
POST /api/applications

теперь бэкенд принимает json и возвращает json

В application_api_handler.go:
- проверяем метод POST
- читаем JSON из r.Body
- Decode заполняет структуру Application
- валидируем данные
- отправляем JSON ответ


***Важно!***
json.NewDecoder(r.Body).Decode(&application)

r.Body - тело HTTP-запроса.
В нем лежит JSON, который отправил frontend

Decode берет JSON и заполняет структуру application

&application передается потому что Decode должен изменить сам объект application!

***создал domain/api_response.go***

APIResponse - это структура ответа API

Поля:
- Status
- Message

У полей есть json tags:
json:"status"
json:"message"

Это нужно чтобы в JSON ответе поля были маленькими буквами:

{
  "status": "success",
  "message": "Заявка успешно принята"
}

***подключил script.js***

script.js перехватывает отправку формы

event.preventDefault() отменяет стандартную отправку формы браузером

fetch отправляет запрос вручную на /api/applications

fetch отправляет данные в формате JSON:
- method: POST
- Content-Type: application/json
- body: JSON.stringify(formData)

async и await нужны потому что HTTP-запрос выполняется не мгновенно
await ждет ответ от сервера, но не ломает работу браузера

response.json() читает JSON ответ от backend

/////////////////////////////////////////////////////////
-----------------------third commit---------------------
////////////////////////////////////////////////////////

***создал application_service***

service — это слой бизнес логики

он НЕ хранит данные, а решает что с ними делать

он содержит:
    - структуру ApplicationService — объект, который работает с заявками
    - внутри него есть repo (repository), через который происходит сохранение
    - конструктор NewApplicationService для создания service
    - метод Create — принимает заявку и передает ее в repository

ВАЖНО:
service не знает как именно хранятся данные (в памяти или в БД)
он просто вызывает repo.Save()

---

***создал application_repository***

repository — слой хранения данных

это интерфейс, который описывает:
    - метод Save для сохранения заявки

ВАЖНО:
интерфейс не содержит логики, только описание

---

***создал memory_application_repository***

это временная реализация repository

он содержит:
    - структуру MemoryApplicationRepository — хранит массив заявок
    - конструктор NewMemoryApplicationRepository
    - метод Save — добавляет заявку в массив через append

ВАЖНО:
данные хранятся только в памяти
при перезапуске сервера они исчезают

---

***изменил application_service***

теперь service использует repository:

ApplicationService содержит внутри repo

метод Create:
    принимает application
    вызывает repo.Save(application)

---

***изменил application_api_handler***

теперь handler НЕ работает напрямую с данными

в handler:
    создается repo (memory repository)
    создается service
    handler передает application в service

pipeline теперь такой:

handler
↓
validation
↓
service
↓
repository
↓
данные сохраняются