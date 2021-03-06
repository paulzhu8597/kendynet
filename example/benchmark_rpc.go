package main

import(
	"time"
	"sync/atomic"
	"strconv"
	"fmt"
	"os"
	"github.com/sniperHW/kendynet/example/testproto"
	"github.com/sniperHW/kendynet/example/test_rpc"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/kendynet/rpc"
	"github.com/sniperHW/kendynet"
	"github.com/sniperHW/kendynet/golog"
)

func server(service string) {
	count := int32(0)
	go func() {
		for {
			time.Sleep(time.Second)
			tmp := atomic.LoadInt32(&count)
			atomic.StoreInt32(&count,0)
			fmt.Printf("count:%d\n",tmp)			
		}
	}()

	server := test_rpc.NewRPCServer()
	//注册服务
	server.RegisterMethod("hello",func (replyer *rpc.RPCReplyer,arg interface{}){
		atomic.AddInt32(&count,1)
		world := &testproto.World{World:proto.String("world")}
		replyer.Reply(world,nil)
	})
	server.Serve(service)
}

func client(service string,count int) {
	hello := &testproto.Hello{Hello:proto.String("hello")}
	caller := test_rpc.NewCaller("hello")
	var onResp func(ret interface{},err error)
	onResp = func(ret interface{},err error){
		if nil != ret {
			err := caller.Call(hello,onResp)
			if err != nil {
				fmt.Printf("%s\n",err.Error())
				return
			}
		} else if nil != err {
			fmt.Printf("%s\n",err.Error())
		}
	}

	for i := 0; i < count ; i++ {
		err := caller.Dial(service,10 * time.Second)
		if err != nil {
			fmt.Printf("%s\n",err.Error())
		} else {
			caller.Call(hello,onResp)
		}
	}
}

func main(){

	outLogger := golog.NewOutputLogger("log","kendynet",1024*1024*1000)
	kendynet.InitLogger(outLogger)
	rpc.InitLogger(outLogger)

	if len(os.Args) < 3 {
		fmt.Printf("usage ./pingpong [server|client|both] ip:port clientcount\n")
		return
	}


	mode := os.Args[1]

	if !(mode == "server" || mode == "client" || mode == "both") {
		fmt.Printf("usage ./pingpong [server|client|both] ip:port clientcount\n")
		return
	}

	service := os.Args[2]

	sigStop := make(chan bool)

	if mode == "server" || mode == "both" {
		go server(service)
	}

	if mode == "client" || mode == "both" {
		if len(os.Args) < 4 {
			fmt.Printf("usage ./pingpong [server|client|both] ip:port clientcount\n")
			return
		}
		connectioncount,err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		//让服务器先运行
		time.Sleep(time.Second)
		go client(service,connectioncount)

	}

	_,_ = <- sigStop

	return

}


