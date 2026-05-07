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

# Second commit(validation, json and api)

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


# Third commit(temporary database(sql), service for work with data)


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




# Fourth commit(migrations, config server, sql injection, handler get and general handler)


feat(database): add MySQL storage and application listing

Add database migration for applications table and extend Application model
with ID field.

Add config layer for loading server and database settings from environment.

Create MySQL connection setup and MySQLApplicationRepository for saving
applications into the database.

Replace memory repository flow with MySQL repository and inject repository
from main.go into handlers.

Add GetAll repository and service methods, GET applications handler and
general ApplicationsHandler to support both GET and POST on /api/applications.



***подготовил структуру базы данных***

создал migrations/001_create_applications.sql

в нем описана таблица applications

таблица содержит:
    -id - уникальный номер заявки
    -name - имя пользователя
    -email - email пользователя
    -created_at - дата создания заявки

***изменил domain.Application***

добавил поле ID

ID нужен потому что после сохранения заявки в базе данных
у каждой заявки должен быть свой уникальный номер


***создал config layer***

создан .env файл для хранения настроек проекта:
    - SERVER_PORT
    - DB_USER
    - DB_PASSWORD
    - DB_HOST
    - DB_PORT
    - DB_NAME

ВАЖНО:
настройки вынесены отдельно от кода

---

***создал config.go***

config отвечает за загрузку настроек проекта

LoadConfig:
    - загружает .env через godotenv.Load()
    - собирает настройки в структуру Config
    - возвращает готовый объект Config

создана функция getEnv:
    - получает переменные окружения через os.Getenv()
    - если переменной нет, возвращает значение по умолчанию

---

***изменил main.go***

main.go теперь:
    - загружает config
    - получает порт сервера из .env
    - запускает сервер через cfg.ServerPort

---

***создал app.NewDB***

NewDB отвечает за подключение к MySQL

внутри:
    - собирается DSN строка подключения
    - sql.Open создает connection pool
    - db.Ping проверяет доступность БД

ВАЖНО:
database/sql работает через mysql driver

используется side-effect import:

_ "github.com/go-sql-driver/mysql"

driver регистрируется внутри database/sql
и после этого становится доступен sql.Open("mysql", ...)

---

***создал mysql_application_repository***

repository теперь работает с реальной MySQL базой данных

структура MySQLApplicationRepository содержит:
    - db *sql.DB

метод Save:
    - выполняет INSERT INTO applications
    - сохраняет name и email
    - получает LastInsertId()
    - записывает ID обратно в application

ВАЖНО:
repository теперь сохраняет данные не в массив, а в MySQL

---

***изменил architecture flow***

раньше:
handler -> service -> memory repository

теперь:
handler -> service -> mysql repository -> MySQL

---

***добавил dependency injection***

repository создается в main.go
и передается в handlers через SetApplicationRepository

handler больше не создает repository самостоятельно

main.go теперь собирает приложение:
    - config
    - db
    - repository
    - handlers


***добавил получение заявок из базы данных***

добавил метод GetAll в repository interface

теперь repository умеет:
    - Save — сохранять заявку
    - GetAll — получать список всех заявок

***изменил mysql_application_repository***

добавил метод GetAll

он выполняет SQL запрос:

    SELECT id, name, email
    FROM applications
    ORDER BY id DESC

Query используется для SELECT запросов

rows.Next() проходит по строкам результата

rows.Scan() переносит данные из строки БД в структуру Application

append добавляет каждую заявку в общий массив applications

ВАЖНО:
Exec используется для INSERT / UPDATE / DELETE
Query используется для SELECT

---

***изменил service***

добавил метод GetAll

service вызывает repository.GetAll()

---

***добавил GetApplicationsHandler***

handler обрабатывает GET запрос и возвращает список заявок в JSON

ответ имеет вид:

{
  "status": "success",
  "data": [...]
}

---

