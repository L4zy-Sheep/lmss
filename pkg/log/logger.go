package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	NOCOLOR = "\033[0m[INFO] %s\033[0m"
	RED     = "\033[31m[ERROR] %s\033[0m"
	GREEN   = "\033[32m[SUCCESS] %s\033[0m"
	YELLOW  = "\033[33m[WARN] %s\033[0m"
)

var (
	out io.Writer = os.Stdout
)

//没写完

func Error(s string) {
	out.Write([]byte(fmt.Sprintf(RED, time.Now().Format(time.DateTime)+"\t"+s+"\n")))
}

func Warn(s string) {
	out.Write([]byte(fmt.Sprintf(YELLOW, time.Now().Format(time.DateTime)+"\t"+s+"\n")))
}

func Info(s string) {
	out.Write([]byte(fmt.Sprintf(NOCOLOR, time.Now().Format(time.DateTime)+"\t"+s+"\n")))
}

func Success(s string) {
	out.Write([]byte(fmt.Sprintf(GREEN, time.Now().Format(time.DateTime)+"\t"+s+"\n")))
}
