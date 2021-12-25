package transaction

type TxInputValidator interface {
	ValidSignature(string) bool
}

type TxInput struct {
	ID          []byte `json:"id"`
	OutputIndex uint64 `json:"output_index"`
	Signature   string `json:"signature"`
}

func (i *TxInput) ValidSignature(data string) bool {
	return i.Signature == data
}
