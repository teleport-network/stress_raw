/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
	_ "github.com/teleport/stress-raw/logger"
	"github.com/teleport/stress-raw/rpcClient"
	"github.com/teleport/stress-raw/utils"
)

// runCmd represents the run command
var (
	userNum int
	rpcList = []string{
		"https://evm-rpc2.qa.davionlabs.com",
		"https://evm-rpc2.qa.davionlabs.com",
		"https://evm-rpc2.qa.davionlabs.com",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {

			wg := new(sync.WaitGroup)
			// 初始化 rpcUrl slice
			for i, r := range rpcList {
				// 创建rpc链接
				rc, err := rpcClient.NewChain(r)
				if err != nil {
					log.Fatal(err)
				}

				// 遍历需要的文件
				for j := startNum + i*userNum; j < startNum+(i+1)*userNum; j++ {

					wg.Add(1)
					// 获取序列化后的数据
					filePath := fmt.Sprintf("%v/%v.csv", csvDir, j)

					txnData, err := utils.ReadCsvByLine(filePath, times+1)
					if err != nil {
						log.Fatal(err)
					}

					go func(t [][]string, rc *rpcClient.Chain, f string) {

						// 循环请求
						for _, txn := range t {
							txHash, nonce, err := rc.SendRawTxn(txn[1])
							if err != nil {
								if strings.Contains(fmt.Sprintf("%v", err), "tx already in mempool") {
									log.Println(err)
								} else {
									log.Fatal(err)
								}
							} else {
								fmt.Printf("txHash: %v,nonce: %v, filePath:%v\n", txHash, nonce, f)
							}
						}
						wg.Done()
					}(txnData[1:], rc, filePath)
				}
			}
			wg.Wait()
		},
	}
)

func init() {
	runCmd.Flags().IntVarP(&userNum, "users", "u", 100, "the number of simultaneously send txn users")

	rootCmd.AddCommand(runCmd)
}
