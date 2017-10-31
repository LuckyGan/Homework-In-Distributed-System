package main

import (
    "net/rpc"
	"time"
	"fmt"
	"os"
	"net"
	"errors"
)

type Args struct {
    UsrPwd string
}

type Arith int

/*
 * 格式化打印服务器端时间（精确到毫秒）
 */
func printFormatTime(serverTime *int64) {
	seconds := *serverTime / 1e9
	nanoSeconds := *serverTime % 1e9
	tm := time.Unix(seconds, nanoSeconds)
	fmt.Println(tm.Format("2006-01-02 15:04:05.000"))
}

/*
 * 获取本地服务器端时间的RPC函数
 */
func (t *Arith) GetTimeStamp(args *Args, reply *int64) error {
	if args.UsrPwd == "123" {
		*reply = time.Now().UnixNano()
		printFormatTime(reply)
		return nil
	} else {
		return errors.New("The key is incorrect")
	}
}

/*
 * 统一检测错误方法
 */
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func main() {
    arith := new(Arith)
    rpc.Register(arith)
    rpc.HandleHTTP()

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
        conn, err := listener.Accept()
        if err != nil {
			continue
        }
		go rpc.ServeConn(conn)
    }
}
