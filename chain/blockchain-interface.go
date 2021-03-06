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
 * file description:  blockchain general interface and logic
 * @Author: Gong Zibin
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

import (
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
)

const (
	//InsertBlockSuccess insert block successfully
	InsertBlockSuccess uint32 = 0
	//InsertBlockErrorGeneral general error
	InsertBlockErrorGeneral uint32 = 1
	//InsertBlockErrorNotLinked the block not linked to the chain
	InsertBlockErrorNotLinked uint32 = 2
	//InsertBlockErrorValidateFail block validate fail
	InsertBlockErrorValidateFail uint32 = 3
	//InsertBlockErrorDiffLibLinked different lib block but linked
	InsertBlockErrorDiffLibLinked uint32 = 4
	//InsertBlockErrorDiffLibNotLinked different lib block and not linked in this chain
	InsertBlockErrorDiffLibNotLinked uint32 = 5
)

//HandledBlockCallback call back
type HandledBlockCallback func(*types.Block)

//BlockChainInterface the interface of chain
type BlockChainInterface interface {
	Close()

	HasBlock(hash common.Hash) bool
	GetBlockByHash(hash common.Hash) *types.Block
	GetBlockByNumber(number uint32) *types.Block
	GetHeaderByNumber(number uint32) *types.Header

	HeadBlockTime() uint64
	HeadBlockNum() uint32
	HeadBlockHash() common.Hash
	HeadBlockDelegate() string
	LastConsensusBlockNum() uint32
	GenesisTimestamp() uint64

	ValidateBlock(block *types.Block) uint32
	InsertBlock(block *types.Block) uint32

	RegisterHandledBlockCallback(cb HandledBlockCallback)
}
