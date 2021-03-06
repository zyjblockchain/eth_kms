package collection

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	localcommon "github.com/zyjblockchain/eth_kms/common"
	"github.com/zyjblockchain/eth_kms/models"
	"github.com/zyjblockchain/eth_kms/transaction"
	"github.com/zyjblockchain/sandy_log/log"
	"math/big"
	"sync"
	"time"
)

type AddrTokenInfo struct {
	Address      common.Address
	TokenAddress common.Address
	TokenBalance *big.Int
}

// 1. BatchGetCanCollectAddress批量获取需要归集的地址
func BatchGetCanCollectAddress(contractAddress common.Address, startId, limit uint, client *ethclient.Client) (int, int, error) {
	// 批量拉取表中的数据
	keysMgrs, total, err := models.GetBatchById(startId, limit)
	if err != nil {
		return 0, 0, err
	}
	log.Debugf("kms中的total: %d", total)
	// 遍历KeysMgrs获取有余额的地址
	tokenAddres := make([]*AddrTokenInfo, 0)
	for _, val := range keysMgrs {
		addr := common.HexToAddress(val.Address)
		// 查询地址余额
		balance, err := transaction.GetTokenBalance(addr, contractAddress, client)
		if err != nil {
			// todo 目前对错误不处理， 直接跳过
			log.Errorf("查询地址token余额失败: %v", err)
			continue
		}
		// 判断balance是否为0
		if balance.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		// 有余额则记录
		tokenAddr := &AddrTokenInfo{
			Address:      addr,
			TokenAddress: contractAddress,
			TokenBalance: balance,
		}
		// 记录地址信息
		tokenAddres = append(tokenAddres, tokenAddr)
	}
	// 持久化到数据库
	err = saveByBatch(tokenAddres)
	return len(keysMgrs), len(tokenAddres), err
}

// 2. 批量保存需要归集的地址信息到数据库
func saveByBatch(addInfos []*AddrTokenInfo) error {
	if len(addInfos) < 1 {
		return errors.New("带保存的数据为空")
	}
	records := make([]*models.CollectionAddress, 0)
	for _, val := range addInfos {
		needCollectionAddr := &models.CollectionAddress{
			Address:      val.Address.Hex(),
			State:        0,
			TokenBalance: val.TokenBalance.String(),
			TokenAddress: val.TokenAddress.Hex(),
		}
		records = append(records, needCollectionAddr)
	}
	// 调用批量保存到数据库的接口
	err := models.InsertBatch(records)
	if err != nil {
		log.Errorf("批量插入数据库(CollectionAddress) 失败: %v", err)
	}
	return err
}

