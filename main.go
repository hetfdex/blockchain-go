package main

type Block struct {
	Data     []byte
	Hash     []byte
	PrevHash []byte
}

type BlockChain struct {
	blocks []*Block
}
