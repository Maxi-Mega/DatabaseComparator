#!/bin/bash -eu

VERSION="0.0.1"
OUTPUT_FILENAME="DatabaseComparator-$VERSION"
OUTPUT_DIR="compiled"

echo "Building $OUTPUT_FILENAME ..."
go build -o $OUTPUT_DIR/$OUTPUT_FILENAME . && strip $OUTPUT_DIR/$OUTPUT_FILENAME && xz $OUTPUT_DIR/$OUTPUT_FILENAME
echo "Done !"