// 3. 为归集的地址转gas费用 todo 现在没有时间优化
func SendGasFeeForColAddrProcess(startId, limit uint, client *ethclient.Client) error {
	collAddres, total, err := models.GetBatchAddrById(startId, limit) // todo 优化为只拉取state等于0的记录
	if err != nil {
		log.Errorf("1. SendGasFeeForColAddrProcess error: %v", err)
	}
	log.Infof("2. CollectionAddress total: %d", total)
	// 为状态等于0的记录中的地址发eth
	collectionAddr, err := models.Get(models.CollectionAddrKey)
	if err != nil {
		log.Errorf("3. 从数据库中查询collectionAddress error: %v", err)
		return err
	}
	// 获取gas 发送地址
	gasFromAddr, err := models.Get(models.GasFromAddrKey)
	if err != nil {
		log.Errorf("4. 从数据库中查询gasFromAddr error: %v", err)
		return err
	}
	// 从kms中获取gas 发送地址私钥
	accKey, err := localcommon.GetAccountKeyByAddr(gasFromAddr)
	if err != nil {
		log.Errorf("5. 从kms中获取gasFromAddr对应的私钥 error: %v", err)
		return err
	}
	gasFromAddrPriv := accKey.Private
	balance, _ := client.BalanceAt(context.Background(), common.HexToAddress(gasFromAddr), nil)
	log.Warnf("gasFromAddr balance: %s", balance.String())

	gasFromAddrNonce, err := client.NonceAt(context.Background(), common.HexToAddress(gasFromAddr), nil)
	if err != nil {
		log.Errorf("6. 获取gasFromAddr Nonce error: %v", err)
		return err
	}
	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Errorf("7. 获取suggest gasPrice 失败： %v", err)
		suggestGasPrice = big.NewInt(500000000) // 默认5Gwei
	}
	var wg sync.WaitGroup
	for _, val := range collAddres {
		if val.State == 0 {
			// 预估交易gas
			tokenAmount, _ := new(big.Int).SetString(val.TokenBalance, 10)
			from := common.HexToAddress(val.Address)
			contractAddress := common.HexToAddress(val.TokenAddress)
			receiver := common.HexToAddress(collectionAddr)
			estGas, err := transaction.EstimateTokenTxGas(client, tokenAmount, from, contractAddress, receiver)
			if err != nil {
				log.Errorf("8. 调用预估交易gas失败：%v", err)
				estGas = 60000
			}
			// 计算需要的gas eth
			needGasEth := new(big.Int).Mul(suggestGasPrice, big.NewInt(int64(estGas)))
			// 保存待发送的交易到交易表中
			// 把交易保存在交易表中
			txRecord := &models.Transaction{
				State:  0,
				TxType: 0,
				From:   gasFromAddr,
				To:     val.Address,
				Amount: needGasEth.String(),
			}
			if err := txRecord.Add(); err != nil {
				log.Errorf("9. 保存交易记录失败：%v", err)
				return err
			}
			// 发送交易
			tx, err := transaction.SendEthTx(gasFromAddrPriv, gasFromAddrNonce, uint64(22000), suggestGasPrice, common.HexToAddress(val.Address), needGasEth, client)
			if err != nil {
				log.Errorf("10. 发送eth 转gas交易error: %v", err)
				// 更新记录为发送失败
				_ = txRecord.Update(models.Transaction{TxHash: tx.Hash().Hex(), State: 2, ErrMsg: err.Error()})
				return err
			} else {
				// 增加from的nonce
				gasFromAddrNonce++
			}
			// 更新记录
			txInfo, err := tx.MarshalJSON()
			if err != nil {
				log.Errorf("marshal tx err: %v", err)
				return err
			}
			if err := txRecord.Update(models.Transaction{TxHash: tx.Hash().Hex(), State: 1, TxInfo: string(txInfo)}); err != nil {
				log.Errorf("11. 更新transaction数据库error: %v", err)
				return err
			}
			// 开启一个协成来监听交易是否发送成功
			wg.Add(1)
			go func() {
				defer wg.Done()
				count := 0
				for {
					if count > 12 {
						// 如果2分钟都没有查到则状态设置为超时
						if err := txRecord.Update(models.Transaction{State: 4}); err != nil {
							log.Errorf("12. 交易监听超时，交易hash: %s", tx.Hash().Hex())
							return
						}
						return
					}
					count++
					time.Sleep(15 * time.Second) // 休眠15s之后去链上查交易
					log.Warnf("13. 开始监听交易 txHash: %s", tx.Hash().Hex())
					_, isPending, err := client.TransactionByHash(context.Background(), tx.Hash())
					if (err != nil && err == ethereum.NotFound) || isPending {
						// 休眠10s之后再查找
						time.Sleep(10 * time.Second)
						continue
					} else if err != nil && err != ethereum.NotFound {
						// 报错，修改数据库的交易记录为发送失败状态和失败msg
						if err := txRecord.Update(models.Transaction{State: 2, ErrMsg: err.Error()}); err != nil {
							log.Errorf("14. 更新发送交易失败Transaction数据库error: %v", err)
							return
						}
					}

					// err == nil 并且isPending == false则表示交易已经上链
					if err == nil && !isPending {
						log.Infof("15. 监听到交易记录，并更新数据库 txHash: %s", tx.Hash().Hex())
						// 更新对应交易的状态 todo 和下面的表更新优化成事务更新方式
						if err := txRecord.Update(models.Transaction{State: 3}); err != nil {
							log.Errorf("16. 更新发送交易成功状态Transaction数据库error: %v", err)
							return
						}
						// 更新对应的接收address的eth余额
						balance, err := client.BalanceAt(context.Background(), *tx.To(), nil)
						if err != nil {
							log.Errorf("17. 获取地址余额失败，addr: %s, error: %v", (*tx.To()).Hex(), err)
							return
						}
						// 更新数据库，eth余额和状态置位可以进行归集操作，归集需要设置的gasPrice和gasLimit一并更新到数据库
						if err := models.UpdateInfo(models.CollectionAddress{Address: (*tx.To()).Hex(), State: 1, EthBalance: balance.String(), GasPrice: suggestGasPrice.String(), GasLimit: estGas}); err != nil {
							log.Errorf("18. 更新CollectionAddress数据库失败，addr: %s, error: %v", (*tx.To()).Hex(), err)
							return
						}
						return
					}
				}
			}()
		}
	}
	wg.Wait()
	return nil
}

