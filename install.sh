#!/bin/sh

REPO=${REPO:-SimonHylander/diskotective}
REMOTE=${REMOTE:-https://github.com/${REPO}}

VERSION="v0.0.1"
BINARY="diskotective"
GOOS="linux"
GOARCH="amd64"

FILENAME="${BINARY}-${VERSION}-${GOOS}-${GOARCH}.tar.gz"

URL="${REMOTE}/releases/download/${VERSION}/${FILENAME}"
TARGET=/usr/local/bin/

curl -L "$URL" | tar zxvf - -C "$TARGET" --overwrite

if [ ! -f "$TARGET/$BINARY" ]; then
    echo "Failed to extract '$FILENAME' to '$TARGET'"
    exit 1
fi

if [ ! -w "$TARGET" ]; then
  echo "Failed to determine a suitable installation path. Please ensure $TARGET is writable and try again."
  exit 1
fi

chmod +x "${TARGET}${BINARY}"

echo "$FILENAME -> ${TARGET}${BINARY}"