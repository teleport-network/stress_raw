/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	_ "github.com/teleport/stress-raw/logger"
	"github.com/teleport/stress-raw/rpcClient"
	"github.com/teleport/stress-raw/utils"
)

// gwtCmd represents the gwt command
var (
	count  int
	gwtCmd = &cobra.Command{
		Use:   "gwt",
		Short: "generate raw transactions",
		Long: `
通过公里的rpc接口去取关键信息
生成裸交易
与nonce一同写入文件中
	`,
		Run: func(cmd *cobra.Command, args []string) {

			if count < 200 {
				log.Fatal("count 数量不能小于分组的组大小: 200")
			}

			// 基础检查,输出csv的文件路径是否存在
			utils.PathNEAC(csvDir)

			// 创建rpc链接
			rc, err := rpcClient.NewChain("https://evm-rpc2.qa.davionlabs.com")
			if err != nil {
				log.Fatal(err)
			}

			// 获取私钥数据
			pKData, err := utils.ReadCsv(filePath)
			if err != nil {
				log.Fatal(err)
			}

			if startNum+count > len(pKData) {
				if err != nil {
					log.Fatal(fmt.Errorf(
						"文件中有效的privateKey只有%v个, 但期望获取的最大数为%v个",
						len(pKData),
						startNum+count,
					))
				}
			}

			// 启用协程创建交易
			var wg = new(sync.WaitGroup)

			// 创建
			errChan := make(chan error, 100)
			closeChan := make(chan bool)
			wg.Add(1)
			go func(report <-chan error, c <-chan bool) {
				select {
				case err := <-report:
					fmt.Println(err)
					if err != nil {
						log.Fatal(err)
					}
				case _, ok := <-c:
					if !ok {
						wg.Done()
					}
				}
			}(errChan, closeChan)

			// 并发分组, 确认组间串行
			wgg := new(sync.WaitGroup)
			for g := 0; g < int(math.Ceil(float64(count)/200)); g++ {

				wgg.Add(200)
				for i := startNum + g*200; i < startNum+g*200+200; i++ {
					wg.Add(1)
					// 启动协程生成数据
					go func(k int, c chan<- error) {
						// 处理私钥 和 fromAddress
						privKeyStr, fromStr := pKData[k][0][2:], pKData[k][1]

						// 查询nonce
						nonce, err := rc.GetNonce(common.HexToAddress(fromStr))
						if err != nil {
							c <- fmt.Errorf("GetNonce err: %v", err)
						}

						records := [][]string{}
						records = append(records, []string{"nonce", "enRawTxn"})

						for j := 0; j < times; j++ {
							txn, err := rc.MakeTxn(privKeyStr, nonce)
							if err != nil {
								c <- fmt.Errorf("MakeTxn err: %v", err)
							}
							txnStr := rc.EnRawTxn(txn)

							records = append(records, []string{strconv.FormatUint(nonce, 10), txnStr})
							nonce++
						}

						if err := utils.WriteCsv(fmt.Sprintf("%v/%v.csv", csvDir, k), records); err != nil {
							c <- fmt.Errorf("WriteCsv err: %v", err)
						}

						wg.Done()
						wgg.Done()
					}(i, errChan)
				}
				wgg.Wait()
			}

			close(closeChan)
			wg.Wait()
		},
	}
)

func init() {
	gwtCmd.Flags().IntVarP(&count, "count", "c", 1000, "privateKey's count")

	rootCmd.AddCommand(gwtCmd)
}
