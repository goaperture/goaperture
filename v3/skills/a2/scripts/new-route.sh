#!/bin/bash
# Создаёт новый роут и автоматически генерирует файлы

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: new-route.sh <NAME> <DESCRIPTION>"
    echo "Example: new-route.sh GetUser 'Получить пользователя'"
    exit 1
fi

NAME="$1"
DESCRIPTION="$2"
SECURE="$3"

# Переходим в папку routes если есть
SCRIPT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
if [ -d "$SCRIPT_DIR/api/routes" ]; then
    cd "$SCRIPT_DIR/api/routes" || exit 1
else
    echo "Ошибка: не найдена папка api/routes"
    exit 1
fi

# Определяем версию (последняя папка v*)
LATEST_VERSION=$(ls -d v* 2>/dev/null | sort -V | tail -1)
if [ -z "$LATEST_VERSION" ]; then
    echo "Ошибка: не найдена версия (папка v*) в api/routes"
    exit 1
fi

cd "$LATEST_VERSION" || exit 1

# Создаём пакет если его нет
if [ ! -d "$NAME" ]; then
    mkdir -p "$NAME"
    echo "Created package: $NAME"
fi

cd "$NAME" || exit 1

# Создаём роут
if [ "$SECURE" = "--secure" ]; then
    a2 add-route "$NAME" "$DESCRIPTION" --sequre
else
    a2 add-route "$NAME" "$DESCRIPTION"
fi

# Генерируем файлы
cd "$SCRIPT_DIR"
a2 generate

echo "✓ Роут $NAME создан и сгенерирован"