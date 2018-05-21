package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bottos-project/core/api"
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/chain/extra"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/role"

	"github.com/bottos-project/core/contract"
	"github.com/bottos-project/core/contract/contractdb"

	cactor "github.com/bottos-project/core/action/actor"
	caapi "github.com/bottos-project/core/action/actor/api"
	"github.com/bottos-project/core/action/actor/transaction"
	actionenv "github.com/bottos-project/core/action/env"
	"github.com/bottos-project/core/transaction"
	"github.com/micro/go-micro"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Load config fail")
		os.Exit(1)
	}

	dbInst := db.NewDbService(config.Param.DataDir, filepath.Join(config.Param.DataDir, "blockchain"), config.Param.OptionDb)
	if dbInst == nil {
		fmt.Println("Create DB service fail")
		os.Exit(1)
	}

	roleIntf := role.NewRole(dbInst)
	contractDB := contractdb.NewContractDB(dbInst)

	nc, err := contract.NewNativeContract(roleIntf)
	if err != nil {
		fmt.Println("Create Native Contract error: ", err)
		os.Exit(1)
	}

	chain, err := chain.CreateBlockChain(dbInst, roleIntf, nc)
	if err != nil {
		fmt.Println("Create BlockChain error: ", err)
		os.Exit(1)
	}

	txStore := txstore.NewTransactionStore(chain, roleIntf)

	actorenv := &actionenv.ActorEnv{
		RoleIntf:   roleIntf,
		ContractDB: contractDB,
		Chain:      chain,
		TxStore:    txStore,
		NcIntf:     nc,
	}
	cactor.InitActors(actorenv)
	//caapi.PushTransaction(2876568)

	//caapi.InitTrxActorAgent()
	var trxPool = transaction.InitTrxPool(actorenv)
	trxactor.SetTrxPool(trxPool)

	if config.Param.ApiServiceEnable {
		repo := caapi.NewApiService(actorenv)

		service := micro.NewService(
			micro.Name("core"),
			micro.Version("2.0.0"),
		)

		service.Init()
		api.RegisterCoreApiHandler(service.Server(), repo)
		if err := service.Run(); err != nil {
			panic(err)
		}
	}

	WaitSystemDown()
}

func WaitSystemDown() {
	exit := make(chan bool, 0)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	defer signal.Stop(sigc)

	go func() {
		<-sigc
		fmt.Println("System shutdown")
		close(exit)
	}()

	<-exit
}
