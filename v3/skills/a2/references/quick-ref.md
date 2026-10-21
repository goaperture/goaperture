# Быстрая справка по a2

## Команды

| Команда | Описание |
|---------|----------|
| `a2 init` | Инициализировать проект |
| `a2 init ./dir` | Инициализировать в директории |
| `a2 generate` | Перегенерировать роуты |
| `a2 generate -a app` | Генерация с указанием модуля |
| `a2 generate -p api` | Генерация с указанием папки |
| `a2 add-route Name "Desc"` | Добавить роут |
| `a2 add-route Name "Desc" --sequre` | Добавить защищённый роут |
| `a2 install-skill` | Установить скилл в `~/.agents/skills/a2` |
| `a2 ent-templates-export` | Создать ent-шаблон `pagination.tmpl` в текущей папке |

## Флаги

| Флаг | Описание |
|------|----------|
| `-n, --name` | Название роута / скилла |
| `-d, --description` | Описание роута |
| `-s, --sequre` | Защищённый роут |

## Минимальный workflow

```bash
# Установка
go install github.com/goaperture/goaperture/v2/cli/a2@latest

# Инициализация
a2 init

# Добавление роута
cd api/routes/v1/Hello
a2 add-route MyRoute "Описание"
cd ../../..

# Генерация
a2 generate

# Установка скилла на другой компьютер
a2 install-skill
```

## Приватные роуты

Флаг `--sequre` при создании роута добавляет `PrivateAccess: true`:

```go
var AdminPanel = aperture.Route[AdminPanelInput, AdminPanelOutput]{
	Description:   "Админ панель",
	PrivateAccess: true, // требует авторизацию
	Handler:       AdminPanelHandler,
}
```

Права проверяются через `auth.Permissions` в `Payload`.

## Пагинация с ent ORM

```bash
# 1. Создать шаблон пагинации в папке с generate.go
cd ent
a2 ent-templates-export

# 2. Добавить шаблон в generate.go
#    --template ./ent

# 3. Перегенерировать ent
```

Хендлер со списком заказов и пагинацией:

```go
func GetOrdersHandler(ctx context.Context, input GetOrdersInput) GetOrdersOutput {
	orders := db.Client.Order.Query().Paginate(ctx, input.Page, 10).AllX(ctx)
	return orders
}
```

`client.WithPagination(ctx)` вызывается фреймворком автоматически для каждого роута — ответ получит `pagination: {page, size, total}`.

## Установка скилла

```bash
# Автоматически через a2
a2 install-skill

# Проверка
ls ~/.agents/skills/a2/
```