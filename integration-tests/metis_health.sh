#!/bin/bash
set -e

while true
do
    peers=$(docker exec metis0 bash -c "metis attach /root/.metis/data/metis.ipc -exec 'admin.peers'")
    block=$(docker exec metis0 bash -c "metis attach /root/.metis/data/metis.ipc -exec 'eth.blockNumber'")

    if [[ -n "$peers" ]] && [[ -n "$block" ]]; then
        break
    fi
done

echo $peers
echo $block
