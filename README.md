
# Тестовое задание:
Написать сервис для поздравлений с днем рождения.

- Цель удобное поздравление сотрудников
- Получения списка сотрудников любым способом(api/ad ldap/прямая регистрация)
- Авторизация
- Возможность подписаться на отписаться от оповещения о дне рождения
- Оповещение о ДР того на кого подписан
- Внешнее взаимодействие (json арi или фронт или тг бот)
- В случае взаимодействия через тг бот (создание группы и добавление в нее всех подписанных)
- В случае взаимодействие через фронт настройка времени оповещения до дня рождения на почту. Предполагаемое время до 3х часов.

Подсказки*
Нужно реализовать (минимум):
- модуль авторизации
- модуль подключения к базе данных(можно просто хранить файл для каждого пользователя)
- модуль подписки
- модуль оповещения(cron)

--- 

# Пояснения к решению

Авторизация происходит при помощи JWT-токенов, которые имеют время жизни 12ч. Для подписи на уведомление о др необходимо указать id пользователя на которого хотите подписаться.

Уведомление о др приходит на почту, указанную при регистрации, за 3 часа до наступления праздника (можно изменить).

Для взаимодействия через фронт необходимо запустить сервис и открыть с помощью браузера файл ```static/index.html```

---

# Установка и запуск

```
git clone https://github.com/Rpqshka/hb-notification.git
```

```
cd hb-notification
```

```
docker-compose up --build
```

Вы можете использовать тестовую конфигурацию, которая находится в файле ```.env```, либо настроить сервис под себя и добавить этот файл в ```.gitignore```

---

# Работа с сервисом

После загрузки БД и запуска сервиса становятся доступны эндпоинты:
- POST   /auth/sign-up      (Регистрация пользователя)
- POST   /auth/sign-in      (Логин)
- GET    /users             (Получение списка всех пользователей)
- GET    /subscriptions     (Получение списка подписок)
- POST   /subscribe         (Подписаться на уведовления)
- DELETE /unsubscribe       (Отписаться от уведомлений)
- POST   /update-cron       (Обновить время до уведомления)

[![postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/24093475-c6ec3bd2-ea37-4c4b-ab16-a1e38d1ce246?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D24093475-c6ec3bd2-ea37-4c4b-ab16-a1e38d1ce246%26entityType%3Dcollection%26workspaceId%3Dfb0d85ac-394b-408f-8648-0747f742ba1c)


---

# Веб интерфейс

[![2024-06-08-154828320.png](https://i.postimg.cc/WbyYRD6T/2024-06-08-154828320.png)](https://postimg.cc/2V47hSKK)

---

# Пример уведовления о ДР

[![notification-example.png](https://i.postimg.cc/9FVH62rM/notification-example.png)](https://postimg.cc/bd65nK0c)

---

# Контакты
Telegram : @rpqshka

email: rpqshka@gmail.com