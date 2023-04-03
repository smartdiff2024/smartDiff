# Program Readme

## Version Requirement

- go-lang: go 1.15 or later



## Usage

### Record

`./geth --datadir ./replay/atcgeth import ./import/0-1M.rlp`

where `datadir` specifies the path to generate the chain and `import` specifies the path to the `blockchain data`. The historical transactions are recorded in the same directory of the `geth`.



#### Obtain more upgrade sequences

`./substate-cli replay --impls --withoutimpl ./order.csv --withimpls ./order2.csv --substatedir ./substate.ethereum 12720347 15335369`

where `withoutimpl` specifies the `csv file` that may not include all the upgrade sequences while `withimpls` specifies the `csv file` that include extra upgrade sequences. `substatedir` specifies the path of historical transactions.



#### Obtain the position of the upgrade transaction

`./substate-cli replay --updatep --withoutimpl order.csv --withimpls order2.csv --substatedir ./substate.ethereum 12720347 15335369 `

where `withoutimpl` specifies the `csv file` that doesn't include the position of the upgrade transaction while `withimpls` specifies the `csv file` that include the position of the upgrade transaction. `substatedir` specifies the path of historical transactions.



#### Obtain the bytecodes of proxy addresses

`./substate-cli replay -getcode --withimpls ./order_result.csv --jsondir ./json --substatedir ./substate.ethereum/ 12720343 15335369`

where `withimpls` specifies the `csv file` containing the contract upgrade sequences, while `jsondir` specifies the `directory path` where the `bytecode json files` are generated. `substatedir` specifies the path of historical transactions.



#### Execute differential testing

`./substate-cli replay -dtrace --jsondir ./json/ --dtracefiles 0x0b59ef7d85f1acc791c937d2e9c40c020c156c6e --dtraceresultdir ./result --oriimpl 0x557DE75A27025815dB74E16EA2B58eb7C2a1360f --replaceimpl 0x6971C1d22cCD76D3ca706523f3685E20faef9071 --firstp 140 --lastp 58 --implblock 12720347 --implp 113 --substatedir ./substate.ethereum/ 12720343 12764589`

where `jsondir` and `dtracefiles` parameters specify the `json file` to be used, while `dtraceresultdir` specifies the directory path where the results will be generated. The `oriimpl` and `replaceimpl` parameters specify the continuous addresses to be used in the replay. The `firstp`, `lastp`, `implblock`, and `implp` parameters specify the transaction numbers to start and end the transaction replay, as well as the transaction number to begin to record the state of proxy contract.



