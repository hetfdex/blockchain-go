package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"

	"github.com/hetfdex/blockchain-go/wallet"
)

const (
	minerRewardAmount = 100

	minerRewardSender = "miner_reward_sender"
)

var (
	errUpdate = errors.New("update not valid: amount exceeds balance")
)

type Transaction struct {
	ID        []byte
	TxInput   TransactionInput
	OutputMap map[string]uint64
}

type TransactionInput struct {
	CreatedAt       time.Time
	SenderAddress   []byte
	Balance         uint64
	SenderSignature []byte
}

func New(senderWallet wallet.Wallet, recipientAddress []byte, amount uint64) (Transaction, error) {
	//Wallet.calculateBalance({chain, address: this.publicKey}) -> if !balance : err
	outputMap := newOutputMap(senderWallet, recipientAddress, amount)

	txInput, err := newTxInput(senderWallet, outputMap)

	if err != nil {
		return Transaction{}, err
	}

	tx := Transaction{
		TxInput:   txInput,
		OutputMap: outputMap,
	}

	tx.setID()

	return tx, nil
}

func MinerReward(minerWallet wallet.Wallet) Transaction {
	outputMap := make(map[string]uint64)

	outputMap[string(minerWallet.PublicKey)] = minerRewardAmount

	tx := Transaction{
		TxInput: TransactionInput{
			CreatedAt:     time.Now().UTC(),
			SenderAddress: []byte(minerRewardSender),
		},
		OutputMap: outputMap,
	}

	tx.setID()

	return tx
}

func (t *Transaction) Update(senderWallet wallet.Wallet, recipientAddress []byte, amount uint64) error {
	if amount > t.OutputMap[string(senderWallet.PublicKey)] {
		return errUpdate
	}

	if t.OutputMap[string(recipientAddress)] == 0 { //nil?
		t.OutputMap[string(recipientAddress)] = amount
	} else {
		t.OutputMap[string(recipientAddress)] = t.OutputMap[string(recipientAddress)] + amount

	}
	t.OutputMap[string(senderWallet.PublicKey)] = t.OutputMap[string(senderWallet.PublicKey)] - amount

	txInput, err := newTxInput(senderWallet, t.OutputMap)

	if err != nil {
		return err
	}
	t.TxInput = txInput

	return nil
}

func (t *Transaction) Valid() bool {
	var outputTotal uint64

	for add := range t.OutputMap {
		outputTotal += t.OutputMap[add]
	}

	if outputTotal != t.TxInput.Balance {
		return false
	}

	/* if !Verify(t.TxInput.SenderAddress, t.OutputMap, t.TxInput.SenderSignature) {
		log.Println(warnValidSignature)

		return false
	} */

	return true
}

func Equal(a Transaction, b Transaction) bool {
	if !bytes.Equal(a.ID, b.ID) {
		return false
	}

	outMapA := a.OutputMap
	outMapB := b.OutputMap

	if len(outMapA) != len(outMapB) {
		return false
	}

	for add := range outMapA {
		if outMapA[add] != outMapB[add] {
			return false
		}
	}

	txInA := a.TxInput
	txInB := b.TxInput

	if !txInA.CreatedAt.Equal(txInB.CreatedAt) {
		return false
	}

	if !bytes.Equal(txInA.SenderAddress, txInB.SenderAddress) {
		return false
	}

	if txInA.Balance != txInB.Balance {
		return false
	}

	if !bytes.Equal(txInA.SenderSignature, txInB.SenderSignature) {
		return false
	}
	return true
}

func newTxInput(senderWallet wallet.Wallet, outputMap map[string]uint64) (TransactionInput, error) {
	senderSignature, err := senderWallet.Sign(outputMap)

	if err != nil {
		return TransactionInput{}, err
	}

	return TransactionInput{
		CreatedAt:       time.Now().UTC(),
		SenderAddress:   senderWallet.PublicKey,
		Balance:         senderWallet.Balance,
		SenderSignature: senderSignature,
	}, nil
}

func newOutputMap(senderWallet wallet.Wallet, recipientAddress []byte, amount uint64) map[string]uint64 {
	outputMap := make(map[string]uint64)

	outputMap[string(recipientAddress)] = amount
	outputMap[string(senderWallet.PublicKey)] = senderWallet.Balance - amount

	return outputMap
}

func (t *Transaction) setID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := json.NewEncoder(&encoded)

	err := encoder.Encode(t)

	if err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())

	t.ID = hash[:]

	return nil
}
