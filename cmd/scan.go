package cmd

import (
	"fmt"
	"lmss/pkg/log"
	"lmss/pkg/utils"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	method      string
	aliveHost   chan string
	ports       []string
	targetPorts []int
)
var supportMethod = []string{"icmp"}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan if target host or targetPorts alive",
	Long: "You can use scan command to scan if ips alive or if ports have opened and what service is\n" +
		"Scan host alive e.g. lmss scan -H 192.168.1.1/24,192.168.1.2 " +
		"Scan targetPorts e.g. lmss scan -H 192.168.1.1/24,192.168.1.2 -P 22,80,3306",
	Run: func(cmd *cobra.Command, args []string) {
		if ports != nil {
			portScan()
		} else {
			hostScan()
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		//method = strings.ToLower(method)//
		method = "icmp" //wait for more
		for _, v := range ports {
			if matched, _ := regexp.MatchString("^[0-9]{1,5}-[0-9]{1,5}$", v); matched {
				tmp := strings.Split(v, "-")
				start, err := strconv.Atoi(tmp[0])
				if err != nil {
					log.Error(fmt.Sprintf("parse port %s err e.g. 22,80,1-65535\n%v", v, err))
					continue
				}
				end, err := strconv.Atoi(tmp[1])
				if err != nil {
					log.Error(fmt.Sprintf("parse port %s err e.g. 22,80,1-65535\n%v", v, err))
					continue
				}
				if start > end {
					start, end = end, start
				}
				if end > 0xffff || end < 0x1 {
					log.Error("target port out of range")
					continue
				}
				for start <= end {
					targetPorts = append(targetPorts, start)
				}
			} else if matched, _ := regexp.MatchString("^[0-9]{1,5}$", v); matched {
				tmp, err := strconv.Atoi(v)
				if err != nil || tmp < 1 || tmp > 0xffff {
					log.Error(fmt.Sprintf("parse port %s err e.g. 22,80,1-65535\n%v", v, err))
					continue
				}
				targetPorts = append(targetPorts, tmp)
			} else {
				continue
			}
		}
	},
}

func hostScan() {
	for host := range hosts {
		switch method {
		case "icmp":
			tp.Start(func() {
				utils.Icmp(host, timeout)
			})
		//case "arp":
		//	if utils.Arp(host) {
		//		log.Success(host)
		//	}
		default:
			log.Error("Unsupported method! Now it just supported icmp")
		}
	}
	tp.Stop()
}

func portScan() {
	for host := range hosts {
		for _, port := range targetPorts {
			tp.Start(func() {
				var res []byte
				var service string
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					goto end
				}
				defer conn.Close()
				conn.Write([]byte("let me see see"))
				conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
				_, err = conn.Read(res)
				if err != nil {
					goto end
				}
				service, err = utils.IdentifyService(res)
				if err != nil {
					log.Error(err.Error())
				} else {
					log.Success(service)
				}
			end:
			})
		}
	}
}

func init() {
	scanCmd.Flags().StringSliceVarP(&ports, "ports", "P", nil, "-P 22,80,3306")
	scanCmd.Flags().StringVarP(&method, "method", "m", "icmp", "-m icmp")
	RootCmd.AddCommand(scanCmd)

}
