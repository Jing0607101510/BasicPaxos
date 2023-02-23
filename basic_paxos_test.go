package basic_paxos

import (
	"testing"
)

func Start(proposer_ids []int, acceptor_ids []int, learner_ids []int) ([]*Proposer, []*Acceptor, []*Learner) {
	proposers := make([]*Proposer, 0)
	acceptors := make([]*Acceptor, 0)
	learners := make([]*Learner, 0)

	for _, pid := range proposer_ids {
		proposers = append(proposers, NewProposer(pid, acceptor_ids))
	}

	for _, aid := range acceptor_ids {
		acceptors = append(acceptors, newAcceptor(aid, learner_ids))
	}

	for _, lid := range learner_ids {
		learners = append(learners, newLearner(lid, acceptor_ids))
	}

	return proposers, acceptors, learners
}

func CleanUp(acceptors []*Acceptor, learners []*Learner) {
	for _, acceptor := range acceptors {
		acceptor.lis.Close()
	}

	for _, learner := range learners {
		learner.lis.Close()
	}
}

func TestSingle(t *testing.T) {
	proposer_ids := []int{1}
	acceptor_ids := []int{10001, 10002, 10003}
	learner_ids := []int{20001}

	proposers, acceptors, learners := Start(proposer_ids, acceptor_ids, learner_ids)
	defer CleanUp(acceptors, learners)

	expect_value := "hello world"
	value := proposers[0].Propose(expect_value)
	if value != expect_value {
		t.Errorf("value=%s, expect_value=%s", value, expect_value)
	}

	learn_value := learners[0].chosen()
	if value != learn_value {
		t.Errorf("value=%s, learn_value=%s", value, learn_value)
	}
}

func TestTwo(t *testing.T) {
	proposer_ids := []int{1, 2}
	acceptor_ids := []int{10001, 10002, 10003}
	learner_ids := []int{20001}

	proposers, acceptors, learners := Start(proposer_ids, acceptor_ids, learner_ids)
	defer CleanUp(acceptors, learners)

	value1 := proposers[0].Propose("hello world")
	value2 := proposers[1].Propose("you are right")

	if value1 != value2 {
		t.Errorf("value1=%s, value2=%s", value1, value2)
	}

	learn_value := learners[0].chosen()
	if learn_value != value1 && learn_value != value2 {
		t.Errorf("value1=%s, value2=%s, learn_value=%s", value1, value2, learn_value)
	}
}