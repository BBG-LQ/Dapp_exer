#!/bin/bash

# 编译合约
solc --bin contract/Count.sol -o ./build/
# 生成合约ABI
solc --abi contract/Count.sol -o ./build/
#
abigen --abi ./build/Count.abi --bin ./build/Count.bin --pkg contract --out ./counter.go