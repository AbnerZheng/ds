package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import (
	"labrpc"
	"bytes"
	"encoding/gob"
	"fmt"
	//"time"
	//"math/rand"
	"strconv"
)

// import "bytes"
// import "encoding/gob"

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

type Log struct {
	Command string
	Term    int
}

type State int8

const (
	FOLLOW State = iota + 1
	LEADER
	CANDIDATE
)

//
// A Go object implementing a singl(
// e Raft peer.
//
type Raft struct {
	mu        sync.Mutex
	peers     []*labrpc.ClientEnd
	persister *Persister
	me        int // index into peers[]

	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// persistent state
	currentTerm int // latest term server has seen (initialized to 0 on first boot, increase monotonically) 单调不减
	voteFor     int // candidateId that received vote in current term (or null if none)
	// log entries;
	// each entry contains command for state machine,
	// and term when entry was received by leader(first index is 1)
	log []Log

	state State // 状态机标识位

	// volatile state on all server
	commitIndex int // index of highest log entry known to be commited(initialized to 0, increase monotonically)
	lastApplied int // index of highest log entry applied to state machine(initialized to 0, increase monotonically)

	// volatile state on leader
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	term = rf.currentTerm
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here.
	// Example:
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.voteFor)
	e.Encode(rf.log)
	data := w.Bytes()
	rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.log)
	d.Decode(&rf.voteFor)
	d.Decode(&rf.currentTerm)
}

//
// example RequestVote RPC arguments structure.
//
type RequestVoteArgs struct {
	Term         int // candidate's term
	CandidateId  int // candidate requesting vote
	LastLogIndex int // index of  candidate's last log entry
	LastLogTerm  int // term of candidate's last log entry
}

//
// example RequestVote RPC reply structure.
//
type RequestVoteReply struct {
	Term        int  // current term, for candidate to update itself
	VoteGranted bool // true mean candidate received vote
}

//
// RequestVote RPC handler.
// 在无leader的情况下， 各个candidate向其他node发送该RPC以进行竞选leader
//
func (rf *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) {
	reply.Term = rf.currentTerm
	reply.VoteGranted = true
	//if args.Term < rf.currentTerm { // 如果自己的term已经被过来拉票的server还高的话，直接否决
	//	reply.VoteGranted = false
	//} else {
	//	if rf.voteFor == -1 || rf.voteFor == args.CandidateId { // go 中没有null，所以以0来表示
	//		if len(rf.log) == 0 {
	//			reply.VoteGranted = true
	//		}else {
	//			lastLog := rf.log[len(rf.log)-1]
	//			if lastLog.Term < args.LastLogTerm {
	//				reply.VoteGranted = true
	//			} else if lastLog.Term == args.LastLogTerm && len(rf.log) < args.LastLogIndex {
	//				reply.VoteGranted = true
	//			} else {
	//				reply.VoteGranted = false
	//			}
	//		}
	//	}
	//}

	//if reply.VoteGranted {
	//	rf.mu.Lock()
	//	rf.voteFor = args.CandidateId
	//	rf.mu.Unlock()
	//}
}

//
// example code to send a RequestVot
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// returns true if labrpc says the RPC was delivered.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

func (rf *Raft) String() string {
	return fmt.Sprintf("me: %d,Term: %d,length of log:%d", rf.me, rf.currentTerm, len(rf.log))
}

func (rf *Raft)stateMachine() {
	count := 0
	c1 := make(chan bool)
	for true{
		switch rf.state {
		case FOLLOW:
			// 如果一定时间内没收到heart beat，则转换为candidate
			rf.mu.Lock()
			rf.state = CANDIDATE
			rf.mu.Unlock()
		case CANDIDATE:
			go func() {
				//randDuration := time.Duration(150 + rand.Intn(150))
				//time.Sleep(randDuration * time.Millisecond)
				c1 <- true
			}()
			// 等待一段时间后，开始竞选
			_ = <-c1
			rf.mu.Lock()
			rf.currentTerm ++
			rf.voteFor = rf.me
			rf.mu.Unlock()

			args := RequestVoteArgs{rf.currentTerm, rf.me, 0, 0}
			if len(rf.log) > 0 {
				lastLog := rf.log[len(rf.log)-1]
				args.LastLogTerm = lastLog.Term
				args.LastLogIndex = len(rf.log)
			}
			for i := range rf.peers {
				if rf.me != i {
					go func(index int) {
						res := &RequestVoteReply{}
						r := rf.sendRequestVote(index, args, res)
						fmt.Printf("RPC result: %t\n", r)
						if res.VoteGranted{
							count += 1
							if count*2 >= len(rf.peers) {
								fmt.Println("I'm the leader")
							}
							fmt.Println("count " + strconv.Itoa(count))
						}
						fmt.Println(strconv.Itoa(rf.me) + "->" + strconv.Itoa(index))
						fmt.Println(args)
						fmt.Println(res)

						fmt.Println(rf)
					}(i)
				}
			}

		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[].
// 在peers中存有自己, 自己可以用peers[me]取得。
// this server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any.
// applyCh 是发送ApplyMsg的channel
// applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
// 长期的任务应该放在goroutines中执行
//

func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	// Your initialization code here.
	rf.state = FOLLOW

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	go rf.stateMachine()
	return rf
}
