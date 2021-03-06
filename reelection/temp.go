// Copyright 2018 The MATRIX Authors as well as Copyright 2014-2017 The go-ethereum Authors
// This file is consisted of the MATRIX library and part of the go-ethereum library.
//
// The MATRIX-ethereum library is free software: you can redistribute it and/or modify it under the terms of the MIT License.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, 
//and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject tothe following conditions:
//
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, 
//WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISINGFROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package reelection

import (
	"errors"
	"math/big"
	"math/rand"
	"sync/atomic"

	"github.com/matrix/go-matrix/ca"
	"github.com/matrix/go-matrix/core/vm"

	"github.com/matrix/go-matrix/common"
	"github.com/matrix/go-matrix/core"
	"github.com/matrix/go-matrix/log"
	"github.com/matrix/go-matrix/mc"
	"github.com/matrix/go-matrix/p2p/discover"
)

func (self *ReElection) TopoUpdate(Q0, Q1, Q2 []mc.TopologyNodeInfo, netTop mc.TopologyGraph, offline []common.Address) []mc.Alternative {
	return self.elect.ToPoUpdate(Q0, Q1, Q2, netTop, offline)
}
func (self *ReElection) NativeUpdate(Q0, Q1, Q2 []mc.TopologyNodeInfo, online mc.TopologyNodeInfo, flag int) ([]mc.TopologyNodeInfo, []mc.TopologyNodeInfo, []mc.TopologyNodeInfo) {
	return self.elect.PrimarylistUpdate(Q0, Q1, Q1, online, flag)
}

//得到特殊交易
func getKeyTransInfo(Height uint64, types string) map[common.Address][]byte {
	//asd()
	ans, err := core.GetBroadcastTxs(big.NewInt(int64(Height)), types)
	if err != nil {
		log.Error(Module, "获取特殊交易失败 Height", Height, "types", types)
	}
	return ans
}

func GetCurrentTopology(height uint64, reqtypes common.RoleType) (*mc.TopologyGraph, error) {

	return ca.GetTopologyByNumber(reqtypes, height)
}

var count int32

func getrand() int {
	var flag bool
	flag = true
	var num int

	atomic.AddInt32(&count, 1)
	q := rand.Intn(15000)
	if flag == true {
		num = 30000 + q
		flag = false
	}
	if flag == false {
		num = 30000 - q
		flag = true
	}

	return num
}

type TestELECT struct {
	Deposit []vm.DepositDetail
	Add     []common.Address
}

var Dminer TestELECT
var Dalidator TestELECT

func Test() {
	Dminer.Deposit, Dminer.Add = GetDepost()
	Dalidator.Deposit, Dalidator.Add = GetDepost()

}

func GetDepost() ([]vm.DepositDetail, []common.Address) {
	var mmreqm mc.MasterMinerReElectionReqMsg
	var item vm.DepositDetail
	//	var str string
	var NodeID discover.NodeID
	var Address common.Address
	mmreqm.SeqNum = 20
	mmreqm.RandSeed = big.NewInt(int64(30))

	for i := 0; i < 64; i++ {
		NodeID[i] = byte(i)
	}

	for k := 0; k < 20; k++ {
		Address[k] = byte(k)
	}
	address_scf := make([]common.Address, 0)
	for j := 0; j < 3; j++ {
		NodeID[63] = byte(int(NodeID[63]) + j)
		Address[19] = byte(int(Address[19]) + j)

		item.Address = Address //common.BytesToAddress([]byte("32f7c8dae96abdae96ab")) //[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
		item.Deposit = big.NewInt(int64(getrand() - j))
		item.NodeID = NodeID //discover.NodeID("discoverNodeIdBytestoaddress") //discover.NodeID.BytesToAddress({1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
		item.OnlineTime = big.NewInt(int64(300 - j))
		mmreqm.MinerList = append(mmreqm.MinerList, item)
		address_scf = append(address_scf, item.Address)
	}

	return mmreqm.MinerList, address_scf
}
func GetAllElectedByHeight(Heigh *big.Int, tp common.RoleType) ([]vm.DepositDetail, error) {

	switch tp {
	case common.RoleMiner:
		//todo：获取所有的抵押交易
		ans, err := ca.GetElectedByHeightAndRole(Heigh, common.RoleMiner)
		log.INFO("從CA獲取礦工抵押交易", "data", ans, "height", Heigh)
		if err != nil {
			return []vm.DepositDetail{}, errors.New("获取矿工交易身份不对")
		}
		return ans, nil
	case common.RoleValidator:
		ans, err := ca.GetElectedByHeightAndRole(Heigh, common.RoleValidator)
		log.Info("從CA獲取驗證者抵押交易", "data", ans, "height", Heigh)
		if err != nil {
			return []vm.DepositDetail{}, errors.New("获取验证者交易身份不对")
		}
		return ans, nil

	default:
		return []vm.DepositDetail{}, errors.New("获取抵押交易身份不对")
	}
}

func GetFound() []vm.DepositDetail {
	return []vm.DepositDetail{}
}
