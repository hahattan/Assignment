package node

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	Debug        = 0
	TimeMultiple = 10
)

type Node struct {
	ID                 int
	mu                 sync.Mutex
	peers              []*Node
	currentTerm        int
	votedFor           int
	state              State
	electionResetEvent time.Time
}

func NewNode(id int, peers []*Node, ready <-chan any) *Node {
	node := &Node{
		ID:       id,
		peers:    peers,
		state:    Follower,
		votedFor: -1,
	}

	go func() {
		<-ready
		fmt.Printf("Member %d: Hi\n", id)
		node.mu.Lock()
		node.electionResetEvent = time.Now()
		node.mu.Unlock()
		node.runElectionTimer()
	}()

	return node
}

func (n *Node) Stop() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = Dead
	n.dlog("becomes dead")
}

// dlog logs a debugging message if Debug > 0.
func (n *Node) dlog(format string, args ...interface{}) {
	if Debug > 0 {
		format = fmt.Sprintf("[%d] ", n.ID) + format
		log.Printf(format, args...)
	}
}

// electionTimeout generates a pseudo-random election timeout duration.
func (n *Node) electionTimeout() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond * TimeMultiple
}

// runElectionTimer implements an election timer. It should be launched whenever
// we want to start a timer towards becoming a candidate in a new election.
//
// This function is blocking and should be launched in a separate goroutine;
// it's designed to work for a single (one-shot) election timer, as it exits
// whenever the Node state changes from follower/candidate or the term changes.
func (n *Node) runElectionTimer() {
	timeoutDuration := n.electionTimeout()
	n.mu.Lock()
	termStarted := n.currentTerm
	n.mu.Unlock()
	n.dlog("election timer started (%v), term=%d", timeoutDuration, termStarted)

	// This loops until either:
	// - we discover the election timer is no longer needed, or
	// - the election timer expires and this Node becomes a candidate
	// In a follower, this typically keeps running in the background for the
	// duration of the Node's lifetime.
	ticker := time.NewTicker(10 * time.Millisecond * TimeMultiple)
	defer ticker.Stop()
	for {
		<-ticker.C

		n.mu.Lock()
		if n.state != Candidate && n.state != Follower {
			n.dlog("[election timer] state=%s, bailing out", n.state)
			n.mu.Unlock()
			return
		}

		if termStarted != n.currentTerm {
			n.dlog("[election timer] term changed from %d to %d, bailing out", termStarted, n.currentTerm)
			n.mu.Unlock()
			return
		}

		// Start an election if we haven't heard from a leader or haven't voted for
		// someone for the duration of the timeout.
		if elapsed := time.Since(n.electionResetEvent); elapsed >= timeoutDuration {
			n.startElection()
			n.mu.Unlock()
			return
		}
		n.mu.Unlock()
	}
}

// startElection starts a new election with this Node as a candidate.
// Expects n.mu to be locked.
func (n *Node) startElection() {
	fmt.Printf("Member %d: I want to be leader\n", n.ID)
	n.state = Candidate
	n.currentTerm += 1
	savedCurrentTerm := n.currentTerm
	n.electionResetEvent = time.Now()
	n.votedFor = n.ID
	n.dlog("starting election (currentTerm=%d)", savedCurrentTerm)

	votesReceived := 0

	// Send RequestVote RPCs to all other servers concurrently.
	for _, peer := range n.peers {
		go func(peer *Node) {
			if peer.ID == n.ID {
				return
			}
			args := RequestVoteRequest{
				Term:        savedCurrentTerm,
				CandidateId: n.ID,
			}
			var reply RequestVoteResponse

			n.dlog("sending RequestVote to %d: %+v", peer.ID, args)
			if err := peer.RequestVote(args, &reply); err == nil {
				n.mu.Lock()
				defer n.mu.Unlock()
				n.dlog("received RequestVoteResponse %+v", reply)

				if n.state != Candidate {
					n.dlog("while waiting for reply, state = %v", n.state)
					return
				}

				if reply.Term > savedCurrentTerm {
					n.dlog("term out of date in RequestVoteResponse")
					n.becomeFollower(reply.Term)
					return
				} else if reply.Term == savedCurrentTerm {
					if reply.VoteGranted {
						votesReceived += 1
						if votesReceived > len(n.peers)/2 {
							// Won the election!
							fmt.Printf("Member %d voted to be leader: (%d > %d/2)\n", n.ID, votesReceived, len(n.peers))
							n.dlog("wins election with %d votes", votesReceived)
							n.startLeader()
							return
						}
					}
				}
			}
		}(peer)
	}

	// Run another election timer, in case this election is not successful.
	go n.runElectionTimer()
}

