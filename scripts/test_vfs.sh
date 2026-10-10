#!/bin/sh
cd "$(dirname "$0")/.." || exit 1
make build vfs || exit 1

echo "== минимальная VFS (один файл)"
./build/emulator -vfs build/vfs/minimal.zip -script examples/vfs_info.txt

echo "== VFS с несколькими файлами"
./build/emulator -vfs build/vfs/multi.zip -script examples/vfs_info.txt

echo "== VFS с вложенностью больше 3 уровней"
./build/emulator -vfs build/vfs/deep.zip -script examples/vfs_info.txt

echo "== VFS не задана"
./build/emulator -script examples/vfs_info.txt

echo "== VFS не существует"
./build/emulator -vfs build/vfs/nope.zip -script examples/vfs_info.txt

echo "== файл не является zip архивом"
./build/emulator -vfs build/vfs/broken.zip -script examples/vfs_info.txt
