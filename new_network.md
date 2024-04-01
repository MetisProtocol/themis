# Init a new network


## 1. Deploy L1 Contract
```
cd {path}/contracts
npm run deploy {network}
npm run verify {network}
```

The address of the deployed contract will be written into the `contractAddresses.json` file.

## 2. Deploy L2 Contract
```
cd {path}/genesis-contracts
npm run deploy {network}
```

The address of the deployed contract will be written into the `contractAddresses.json` file.

## 3. Build themis (or use themis image)
```
make build
```

## 4. Create network configuration file

```
./build/themisd create-network --v 3 --n 0 --output-dir ./mynetwork --home ./mynetwork/home --chain-id mainnet --chain mainnet 
```
In the `mynetwork` directory, three validator node directories will be generated, named `node0`, `node1`, and `node2`, corresponding to the three initialized node configuration directories.

## 5. ID for L1 LockingContract: 1, 2, 3

```
cd {path}/contracts
Compile `scrpts/lock.js`, `seqAddress` `seqPub`, corresponding to the addresses and public keys of `node0`, `node1`, and `node2` (which can be obtained in `mynetwork/signer-dump.json`).
Run `npm run lock {network}`.
It is necessary to ensure that the ID obtained from the on-chain lock-up is consistent with the node order.
```

## 6. Modify default configuration parameters

Configure directories for three nodes and make the following modifications:

- In themis/config/config.toml, change the moniker to the actual node identifier.
- In persistent_peers, replace node0 with the public IP of node0, node1 with the public IP of node1, and node2 with the public IP of node2.
- In themis/config/genesis.json, modify the chainmanager.chain_params parameter to the actual deployed contract addresses:
   ```
   {
        "main_chain_id": "1",   // Mainnet chain_id
        "metis_chain_id": "1088", // Metis chain chain_id
        "metis_token_address": "0x114f836434a9aa9ca584491e7965b16565bf5d7b",     // L1 Metis contract address L1MetisToken
        "staking_manager_address": "0x73d5b3d9c5502953e51e3ddedff81a3e86fa874d", // L1 locking contract address LockingPoolProxy
        "staking_info_address": "0x33cdb54fb5b0a469adb7d294dd868f4b782e2fba",   // L1 info contract address LockingInfo
        "validator_set_address": "0x3C30d5A6B4F29187122eE4142D6627B228D3b59D"   // L2 MetisSequencerSet contract address MetisSequencerSet
    }
   ```
With these modifications, the node configuration files are complete. Copy the configurations for node0, node1, and node2 to their respective servers.

## 7. Backup node files
```
cp -rf {path}/node/themis  /data/backup
```


## 8. Run tss-node

(1) Query the internal IP addresses of three TSS nodes.

```
ifconfig ens5
ens5: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 9001
        inet 10.128.6.70  netmask 255.255.240.0  broadcast 10.128.15.255
        inet6 fe80::ced:9eff:fe02:6ff9  prefixlen 64  scopeid 0x20<link>
        ether 0e:ed:9e:02:6f:f9  txqueuelen 1000  (Ethernet)
        RX packets 790251764  bytes 229166422095 (229.1 GB)
        RX errors 0  dropped 1  overruns 0  frame 0
        TX packets 867541586  bytes 277685644435 (277.6 GB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```

(2) The IDs of the three TSS nodes can be found in the genesis.json file.
```
 "mpc_set": [
        {
          "id": "16Uiu2HAm8K2zq6UU7H69zevFg7MnQw66pkcSH3nttgvyZxvDhHqn",
          "moniker": "0xf2b17653c96ed20c701d2a6020ba4a2ac46c4734",
          "key": "Ar90oR3VFVeSMle9N/4huWkc1MrV+UefjUApv1fW7fAt"
        },
        {
          "id": "16Uiu2HAmKbiLqXmhoGCLNTubuF6qSp8oQBXN51X5r2Mg69rrS4sj",
          "moniker": "0x32649704bed02d28c99005544763720183fa9fd8",
          "key": "A2cq90qQfyP8pOSIkpI0GmI0bc9BPyhEk4Njy3R3w+YC"
        },
        {
          "id": "16Uiu2HAmJh1wK9PhdXoPQeXXqcGmcCeYCS2ubQ26q3LUnrvFgnhQ",
          "moniker": "0xf6b281cfc15465bd5ec8ec682c45c248f41068f3",
          "key": "A1mrE8Fkwbj7a+5p56drPu6yto80qbLuFJXP5q0699FT"
        }
      ]
```

(3) Set p2pBootstrapPeers startup arguments
```
--p2pBootstrapPeers=/ip4/10.128.6.xxx/tcp/4001/p2p/16Uiu2HAm8K2zq6UU7H69zevFg7MnQw66pkcSH3nttgvyZxvDhHqn,/ip4/10.128.21.xxx/tcp/4001/p2p/16Uiu2HAmKbiLqXmhoGCLNTubuF6qSp8oQBXN51X5r2Mg69rrS4sj,/ip4/10.128.28.xxx/tcp/4001/p2p/16Uiu2HAmJh1wK9PhdXoPQeXXqcGmcCeYCS2ubQ26q3LUnrvFgnhQ
```

(4) Copy private key

```
cp ./node/themis/config/priv_validator_key.json  ./tss-node
```

(5) Startup tss-node
```
cd {path}/tss-node

Config docker-compose.yml
--themisRest Repalce nodeIP, tss-node0 is node0'sIP
--p2pBootstrapPeers Replace to (3)

docker-compose up -d
```


## 9. Startup node and bridge
(1) Config envs
```
cd {path}/tss-node
vim docker-compose.yml
Make sure to modify env
```

(2) Startup node
```
docker-compose up -d node && docker-compose logs -f node
```

(3) Startup bridge
```
docker-compose up -d bridge && docker-compose logs -f bridge
```
Check logs


## 10. Change default MPC address
After the bridge is successfully launched, a MPC address will be automatically generated. You can check whether it has been successfully generated by visiting http://127.0.0.1:1317/mpc/latest/0. After confirming the success, it is necessary to update the contract accordingly.

(1) Update MPC address to L1 contract

```
cd {path}/contracts
vi .env
# set MPC_ADDRESS
npm run update_mpc  sepolia
```

(2) Update MPC address to L2 contract

```
cd {path}/genesis-contracts
vi .env
# set MPC_ADDRESS
npm run update_mpc sepolia
```


## 11. Backup tss-node files

```
docker-compose stop tss-node
cp -rf tss-node  /data/backup
docker-compose up -d tss-node
```