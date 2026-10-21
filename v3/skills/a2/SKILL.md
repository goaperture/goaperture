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

## Хук `Prepare` и коллектор `CL`

Помимо `Handler` и `PrivateAccess`, структура `aperture.Route[I, O]` имеет поле `Prepare` — это опциональная функция-хук, которая **не выполняется во время HTTP-запроса**. Она нужна для сбора примеров входных/выходных данных при генерации документации или схемы API.

### Сигнатура

```go
type Prepare[P Input, O Output] = func(collector *CL[P, O])
type CL[P Input, O Output] = collector.Collector[P, O]
```

### Где вызывается

В `aperture.Handle()` при оборачивании роута создаётся `Switch` с полем `PrepareCall: func() collector.RouteDump`. При вызове `PrepareCall`:

1. создаётся пустой `Collector` с `Handler` роута;
2. вызывается `route.Prepare(&cll)`;
3. в конце возвращается `RouteDump` со всеми собранными данными (метод, описание, inputs, outputs, errors, флаг пагинации).

Видно в `api/aperture/doc.go`: `dump := route.PrepareCall()`.

### Методы коллектора `CL[I, O]`

Внутри `Prepare` доступен коллектор с методами:

| Метод | Что делает |
|-------|-----------|
| `Execute(input I) *CL` | Реальный прогон хендлера: создаёт `ctx` через `client.WithPagination`, кладёт `input` в `Inputs`, вызывает `c.Handler(ctx, input)`, сохраняет результат в `Outputs`, обновляет `WithPagination`. Паники через `recover()` ловятся и пишутся в `Errors` — наружу не пробрасываются. |
| `Entry(input I) *CL` | Только регистрирует `input` в `Inputs`, хендлер **не вызывается**. Используется в паре с `Expect`/`ExpectArray` для кейсов, которые нельзя воспроизвести (side-effects в БД, внешние вызовы и т.п.). |
| `Expect(output any) *CL` | Регистрирует ожидаемый выход в `Outputs` как один элемент. Парный к `Entry`. |
| `ExpectArray(output any) *CL` | То же, что `Expect`, но оборачивает результат в `[]any{output}`. Удобно для кейсов с массивами. |
| `GetDump() RouteDump` | Собирает финальный дамп для документации. |

Структура `RouteDump`:

```go
type RouteDump struct {
    Method         string
    Description    string
    AccessKey      string
    Inputs         []any
    Outputs        []any
    Errors         []string
    WithPagination bool
}
```

Поля `Method`, `Description`, `AccessKey` в `GetDump()` берутся из коллектора напрямую (без методов-сеттеров — заполняются автоматически или остаются пустыми).

### Пример из шаблона

```go
type HelloWorldInput struct{}
type HelloWorldOutput interface{ any }

var HelloWorld = aperture.Route[HelloWorldInput, HelloWorldOutput]{
    Description: "HW",
    Handler:     HelloWorldHandler,
    Prepare: func(cl *aperture.CL[HelloWorldInput, HelloWorldOutput]) {
        cl.Execute(HelloWorldInput{})
    },
}
```

### Типичные паттерны

```go
Prepare: func(cl *aperture.CL[MyInput, MyOutput]) {
    // реальный прогон хендлера (может дёрнуть БД)
    cl.Execute(MyInput{Foo: "bar"})

    // ручное описание кейса без вызова хендлера
    cl.Entry(MyInput{Foo: "baz"}).Expect(MyOutput{Result: "qux"})

    // для массивов
    cl.Entry(MyInput{}).ExpectArray([]MyItem{...})
}
```

### Нюансы

- `Execute` **выполняет реальный код хендлера** — учитывайте side-effects (БД, письма, логи). Паники не пробьются наружу, они попадут в `Errors` дампа.
- `Handler` коллектора заполняется автоматически в `handle.go` (внутри `Handle()`), в `Prepare` менять его не нужно.
- Коллектор не thread-safe — `Prepare` предполагается синхронным.
- Если `Prepare` не указан, `PrepareCall` вернёт пустой дамп — документация/схема для роута будет без примеров.
- Для WebSocket-топиков есть аналогичный хук `Prepare` (`ws/aperture/topic.go`) с похожим `PrepareCall`, но работающий с `Collector[Message]`.

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
7. Используй `Prepare` для сбора примеров входных/выходных данных — это улучшает автогенерируемую документацию роута

## Установка скилла на другой компьютер

### Автоматическая установка скила

```bash
# Установка скилла в глобальную директорию
a2 install-skill

# Создаст ~/.agents/skills/a2 и скопирует туда файлы
```
