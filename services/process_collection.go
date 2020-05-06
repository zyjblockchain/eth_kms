package services

import (
	"errors"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/zyjblockchain/eth_kms/collection"
	"github.com/zyjblockchain/eth_kms/common"
	"github.com/zyjblockchain/eth_kms/models"
	"github.com/zyjblockchain/eth_kms/transaction"
	"github.com/zyjblockchain/sandy_log/log"
	"os"
	"strconv"
	"sync"
)

type ProcessCollection struct {
	sync.Mutex
}

// ProcessCollectUSDT USDT资产归集
func (p *ProcessCollection) ProcessCollectUSDT() error {
	p.Lock()
	defer p.Unlock()

	var offset int64
	// 获取上次归集到的offset
	strNum, err := models.Get(models.CollectionAddrOffset)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果err为not fund, 则offset从0开始
			offset = 0
		} else {
			log.Errorf("在kv表中查询collection offset 失败：%v", err)
			return err
		}
	} else {
		offset, err = strconv.ParseInt(strNum, 10, 32)
		if err != nil {
			log.Errorf("strconv.ParseInt err: %v", err)
			return err
		}
	}

	// 1. 筛选出需要归集的地址
	log.Infof("ethereum chain net: %s", os.Getenv("ETH_NET"))
	client := transaction.NewEthClient(os.Getenv("ETH_NET"))
	if client == nil {
		panic("new wth client err")
	}
	defer client.Close()
	contractAddress := ethCommon.HexToAddress(os.Getenv("USDT_ADDRESS"))
	fetchAddrNum, filterAddrNum, err := collection.BatchGetCanCollectAddress(contractAddress, uint(offset), common.BatchLimit, client)
	if err != nil {
		log.Errorf("批量筛选需要归集地址失败：%v", err)
		return err
	}
	// 移动偏移量并保存到数据库
	if fetchAddrNum != 0 {
		offset += int64(fetchAddrNum)
		// 保存数据量
		if err := models.Set(models.CollectionAddrOffset, strconv.FormatInt(offset, 10)); err != nil {
			log.Errorf("保存CollectionAddrOffset到kv数据库err: %v", err)
			return err
		}
	}

	// 2. 为筛选出来需要归集的地址转gas费用
	if filterAddrNum == 0 {
		// 没有需要归集的地址，则直接返回
		log.Errorf("没有需要归集的地址")
		return errors.New("没有需要归集的地址")
	}
	// 获取gas fee的上次偏移位置
	var gasFeeOffset int64
	str, err := models.Get(models.SendGasFeeOffset)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			gasFeeOffset = 0
		} else {
			log.Errorf("在kv表中查询SendGasFeeOffset offset 失败：%v", err)
			return err
		}
	} else {
		gasFeeOffset, err = strconv.ParseInt(str, 10, 32)
		if err != nil {
			log.Errorf("strconv.ParseInt err: %v", err)
			return err
		}
	}
	// 为归集地址转gas费用
	if err := collection.SendGasFeeForColAddrProcess(uint(gasFeeOffset), uint(filterAddrNum), client); err != nil {
		log.Errorf("为归集地址转gas fee失败：%v", err)
		return err
	}

	// 3. 归集到接收地址
	if err := collection.ColleTokenProcess(uint(gasFeeOffset), uint(filterAddrNum), client); err != nil {
		log.Errorf("归集到接收地址失败：%v", err)
		return err
	}

	// 4. 移动偏移量并保存数据库
	gasFeeOffset += int64(filterAddrNum)
	if err := models.Set(models.SendGasFeeOffset, strconv.FormatInt(gasFeeOffset, 10)); err != nil {
		log.Errorf("保存gasFeeOffset 到kv表失败：%v", err)
		return err
	}
	return nil
}
