// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Count{
    uint256 public count;

    address public owner;


    constructor(uint256 _initcount){
        owner=msg.sender;
        count=_initcount;
    }

    modifier onlyOwer(){
        require(msg.sender == owner, "Counter: caller is not the owner");
        _;
    }

    function addOne() external returns (uint256)  {

        count += 1;
        return count;
    }

    // 
    function setCount(uint256 _number) external onlyOwer{

        count = _number;
        
    }

    // 
    function reset() external onlyOwer{

        count = 0;  
    }


    function getCount() external view returns (uint256) {
        return count;
    }




    



}