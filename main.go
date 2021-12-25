package main

import (
	"fmt"
	"os"

	"github.com/hetfdex/blockchain-go/orchestrator"
)

func main() {
	_ = os.RemoveAll("./tmp")

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
		err = orchestrator.AddBlock(bc, fmt.Sprintf("%d", i))

		if err != nil {
			panic(err)
		}
	}

	err = orchestrator.PrintBlocks(bc)

	if err != nil {
		panic(err)
	}
	_ = os.RemoveAll("./tmp")
}
