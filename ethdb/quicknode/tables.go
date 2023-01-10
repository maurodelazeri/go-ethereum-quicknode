package quicknode

type EvmData struct {
	Key   []byte `gorm:"primary_key"`
	Value []byte
}
