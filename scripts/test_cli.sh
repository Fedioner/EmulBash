#!/bin/sh
cd "$(dirname "$0")/.." || exit 1
make build vfs || exit 1

echo "== без параметров"
./build/emulator

echo "== только -vfs"
./build/emulator -vfs build/vfs/multi.zip

echo "== только -script"
./build/emulator -script examples/start.txt

echo "== -vfs и -script"
./build/emulator -vfs build/vfs/multi.zip -script examples/start.txt

echo "== все параметры сразу"
./build/emulator -vfs build/vfs/minimal.zip -script examples/start.txt -config examples/config.ini
