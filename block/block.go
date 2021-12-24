package block

const (
	genesisBlockData = "genesis_block_data"
)

type Block struct {
	Data     []byte `json:"data"`
	PrevHash []byte `json:"prev_hash"`
	Hash     []byte `json:"hash"`
	Nonce    uint64 `json:"nonce"`
}

func New(data string, prevHash []byte) *Block {
	return &Block{
		Data:     []byte(data),
		PrevHash: prevHash,
		Hash:     []byte{},
		Nonce:    0,
	}
}

func NewGenesis() *Block {
	return New(genesisBlockData, []byte{})
}
