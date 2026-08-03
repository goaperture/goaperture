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

### Защищённые роуты (приватные)

Добавь флаг `--sequre` для роутов, требующих авторизации:

```bash
a2 add-route AdminRoute "Админская функция" --sequre
```

Права доступа проверяются через `auth.Permissions` в `Payload`.

#### Приватный доступ вручную (`PrivateAccess: true`)

Чтобы сделать роут приватным (требует авторизации), задай поле `PrivateAccess: true` в описании роута:

```go
var AdminPanel = aperture.Route[AdminPanelInput, AdminPanelOutput]{
	Description:   "Админ панель",
	PrivateAccess: true, // <- роут требует авторизацию
	Handler:       AdminPanelHandler,
}
```

Поле `PrivateAccess bool` — часть структуры `aperture.Route` (и `aperture.Switch` для WebSocket). При `PrivateAccess: true` фреймворк:

- проверяет авторизацию перед вызовом хендлера (см. `api/auth.go`);
- блокирует доступ без валидного токена/прав;
- права проверяются через `auth.Permissions` в `Payload`.

То же поле используется для приватных WebSocket-эндпоинтов в `api/ws/...`.

## Пагинация с ent ORM

GoAperture имеет встроенную поддержку пагинации для ent ORM. Ответ автоматически получает мета-данные `pagination` (`page`, `size`, `total`).

### Установка дополнения `ent-templates-export`

1. Перейди в папку проекта (где лежит `generate.go` для ent) и выполни:

   ```bash
   a2 ent-templates-export
   ```

   Команда создаст файл `pagination.tmpl` с шаблонами, которые добавляют метод `Paginate` к ent-запросам, а также импорты `api/client` и `responce`.

2. Добавь шаблон в `generate.go`:

   ```bash
   # в директиве //go:generate entc generate ...
   --template ./<имя_папки>/

   # например
   --template ./ent
   ```

   При необходимости добавь фичи sql:

   ```bash
   --template ./<имя_папки>/ --feature sql/execquery,sql/upsert
   ```

3. Перегенерируй ent-код.

### Использование `Paginate` в хендлере

Метод `Paginate(ctx, page, size)` добавляет к запросу `LIMIT`/`OFFSET`, считает общее количество записей и кладёт результат в контекст через `client.SetPagination`. Пагинация попадёт в ответ автоматически — фреймворк вызывает `client.WithPagination(ctx)` для каждого роута перед вызовом хендлера.

Минималистичный хендлер для получения списка заказов с пагинацией:

```go
func GetOrdersHandler(ctx context.Context, input GetOrdersInput) GetOrdersOutput {
	orders := db.Client.Order.Query().Paginate(ctx, input.Page, 10).AllX(ctx)
	return orders
}
```

Где:

- `input.Page` — номер страницы из входных данных роута;
- `10` — размер страницы;
- `db.Client` — клиент ent-базы данных из проекта.

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
5. Для приватного доступа укажи `PrivateAccess: true` в описании роута
6. Для пагинации с ent ORM используй `a2 ent-templates-export` и метод `Paginate`

## Установка скилла на другой компьютер

### Автоматическая установка скила

```bash
# Установка скилла в глобальную директорию
a2 install-skill

# Создаст ~/.agents/skills/a2 и скопирует туда файлы
```
