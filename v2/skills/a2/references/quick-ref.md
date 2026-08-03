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

## Установка скилла

```bash
# Автоматически через a2
a2 install-skill

# Проверка
ls ~/.agents/skills/a2/
```