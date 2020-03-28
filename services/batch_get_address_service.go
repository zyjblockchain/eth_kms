package services

import "eth_kms/models"

type BatchGetAddrInfo struct {
	StartId uint `form:"startId" json:"startId"`
	Limit   uint `form:"limit" json:"limit"`
}

// BatchGetAddr 批量拉取地址
func (b *BatchGetAddrInfo) BatchGetAddr() ([]string, uint, error) {
	keysMgrArr, total, err := models.GetBatchById(b.StartId, b.Limit)
	addresses := make([]string, 0, 0)
	for _, keys := range keysMgrArr {
		addresses = append(addresses, keys.Address)
	}
	return addresses, total, err
}
