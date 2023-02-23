package basic_paxos

import (
	"fmt"
	"net"
	"net/rpc"
	"log"
)

type Acceptor struct {
	lis net.Listener
	id int
	min_proposal int
	accepted_number int
	accepted_value interface{}
	learners []int
}

func newAcceptor(id int, learners []int) *Acceptor {
	acceptor := &Acceptor{
		id: id,
		learners: learners,
	}
	acceptor.server()
	return acceptor
}

func (a *Acceptor) Prepare(req *Request, resp *Response) error {
	if req.Number <= a.min_proposal { // TODO: 包不包括等号----->不包含
		resp.Ok = false
	} else {
		a.min_proposal = req.Number
		resp.Ok = true
		resp.Accepted_number = a.accepted_number
		resp.Accepted_value = a.accepted_value
	}
	return nil
}

func (a *Acceptor) Accept(req *Request, resp *Response) error {
	if req.Number < a.min_proposal { // TODO: 同上
		resp.Ok = false
	} else {
		resp.Ok = true
		a.min_proposal = req.Number
		a.accepted_number = req.Number
		a.accepted_value = req.Value
		for _, lid := range a.learners {
			learn_req := &Request{
				From: a.id,
				To: lid,
				Number: req.Number,
				Value: req.Value,
			}
			learn_resp := new(Response)
			addr := fmt.Sprintf("127.0.0.1:%d", lid)
			ok := Call(addr, "Learner.Learn", learn_req, learn_resp)
			if !ok {
				continue
			}
		}
	}
	return nil
}

func (a *Acceptor) close() {
	a.lis.Close()
}

func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	addr := fmt.Sprintf(":%d", a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		fmt.Printf("Acceptor listen %s fail", addr)
		log.Fatal("listen error: ", e)
	}
	a.lis = l

	go func() {
		for {
			conn, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}