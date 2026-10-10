#!/bin/sh
cd "$(dirname "$0")/.." || exit 1
mkdir -p build/vfs
for name in minimal multi deep; do
	rm -f "build/vfs/$name.zip"
	(cd "vfs/$name" && zip -r -q -X "../../build/vfs/$name.zip" . -x '*.DS_Store') || exit 1
	echo "build/vfs/$name.zip"
done
echo "это не zip архив" > build/vfs/broken.zip
echo "build/vfs/broken.zip"
