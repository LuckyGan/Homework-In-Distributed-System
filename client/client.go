package main

import (
    "net/rpc"
    "fmt"
	"os"
	"time"
	"strings"
)

type Args struct {
    UsrPwd string
}

/*
 * 获取服务器端的时间
 */
func getServerTime(client *rpc.Client) int64 {
	args := Args {os.Args[2]}
	var reply int64

	startClientTime := time.Now().UnixNano() //记录客户端启动时间
	err := client.Call("Arith.GetTimeStamp", args, &reply)
	endClientTime := time.Now().UnixNano()   //记录客户端获取响应时间

	if err != nil {
		if err.Error() == "The key is incorrect" {
			fmt.Println(err.Error())
		} else {
			fmt.Println("call error.")
		}
		os.Exit(1)
	}
	return reply + (endClientTime - startClientTime) / 2
}

/*
 * 格式化打印服务器端时间（精确到纳秒）
 */
func printFormatTime(serverTime *int64) {
	seconds := *serverTime / 1e9
	nanoSeconds := *serverTime % 1e9
	tm := time.Unix(seconds, nanoSeconds)
	fmt.Println("Time on " + os.Args[1][0:strings.Index(os.Args[1], ":")] + " is " + tm.Format("2006-01-02 15:04:05.000"))
}

func main()  {
	if len(os.Args) != 3 {
		fmt.Println("args error. ")
		os.Exit(1)
	}
	service := os.Args[1]
    client, err := rpc.Dial("tcp", service)
    if err != nil {
		fmt.Println("connect error.")
		os.Exit(1)
    }

	for {
		serverTime := getServerTime(client)
		printFormatTime(&serverTime)
		time.Sleep(time.Second) //每秒刷新一次
	}
}
