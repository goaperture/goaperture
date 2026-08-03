---
name: a2
description: CLI утилита для работы с фреймворком GoAperture. Инициализирует проекты, генерирует роуты и добавляет эндпоинты. Используй для создания и настройки file-based API проектов.
license: MIT
compatibility: Requires Go 1.21+, a2 CLI installed globally
metadata:
  embed: skills/a2
  version: v2.0
  author: GoAperture team
---

# a2 - GoAperture CLI

## Описание

a2 (Aperture API) — это CLI утилита для фреймворка GoAperture v2, которая позволяет создавать API на основе файловой структуры. Она упрощает работу с роутами, генерацией кода и управлением эндпоинтами.

## Установка

Перед первым использованием установи утилиту:

```bash
go install github.com/goaperture/goaperture/v2/cli/a2@latest
```

Убедись, что `$GOPATH/bin` в PATH, чтобы команда `a2` была доступна.

## Команды

### init

Инициализирует базовую структуру проекта:

```bash
a2 init           # Инициализация в текущей директории
a2 init ./myproj  # Инициализация в указанной директории
```

Создаёт:
- Шаблоны сервера (`api/server.go`)
- Шаблоны аутентификации (`api/auth.go`)
- Примеры роутов (`api/routes/v1/Hello/`)

После инициализации выполните:
```bash
cd .
go mod init app
go mod tidy
a2 generate
```

### generate

Перегенерирует роуты и WebSocket эндпоинты:

```bash
a2 generate              # С генерацией из текущей директории
a2 generate -a myapp     # С указанием модуля go
a2 generate -p api       # С указанием папки с роутами
```

Используйте эту команду после добавления или изменения роутов.

### add-route

Создаёт новый эндпоинт:

```bash
a2 add-route <NAME> <DESCRIPTION>
# или с флагами
a2 add-route -n <NAME> -d <DESCRIPTION> [-s]
```

Параметры:
- `-n, --name` — Название роута (обязательный)
- `-d, --description` — Описание эндпоинта (обязательный)
- `-s, --sequre` — Создать защищённый роут (требует авторизации)

Пример:
```bash
cd api/routes/v1/Hello
a2 add-route GetUser "Получить пользователя по ID"
a2 add-route -n AdminPanel -d "Админ панель" --sequre
```

### install-skill

Устанавливает скилл a2 в глобальную директорию:

```bash
a2 install-skill

# Создаст ~/.agents/skills/a2 и скопирует туда файлы скилла
# После установки используйте: /skill:a2
```

Опции:
- `-n, --name` — Название скилла для установки (по умолчанию: a2)

## Структура проекта

```
project/
├── api/
│   ├── server.go          # Конфигурация сервера (порт, auth)
│   ├── auth.go            # Аутентификация и permissions
│   ├── routes/
│   │   └── v1/
│   │       └── Hello/
│   │           ├── HelloWorld.go
│   │           └── PrivateWorld.go
│   └── ws/
│       └── v1/
│           └── ...
├── main.go
└── go.mod
```

## Работа с роутами

### Создание роута

1. Перейди в папку с версиями (например `api/routes/v1/`):
   ```bash
   cd api/routes/v1/Hello
   ```

2. Добавь новый роут:
   ```bash
   a2 add-route MyRoute "Описание роута"
   ```

3. Сгенерируй обновлённые файлы:
   ```bash
   a2 generate
   ```

### Защищённые роуты

Добавь флаг `--sequre` для роутов, требующих авторизации:

```bash
a2 add-route AdminRoute "Админская функция" --sequre
```

Права доступа проверяются через `auth.Permissions` в `Payload`.

## Полезные комбинации

```bash
# Полная инициализация нового проекта
go install github.com/goaperture/goaperture/v2/cli/a2@latest
a2 init
go mod init app && go mod tidy
a2 generate

# Добавление и генерация
cd api/routes/v1/Hello
a2 add-route NewEndpoint "Новое эндпоинт"
cd ../..
a2 generate

# Установка скилла на другой компьютер
a2 install-skill
```

## Советы

1. Всегда запускай `a2 generate` после изменения роутов
2. Названия роутов должны быть в CamelCase
3. Описания на русском языке
4. Защищённые роуты требуют наличия пользователя с правами доступа

## Установка скилла на другой компьютер

### Встроенный способ (через pi)

Скилл доступен как embed-ресурс в проекте. Используйте в другом проекте через settings.json:

```json
{
  "skills": ["./skills/a2"]
}
```

### Автоматическая установка

```bash
# Установка скилла в глобальную директорию
a2 install-skill

# Создаст ~/.agents/skills/a2 и скопирует туда файлы
```

### Ручная установка

```bash
# Скопировать скилл в глобальную директорию
mkdir -p ~/.agents/skills/a2
cp -r /path/to/project/skills/a2/* ~/.agents/skills/a2/
```

## Embed информация

Этот скилл может быть встроен (embed) в другие проекты:
- Каталог: `skills/a2`
- Основной файл: `SKILL.md`
- Скрипты: `scripts/new-route.sh`
- Справка: `references/quick-ref.md`

Для встраивания в Go проект используйте:

```go
//go:embed skills/a2/*
var skillFS embed.FS
```

## Проверка установки

```bash
ls ~/.agents/skills/a2/
# SKILL.md
# scripts/new-route.sh
# references/quick-ref.md
```