type RequestVoteRequest struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteResponse struct {
	Term        int
	VoteGranted bool
}

func (n *Node) RequestVote(req RequestVoteRequest, resp *RequestVoteResponse) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state == Dead {
		return nil
	}
	n.dlog("RequestVote: %+v [currentTerm=%d, votedFor=%d]", req, n.currentTerm, n.votedFor)

	if req.Term > n.currentTerm {
		n.dlog("... term out of date in RequestVote")
		n.becomeFollower(req.Term)
	}

	if n.currentTerm == req.Term &&
		(n.votedFor == -1 || n.votedFor == req.CandidateId) {
		fmt.Printf("Member %d: Accept member %d to be leader\n", n.ID, req.CandidateId)
		resp.VoteGranted = true
		n.votedFor = req.CandidateId
		n.electionResetEvent = time.Now()
	} else {
		resp.VoteGranted = false
	}
	resp.Term = n.currentTerm
	n.dlog("... RequestVote resp: %+v", resp)
	return nil
}

// startLeader switches n into a leader state and begins process of heartbeats.
// Expects n.mu to be locked.
func (n *Node) startLeader() {
	n.state = Leader
	n.dlog("becomes Leader; term=%d, log=%v", n.currentTerm, "idk")

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond * TimeMultiple)
		defer ticker.Stop()

		// Send periodic heartbeats, as long as still leader.
		for {
			n.leaderSendHeartbeats()
			<-ticker.C

			n.mu.Lock()
			if n.state != Leader {
				n.mu.Unlock()
				return
			}
			n.mu.Unlock()
		}
	}()
}

// leaderSendHeartbeats sends a round of heartbeats to all peers, collects their
// replies and adjusts n's state.
func (n *Node) leaderSendHeartbeats() {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return
	}
	savedCurrentTerm := n.currentTerm
	n.mu.Unlock()

	for _, peer := range n.peers {
		if peer.ID == n.ID {
			return
		}
		args := AnswerHeartbeatRequest{
			Term:     savedCurrentTerm,
			LeaderId: n.ID,
		}
		go func(peer *Node) {
			n.dlog("sending AnswerHeartBeat to %v: ni=%d, args=%+v", peer.ID, 0, args)
			var reply AnswerHeartbeatResponse
			if err := peer.AnswerHeartBeat(args, &reply); err == nil {
				n.mu.Lock()
				defer n.mu.Unlock()
				if reply.Term > savedCurrentTerm {
					n.dlog("term out of date in heartbeat reply")
					n.becomeFollower(reply.Term)
					return
				}
			}
		}(peer)
	}
}

type AnswerHeartbeatRequest struct {
	Term     int
	LeaderId int

	PrevLogIndex int
	PrevLogTerm  int
	LeaderCommit int
}

type AnswerHeartbeatResponse struct {
	Term    int
	Success bool
}

func (n *Node) AnswerHeartBeat(req AnswerHeartbeatRequest, resp *AnswerHeartbeatResponse) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state == Dead {
		return nil
	}
	n.dlog("AnswerHeartBeat: %+v", req)

	if req.Term > n.currentTerm {
		n.dlog("... term out of date in AnswerHeartBeat")
		n.becomeFollower(req.Term)
	}

	resp.Success = false
	if req.Term == n.currentTerm {
		if n.state != Follower {
			n.becomeFollower(req.Term)
		}
		n.electionResetEvent = time.Now()
		resp.Success = true
	}

	resp.Term = n.currentTerm
	n.dlog("AnswerHeartBeat resp: %+v", *resp)
	return nil
}

// becomeFollower makes n a follower and resets its state.
// Expects n.mu to be locked.
func (n *Node) becomeFollower(term int) {
	n.dlog("becomes follower with term=%d", term)
	n.state = Follower
	n.currentTerm = term
	n.votedFor = -1
	n.electionResetEvent = time.Now()

	go n.runElectionTimer()
}
