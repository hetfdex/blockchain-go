package transaction

type TxOutputValidator interface {
	ValidPublicKey(string) bool
}

type TxOutput struct {
	Value     uint64 `json:"value"`
	PublicKey string `json:"public_key"`
}

func (o *TxOutput) ValidPublicKey(data string) bool {
	return o.PublicKey == data
}
