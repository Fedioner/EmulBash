#!/bin/sh
cd "$(dirname "$0")/.." || exit 1
make build || exit 1

echo "== без параметров"
./build/emulator

echo "== только -vfs"
./build/emulator -vfs vfs/test.zip

echo "== только -script"
./build/emulator -script examples/start.txt

echo "== -vfs и -script"
./build/emulator -vfs vfs/test.zip -script examples/start.txt

echo "== все параметры сразу"
./build/emulator -vfs other.zip -script examples/start.txt -config examples/config.ini