// 4. 归集到接收地址
func ColleTokenProcess(startId, limit uint, client *ethclient.Client) error {
	// 获取接收归集token的地址
	colAddr, err := models.Get(models.CollectionAddrKey)
	if err != nil {
		log.Errorf("0. 从数据库中读取CollectionAddrKey error: %v", err)
	}
	log.Infof("1. 归集接收地址：%s", colAddr)
	// 从CollectionAddress表中读取还未归集的地址
	collAddrInfoes, total, err := models.GetBatchAddrById(startId, limit)
	if err != nil {
		log.Errorf("2. 批量获取CollectionAddress error: %v", err)
	}
	log.Infof("total: %d", total)

	var wg sync.WaitGroup
	// 筛选出可以进行归集的记录进行归集操作
	for _, val := range collAddrInfoes {
		if val.State == 1 { // 只对状态为1的记录进行归集操作
			log.Infof("3. 准备开始归集的record: %s", val.Println())
			// 发送token交易到接收归集的账户中
			fromAddr := common.HexToAddress(val.Address)
			// 从kms中获取private
			ak, err := localcommon.GetAccountKeyByAddr(val.Address)
			if err != nil {
				log.Errorf("4. 从kms中获取私钥失败，address: %s, error: %v", val.Address, err)
				continue
			}
			fromPriv := ak.Private
			tokenReceiver := common.HexToAddress(colAddr)
			contractAddr := common.HexToAddress(val.TokenAddress)
			nonce, err := client.NonceAt(context.Background(), fromAddr, nil)
			if err != nil {
				log.Errorf("5. 获取nonce失败，address: %s, error: %v", fromAddr.Hex(), err)
				continue
			}
			gasPrice, b := new(big.Int).SetString(val.GasPrice, 10)
			if !b {
				log.Errorf("6. gasPrice转big Int失败, string.GasPrice: %s", val.GasPrice)
				continue
			}
			tokenAmount, b := new(big.Int).SetString(val.TokenBalance, 10)
			if !b {
				log.Errorf("7. TokenBalancee转big Int失败, string.TokenBalance: %s", val.GasPrice)
				continue
			}

			// 保存待发送的交易到交易表中
			// 把交易保存在交易表中
			txRecord := &models.Transaction{
				State:  0,
				TxType: 1,
				From:   val.Address,
				To:     colAddr,
				Amount: tokenAmount.String(),
			}
			if err := txRecord.Add(); err != nil {
				log.Errorf("8. 保存交易记录失败：%v", err)
				return err
			}
			// 设置归集表对应的记录为归集进行中
			if err := models.UpdateInfo(models.CollectionAddress{Address: fromAddr.Hex(), State: 2}); err != nil {
				log.Errorf("8.1. 更新CollectionAddress数据库失败，addr: %s, error: %v", fromAddr.Hex(), err)
				return err
			}

			// 发送交易
			tx, err := transaction.SendTokenTx(fromPriv, nonce, val.GasLimit, gasPrice, tokenReceiver, contractAddr, tokenAmount, client)
			if err != nil {
				log.Errorf("9. 发送txType = %d 交易error: %v", txRecord.TxType, err)
				// 更新记录为发送失败
				_ = txRecord.Update(models.Transaction{TxHash: tx.Hash().Hex(), State: 2, ErrMsg: err.Error()})
				return err
			}
			// 更新记录
			txInfo, err := tx.MarshalJSON()
			if err != nil {
				log.Errorf("10. marshal tx err: %v", err)
				return err
			}
			if err := txRecord.Update(models.Transaction{TxHash: tx.Hash().Hex(), State: 1, TxInfo: string(txInfo)}); err != nil {
				log.Errorf("11. 更新transaction数据库error: %v", err)
				return err
			}

			// 监听交易是否成功
			// 开启一个协程来监听交易是否发送成功
			wg.Add(1)
			go func() {
				defer wg.Done()
				count := 0
				for {
					if count > 12 {
						// 如果2分钟都没有查到则状态设置为超时
						if err := txRecord.Update(models.Transaction{State: 4}); err != nil {
							log.Errorf("12. 交易监听超时，交易hash: %s", tx.Hash().Hex())
							return
						}
						return
					}
					count++
					time.Sleep(15 * time.Second) // 休眠15s之后去链上查交易
					log.Warnf("13. 开始监听交易 txHash: %s", tx.Hash().Hex())
					_, isPending, err := client.TransactionByHash(context.Background(), tx.Hash())
					if (err != nil && err == ethereum.NotFound) || isPending {
						// 休眠10s之后再查找
						time.Sleep(10 * time.Second)
						continue
					} else if err != nil && err != ethereum.NotFound {
						// 报错，修改数据库的交易记录为发送失败状态和失败msg
						if err := txRecord.Update(models.Transaction{State: 2, ErrMsg: err.Error()}); err != nil {
							log.Errorf("14. 更新发送交易失败Transaction数据库error: %v", err)
							return
						}
					}

					// err == nil 并且isPending == false则表示交易已经上链
					if err == nil && !isPending {
						log.Infof("15. 监听到交易记录，并更新数据库 txHash: %s", tx.Hash().Hex())
						// 更新对应交易的状态 todo 和下面的表更新优化成事务更新方式
						if err := txRecord.Update(models.Transaction{State: 3}); err != nil {
							log.Errorf("16. 更新发送交易成功状态Transaction数据库error: %v", err)
							return
						}

						// 更新from的ethBalance
						ethBalance, err := client.BalanceAt(context.Background(), common.HexToAddress(txRecord.From), nil)
						if err != nil {
							log.Errorf("17. 获取地址eth余额失败，addr: %s, error: %v", txRecord.From, err)
							return
						}
						// 获取from token余额
						toeknBalance, err := transaction.GetTokenBalance(fromAddr, contractAddr, client)
						if err != nil {
							log.Errorf("18. 获取地址token余额失败，addr: %s, error: %v", fromAddr.Hex(), err)
							return
						}
						// 更新数据库，eth余额和状态
						if err := models.UpdateInfo(models.CollectionAddress{Address: fromAddr.Hex(), State: 3, EthBalance: ethBalance.String(), TokenBalance: toeknBalance.String()}); err != nil {
							log.Errorf("19. 更新CollectionAddress数据库失败，addr: %s, error: %v", fromAddr.Hex(), err)
							return
						}
						return
					}
				}
			}()
		}
	}
	wg.Wait()
	return nil
}
