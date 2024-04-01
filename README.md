# Themis

Validator node for Metis Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source 

Make sure your have go1.17+ already installed

### Install 
```bash 
$ make install
```
### Init-themis 
```bash 
$ themisd init
$ themisd init --chain=mainnet        Will init with genesis.json for mainnet
$ themisd init --chain=testnet         Will init with genesis.json for testnet
```
### Run-themis 
```bash 
$ themisd start
```
#### Usage
```
$ themisd start                       Will start for mainnet by default
$ themisd start --chain=mainnet       Will start for mainnet
$ themisd start --chain=testnet       Will start for testnet
$ themisd start --chain=local         Will start for local with NewSelectionAlgoHeight = 0
```

### Run rest server
```bash 
$ themisd rest-server 
```

### Run bridge
```bash 
$ themisd bridge 
```

### Develop using Docker

You can build and run Themis using the included Dockerfile in the root directory:

```bash
docker build -t themis .
docker run themis
```


