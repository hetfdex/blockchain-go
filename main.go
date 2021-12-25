package main

import (
	"github.com/hetfdex/blockchain-go/orchestrator"
	"github.com/hetfdex/blockchain-go/transaction"
)

func main() {
	db, wrapper, err := orchestrator.InitDb()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	bc, err := orchestrator.InitBlockchain(wrapper)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		err = orchestrator.AddBlock(bc, []transaction.Transaction{})

		if err != nil {
			panic(err)
		}
	}

	err = orchestrator.PrintBlocks(bc)

	if err != nil {
		panic(err)
	}
}
