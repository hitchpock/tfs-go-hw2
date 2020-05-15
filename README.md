# Домашнее задание по Golang.

## Структура репозитория

```
.
├── main.go              # файл с исходными кодами проекта

```

## Задание
В инвестициях есть социальная сеть - Пульс. Для развития платформы необходимо анализировать пользователей и их взаимоотношения.


### Цель
Найти через сколько рукопожатий знакомы пользователи


### Вход
В файле users.json содержится массив структур. Формат структуры: Nick (ник), Email (email), Created_at (Время создания), Subscribers (подписчики). В файле input.csv хранятся email-ы пользователей, "расстояние" между которыми требуется найти


### Выход
Файл result.json - это массив структур. Формат структуры:

{

 "id": 0, - номер строки из файла input.csv

 "from": "email1",

 "to": "email2",

"path": [{"email": "example@email.com", "created_at": "2017-11-30"}]

}


### Примечания

* Необходимо найти НАИМЕНЬШЕЕ количество рукопожатий
* Путь должен быть отсортирован в правильной последовательности, то есть from->path[0]->...path[len(path)-1]->to
* Если пользователи никак не могуть быть знакомы или знакомы лично , то массива path быть не должно