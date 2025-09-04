package cmd

import (
	"lmss/pkg/utils"
	"log"

	"github.com/spf13/cobra"
)

var (
	hosts   chan string
	ips     []string
	thread  int
	tp      *utils.ThreadPool
	timeout uint8
)

var RootCmd = &cobra.Command{
	Use: "lmss",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) { //当前命令以及所有子命令执行前都会执行
		tp = utils.NewPool(thread)
		hosts = make(chan string, thread)
		go utils.ParseCIDR(ips, hosts)
	},
	PreRun: func(cmd *cobra.Command, args []string) { //当前命令执行前才会执行
		log.Println("cmd.PreRun")
	},
}

func init() {
	RootCmd.PersistentFlags().StringSliceVarP(&ips, "ips", "H", nil, "-H 192.168.1.1/24,192.168.1.2")
	RootCmd.PersistentFlags().IntVarP(&thread, "thread", "T", 5, "-T 5	--default 5 threads")
	RootCmd.PersistentFlags().Uint8VarP(&timeout, "timeout", "t", 100, "-t 100	--default 100ms")
}
