# Сервис подбора фильмов MOVIEWORLD


api

actors:
1. GET /actors/ - список всех актеров
2. GET /actor/{ACTOR_ID} информация о конкретном актере

films:
1. GET /films - список всех фильмов, может принимать query параметры
2. GET /films/by/{ACTOR_ID} список фильмов в которых снимался актер с таким айди
3. GET /film/{FILM_ID} информация о конкретном фильме
4. GET /films/soon/ список предстоящих релизов
5. GET /films/favourite - избранные фильмы пользователя
6. POST /films/favourite/{FILM_ID} - добавить фильм в избранное
7. DELETE /films/favourite/{FILM_ID} - удаление фильма из избранного
8. GET /film/{FILM_ID}/actors список актеров сыгравших в фильме
9. GET /film/{FILM_ID}/genres список жанров фильма

auth:
1. POST /register - регистрация
2. POST /login - вход по логину и паролю

review:
1. POST /review/{FILM_ID} - оставить отзыв
2. DELETE /review/{REVIEW_ID} - удалить отзыв
3. PUT /review/{REVIEW_ID} - изменить отзыв
4. GET /review/{FILM_ID} - получение отзывов о фильме

search of actors and films
1. GET /search/{DATA} - регистронезависимый поиск актеров и фильмов, где есть вхождение строки DATA в названии фильма или его режиссера или в имени + фамилии актера
