package basic_paxos

import (
	"fmt"
)

type Proposer struct {
	id int
	round int
	number int
	acceptors []int
}

func NewProposer(id int, acceptors []int) *Proposer {
	proposer := &Proposer{
		id: id,
		acceptors: acceptors,
	}
	return proposer
}

func (p *Proposer) Propose(value interface{}) interface{} {
	p.round++
	p.number = p.GetNumber()

	prepare_count := 0
	max_number := 0
	
	// phase1
	for _, aid := range p.acceptors {
		req := &Request{
			Number: p.number,
			From: p.id,
			To: aid,
		}
		resp := new(Response)
		ip := fmt.Sprintf("127.0.0.1:%d", aid)
		ok := Call(ip, "Acceptor.Prepare", req, resp)
		if !ok {
			continue
		}
		if resp.Ok {
			prepare_count++
			if (resp.Accepted_number > max_number) {
				max_number = resp.Accepted_number
				value = resp.Accepted_value
			}
		}
		// TODO: 是否到达超过时break循环----> 需要
		if prepare_count >= p.GetMajority() {
			break
		}
	}

	// phase2
	accept_count := 0
	if prepare_count >= p.GetMajority() {
		for _, aid := range p.acceptors {
			req := &Request{
				Number: p.number,
				Value: value,
				From: p.id,
				To: aid,
			}
			resp := new(Response)
			addr := fmt.Sprintf("127.0.0.1:%d", aid)
			ok := Call(addr, "Acceptor.Accept", req, resp)
			if !ok {
				continue
			}
			if resp.Ok {	// TODO: acceptor在Accept阶段返回什么？
				accept_count++
			}
		}
	}

	if accept_count >= p.GetMajority() {
		return value
	} else {
		return nil
	}
}

func (p *Proposer) GetMajority() int {
	return len(p.acceptors) / 2 + 1
}

func (p *Proposer) GetNumber() int {
	return (p.round << 16) | p.id
}