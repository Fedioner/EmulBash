#!/bin/sh
cd "$(dirname "$0")/.." || exit 1
make build || exit 1

echo "== только -config (vfs и script берутся из файла)"
./build/emulator -config examples/config.ini

echo "== -config и -vfs (vfs из командной строки важнее)"
./build/emulator -config examples/config.ini -vfs other.zip

echo "== -config и -script (script из командной строки важнее)"
./build/emulator -config examples/config.ini -script examples/start_error.txt

echo "== скрипт с exit, окно сразу закроется"
./build/emulator -config examples/config.ini -vfs other.zip -script examples/start_exit.txt
echo "код выхода: $?"
