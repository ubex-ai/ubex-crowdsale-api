## Development only (ubuntu 16.04 tested)

source zeppelin contracts

```console
cd $GOPATH/src/ubex-api
rsync -vax vendor/github.com/OpenZeppelin/zeppelin-solidity/contracts/ solidity/zeppelin/
```

install latest solc

```console
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
```

build go bindings

```console
$GOPATH/src/ubex-api/solidity/bin/abigen --sol solidity/UbexCrowdsale.sol --pkg=ubex_crowdsale --out=solidity/bindings/ubex_crowdsale/UbexCrowdsale.go
```