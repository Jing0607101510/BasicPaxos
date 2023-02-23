package basic_paxos

import (
	"net"
	"net/rpc"
	"fmt"
	"log"
)

type Learner struct {
	lis net.Listener
	id int
	msg map[int]Request
}

func newLearner(id int, acceptors []int) *Learner {
	learner := &Learner{
		id: id,
		msg: make(map[int]Request),
	}

	for _, acceptor := range acceptors {
		learner.msg[acceptor] = Request{
			Number: 0,
			Value: nil,
		}
	}

	learner.server()
	return learner
}

func (l *Learner) Learn(req *Request, resp *Response) error {
	if req.Number > l.msg[req.From].Number { // TODO: 这里的条件是什么？
		resp.Ok = true
		l.msg[req.From] = *req
	} else {
		resp.Ok = false
	}
	return nil
}

func (l *Learner) chosen() interface{} {
	// 统计各个number的数量，以及各个number对应的Response
	nubmer_count := make(map[int]int)
	number_to_response := make(map[int]Request)

	for _, accepted := range l.msg {
		if accepted.Number != 0 {
			nubmer_count[accepted.Number]++
			number_to_response[accepted.Number] = accepted
		}
	}
	
	for number, count := range nubmer_count {
		if count >= l.getMajority() {
			return number_to_response[number].Value
		}
	}

	return nil
}

func (l *Learner) close() {
	l.lis.Close()
}

func (l *Learner) getMajority() int {
	return len(l.msg)/2 + 1
}

func (l *Learner) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(l)
	addr := fmt.Sprintf(":%d", l.id)
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		fmt.Printf("Learner listen %s fail", addr)
		log.Fatal("listen error: ", e)
	}
	l.lis = lis

	go func() {
		for {
			conn, err := l.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}