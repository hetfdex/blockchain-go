package wallet

import "encoding/json"

const (
	startingBalance = 1000
)

type Wallet struct {
	Balance   uint64
	Keypair   []byte
	PublicKey []byte
}

func New() Wallet {
	kp := []byte("tbd") //ec.genKeyPair()

	return Wallet{
		Balance:   startingBalance,
		Keypair:   kp,
		PublicKey: kp, //keyPair.getPublic().encode("hex")
	}
}

func (w *Wallet) Sign(outputMap map[string]uint64) ([]byte, error) {
	res, err := json.Marshal(outputMap)

	if err != nil {
		return nil, err
	}
	return res, nil //this.keyPair.sign(cryptoHash(data))
}

//func CalculateBalance() {} -> here?
