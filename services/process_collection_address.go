package services

import "github.com/zyjblockchain/eth_kms/models"

// 归集资产接收地址
type SetCollectionAddr struct {
	Address string `form:"address" json:"address"`
}

// SetCollectionAddr
func (s *SetCollectionAddr) SetCollectionAddr() error {
	return models.Set(models.CollectionAddrKey, s.Address)
}

type GetCollectionAddr struct {
}

// GetCollectionAddr
func (g *GetCollectionAddr) GetCollectionAddr() (string, error) {
	return models.Get(models.CollectionAddrKey)
}
