#!/usr/bin/env sh

if [ "$1" = 'themiscli' ]; then
    shift
    exec themiscli --home=$THEMIS_DIR "$@"
fi

exec themisd --home=$THEMIS_DIR "$@"
