

package transaction

import (
	"time"
	"sync"
	"fmt"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/action/message"
	
)


var (
	trxExpirationCheckInterval    = time.Minute     // Time interval for check expiration pending transactions
	trxExpirationTime             = time.Minute     // Pending Trx max time , to be delete
)



type TrxPool struct {
	pending     map[common.Hash]*types.Transaction       
	expiration  map[common.Hash]time.Time    // to be delete
	
	mu           sync.RWMutex
	quit chan struct{}
}


func InitTrxPool() *TrxPool {
	
	// Create the transaction pool
	pool := &TrxPool{
		pending:      make(map[common.Hash]*types.Transaction),
		expiration:   make(map[common.Hash]time.Time),
		quit:         make(chan struct{}),
	}

	go pool.expirationCheckLoop()

	return pool
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (pool *TrxPool) expirationCheckLoop() {	
	expire := time.NewTicker(trxExpirationCheckInterval)
	defer expire.Stop()

	for {
		select {
		case <-expire.C:
			pool.mu.Lock()

			var currentTime = time.Now()
			for txHash := range pool.expiration {

				if (currentTime.After(pool.expiration[txHash])) {
					delete(pool.expiration, txHash)
					delete(pool.pending, txHash)					
				}
				
			}
			pool.mu.Unlock()

		case <-pool.quit:
			return
		}
	}
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (pool *TrxPool) addTransaction(trx *types.Transaction) {	
	trxHash := trx.Hash()
	pool.pending[trxHash] = trx
	//pool.expiration = time.Now()
}



func (pool *TrxPool) Stop() {
	
	close(pool.quit)

	fmt.Println("Transaction pool stopped")
}

func (pool *TrxPool)CheckTransactionBaseConditionFromFront(){

	/* check max pending trx num */
	/* check account validate */
	/* check signature */

}


func (pool *TrxPool)CheckTransactionBaseConditionFromP2P(){	

}



// HandlTransactionFromFront handles a transaction from front
func (pool *TrxPool)HandleTransactionFromFront(trx *types.Transaction) {
	
    pool.CheckTransactionBaseConditionFromFront()
	//start db session
	ApplyTransaction(trx)

	//add to pending

	//revert db session

	//tell P2P actor to notify trx	
}


// HandlTransactionFromP2P handles a transaction from P2P
func (pool *TrxPool)HandleTransactionFromP2P(trx *types.Transaction) {

	pool.CheckTransactionBaseConditionFromP2P()

	// start db session
	ApplyTransaction(trx)	

	pool.addTransaction(trx)

	//revert db session	
}



func (pool *TrxPool)HandlePushTransactionReq(TrxSender message.TrxSenderType, trx *types.Transaction){

	if (message.TrxSenderTypeFront == TrxSender){ 
		pool.HandleTransactionFromFront(trx)
	} else if (message.TrxSenderTypeP2P == TrxSender) {
		pool.HandleTransactionFromP2P(trx)
	}	
}
