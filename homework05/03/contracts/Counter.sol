// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title Counter - 简单的计数器合约
/// @notice 允许任何人递增计数器并查询当前值
contract Counter {
    // 计数器存储在 slot 0
    uint256 private _count;
    // 合约部署者地址，存储在 slot 1
    address public owner;

    /// @notice 每次递增时触发的事件
    event CountIncremented(address indexed incrementer, uint256 newCount);

    constructor() {
        owner = msg.sender;
        _count = 0;
    }

    /// @notice 将计数器加 1，并触发事件
    function increment() public {
        _count += 1;
        emit CountIncremented(msg.sender, _count);
    }

    /// @notice 返回当前计数器值
    function getCount() public view returns (uint256) {
        return _count;
    }
}
