// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  contract
 * @Author: Gong Zibin
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"math/big"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
	log "github.com/cihub/seelog"
)



//NewNativeContract is to create a new native contract
func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	intf, err := NewNativeContractHandler()
	if err != nil {
		return nil, err
	}
	roleIntf.SetScheduleDelegate(&role.ScheduleDelegate{CurrentTermTime: big.NewInt(2)})

	return intf, nil
}

func newTransaction(contract string, method string, param []byte) *types.Transaction {
	trx := &types.Transaction{
		Sender:   contract,
		Contract: contract,
		Method:   method,
		Param:    param,
	}

	return trx
}

//NativeContractInitChain is to init
func NativeContractInitChain(ldb *db.DBService, roleIntf role.RoleInterface, ncIntf NativeContractInterface) ([]*types.Transaction, error) {
	err := CreateNativeContractAccount(roleIntf)
	if err != nil {
		return nil, err
	}

	var trxs []*types.Transaction

	// construct trxs
	var i int
	Abi := abi.GetAbi()
	
	for i = 0; i < len(config.Genesis.InitDelegates); i++ {
		name := config.Genesis.InitDelegates[i].Name
		
		// 1, new account trx
		mapstruct := make(map[string]interface{})
		abi.Setmapval(mapstruct, "name", name)
		abi.Setmapval(mapstruct, "pubkey", config.Genesis.InitDelegates[i].PublicKey)
		nparam, err2   := abi.MarshalAbiEx(mapstruct, Abi, config.BOTTOS_CONTRACT_NAME, "newaccount")
		if err2 != nil {
			log.Info("abi.MarshalAbiEx failed for new account:", name)
			continue
		}

		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx

		mapstruct2 := make(map[string]interface{})
		abi.Setmapval(mapstruct2, "from", config.BOTTOS_CONTRACT_NAME)
		abi.Setmapval(mapstruct2, "to", name)
		abi.Setmapval(mapstruct2, "value", uint64(config.Genesis.InitDelegates[i].Balance))
		tparam, err3   := abi.MarshalAbiEx(mapstruct2, Abi, config.BOTTOS_CONTRACT_NAME, "transfer")
		if err3 != nil {
			log.Info("abi.MarshalAbiEx failed for transfer with account:", name)
			continue
		}

		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate

		mapstruct3 := make(map[string]interface{})
		abi.Setmapval(mapstruct3, "name", name)
		abi.Setmapval(mapstruct3, "pubkey", config.Genesis.InitDelegates[i].PublicKey)
		sparam, err4   := abi.MarshalAbiEx(mapstruct3, Abi, config.BOTTOS_CONTRACT_NAME, "setdelegate")
		if err4 != nil {
			log.Info("abi.MarshalAbiEx failed for setdegelage with account:", name)
			continue
		}
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "setdelegate", sparam)
		trxs = append(trxs, trx)
	}

	// init CoreState delegates
	coreState, _ := roleIntf.GetCoreState()
	for i = 0; i < int(config.BLOCKS_PER_ROUND); i++ {
		name := config.Genesis.InitDelegates[i].Name

		coreState.CurrentDelegates = append(coreState.CurrentDelegates, name)
	}
	roleIntf.SetCoreState(coreState)

	return trxs, nil
}

//CreateNativeContractAccount is to create native contract account
func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	pubkey, _ := common.HexToBytes(config.Param.KeyPairs[0].PublicKey)
	a := abi.CreateNativeContractABI()
	abijson, _ := abi.AbiToJson(a)
	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  config.Genesis.GenesisTime,
		PublicKey:   pubkey,
		ContractAbi: []byte(abijson),
	}
	roleIntf.SetAccount(bto.AccountName, bto)

	// balance
	var initSupply uint64
	initSupply, err = safemath.Uint64Mul(config.BOTTOS_INIT_SUPPLY, config.BOTTOS_SUPPLY_MUL)
	if err != nil {
		return err
	}

	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance:     initSupply,
	}
	roleIntf.SetBalance(bto.AccountName, balance)

	// staked_balance
	stakedBalance := &role.StakedBalance{
		AccountName: bto.AccountName,
	}
	roleIntf.SetStakedBalance(bto.AccountName, stakedBalance)

	return nil
}
