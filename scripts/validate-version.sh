#!/bin/bash

tar zxf out/package/nagiosfoundation-linux-amd64-*.tgz
CHECK_VERSION=$(./godelw project-version)
EXIT=0

for CHECK_NAME in $(ls bin); do
    OUTPUT=$(bin/$CHECK_NAME version)

    echo $OUTPUT

    OUT_NAME=$(echo $OUTPUT | cut -d \  -f 1)
    OUT_VERSION=$(echo $OUTPUT | cut -d \  -f 3)

    if [ "$OUT_NAME" != "$CHECK_NAME" ]; then
        echo "Check name is $CHECK_NAME but output name is $OUT_NAME"
        EXIT=1
    elif [ "$OUT_VERSION" != "$CHECK_VERSION" ]; then
        echo "Version is $CHECK_VERSION but output version is $OUT_VERSION"
        EXIT=1
    fi
done

exit $EXIT
