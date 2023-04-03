# Program Readme

## Version Requirement 

- IDA: IDA 7.2 or IDA7.3
- ida-evm: lateset
- idaphora: diaphora-1.2



## Usage

`check`: generate a result db from `bytecode.json`, it will split the `bytecode.json ` into some `bytecode` files which can be decompiled by `IDA` with `ida-evm`  and invoke `diaphora` to compare the decompiled result.

`analy_check_result`: Analyze similarity results and determine whether the results comply with refactoring consistency.