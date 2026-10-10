#!/bin/sh
cd "$(dirname "$0")/.." || exit 1
make build vfs || exit 1

echo "== неизвестный параметр"
./build/emulator -abc
echo "код выхода: $?"

echo "== параметр без значения"
./build/emulator -vfs
echo "код выхода: $?"

echo "== лишний аргумент"
./build/emulator -script examples/start.txt lishnee
echo "код выхода: $?"

echo "== конфиг не существует"
./build/emulator -config examples/nope.ini
echo "код выхода: $?"

echo "== неизвестный ключ в конфиге"
./build/emulator -config examples/bad_key.ini
echo "код выхода: $?"

echo "== строка без = в конфиге"
./build/emulator -config examples/broken.ini
echo "код выхода: $?"

echo "== скрипт не существует"
./build/emulator -vfs build/vfs/multi.zip -script examples/nope.txt

echo "== ошибка в скрипте, выполнение останавливается"
./build/emulator -script examples/start_error.txt
