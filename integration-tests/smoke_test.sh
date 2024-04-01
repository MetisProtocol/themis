#!/bin/bash
set -e

balanceInit=$(docker exec metis0 bash -c "metis attach /root/.metis/data/metis.ipc -exec 'Math.round(web3.fromWei(eth.getBalance(eth.accounts[0])))'")

stateSyncFound="false"
SECONDS=0
start_time=$SECONDS

while true
do
  
    balance=$(docker exec metis0 bash -c "metis attach /root/.metis/data/metis.ipc -exec 'Math.round(web3.fromWei(eth.getBalance(eth.accounts[0])))'")

    if ! [[ "$balance" =~ ^[0-9]+$ ]]; then
        echo "Something is wrong! Can't find the balance of first account in metis network."
        exit 1
    fi

    if (( $balance > $balanceInit )); then
        if [ $stateSyncFound != "true" ]; then 
            stateSyncTime=$(( SECONDS - start_time ))
            stateSyncFound="true" 
        fi      
    fi


    if [ $stateSyncFound == "true" ] ; then
        break
    fi    

done
echo "Both state sync  went through. All tests have passed!"
echo "Time taken for state sync: $(printf '%02dm:%02ds\n'  $(($stateSyncTime%3600/60)) $(($stateSyncTime%60)))"
