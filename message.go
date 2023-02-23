package basic_paxos

import (
	"net/rpc"
	"fmt"
)

type Request struct {
	Number int
	Value interface{}
	From int
	To int
}

type Response struct {
	Ok bool
	Accepted_number int
	Accepted_value interface{}
}

func Call(ip string, func_name string, req interface{}, resp interface{}) bool { // TODO: 参数是否需要*; 返回值
	client, err := rpc.Dial("tcp", ip)
	if err != nil {
		fmt.Printf("Dial ip %s func_name %s fail\n", ip, func_name)
		return false
	}

	defer client.Close()

	err = client.Call(func_name, req, resp)
	if err != nil {
		fmt.Printf("Call ip %s func_name %s fail, err %v\n", ip, func_name, err)
		return false
	}

	return true
}