***добавил общий ApplicationsHandler***

теперь один endpoint /api/applications работает по разным HTTP методам:

GET /api/applications
    получить список заявок

POST /api/applications
    создать новую заявку

ВАЖНО:
один URL может выполнять разные действия в зависимости от HTTP метода



# Fifth commit(put and delete, full crud api)

***добавил обновление заявки***

добавил метод Update в repository interface

теперь repository умеет:
    - Save — создать заявку
    - GetAll — получить все заявки
    - Update — изменить заявку по ID

---

***изменил mysql_application_repository***

добавил метод Update

он выполняет SQL запрос:

    UPDATE applications
    SET name = ?, email = ?
    WHERE id = ?

ВАЖНО:
PUT используется для изменения уже существующей записи

адрес теперь может быть динамическим:

    PUT /api/applications/1

это значит:
    изменить заявку с ID = 1

---

***добавил UpdateApplicationHandler***

handler:
    - достает ID из URL
    - проверяет что ID корректный
    - читает JSON из r.Body
    - валидирует данные
    - вызывает service.Update()
    - возвращает JSON ответ

---

***добавил удаление заявки***

добавил метод Delete в repository interface

добавил метод Delete в service

добавил метод Delete в mysql_application_repository

SQL запрос:

    DELETE FROM applications
    WHERE id = ?

---

***добавил DeleteApplicationHandler***

handler:
    - достает ID из URL
    - проверяет ID
    - вызывает service.Delete()
    - возвращает JSON ответ

---

***изменил ApplicationsHandler***

теперь один endpoint поддерживает четыре метода:

GET /api/applications
    получить список заявок

POST /api/applications
    создать заявку

PUT /api/applications/{id}
    изменить заявку по ID

DELETE /api/applications/{id}
    удалить заявку по ID

---

***итог***

теперь API умеет полный CRUD:

Create — POST
Read — GET
Update — PUT
Delete — DELETE

!!ВАЖНО!!
Exec используется для INSERT / UPDATE / DELETE
Query используется для SELECT


# Sixth commit(auth base, users and generated credentials)

***начал делать авторизацию пользователя***

создал таблицу users

таблица содержит:
    - id — уникальный номер пользователя
    - application_id — связь пользователя с заявкой
    - login — логин пользователя
    - password_hash — хеш пароля
    - created_at — дата создания пользователя

ВАЖНО:
пароль в базу данных не сохраняется в открытом виде
в базу сохраняется только password_hash

---

***создал domain.User***

User — это сущность пользователя

она содержит:
    - ID
    - ApplicationID
    - Login
    - PasswordHash

---

***создал security/password.go***

security отвечает за работу с логином и паролем

добавил функции:
    - GenerateLogin()
    - GeneratePassword()
    - HashPassword()
    - CheckPassword()

GenerateLogin создает логин автоматически

GeneratePassword создает случайный пароль

HashPassword превращает пароль в хеш

CheckPassword сравнивает введенный пароль с хешем из базы данных

---

***изменил создание заявки***

теперь при создании заявки backend делает не только INSERT в applications

теперь flow такой:

POST /api/applications
↓
handler
↓
validation
↓
service.Create()
↓
repository.Save()
↓
создание login и password
↓
HashPassword(password)
↓
repository.CreateUser()
↓
MySQL

---

***добавил CreateUser в repository***

repository теперь умеет создавать пользователя для заявки

SQL запрос:

    INSERT INTO users (application_id, login, password_hash)
    VALUES (?, ?, ?)

---

***изменил APIResponse***

добавил поле Data

теперь API может возвращать не только status и message,
но и дополнительные данные

например после создания заявки backend возвращает:

{
  "status": "success",
  "message": "Заявка успешно принята",
  "data": {
    "login": "...",
    "password": "..."
  }
}

ВАЖНО:
обычный пароль показывается пользователю один раз
в базе хранится только хеш пароля