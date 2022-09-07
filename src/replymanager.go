package src

type ReplyManager struct {
	numProcesses int
	numReplies   int
	replyQueue   []int
}

func NewReplyManager(numProcesses int) *ReplyManager {
	return &ReplyManager{
		numProcesses: numProcesses,
		numReplies:   0,
		replyQueue:   make([]int, 0),
	}
}

func (m *ReplyManager) ReceiveReply() bool {
	m.numReplies++
	if m.numReplies == m.numProcesses-1 {
		return true
	}
	return false
}

func (m *ReplyManager) EnqueueProcess(pid int) {
	m.replyQueue = append(m.replyQueue, pid)
}

func (m *ReplyManager) Dequeue() []int {
	return m.replyQueue
}

func (m *ReplyManager) Reset() {
	m.numReplies = 0
	m.replyQueue = make([]int, 0)
}
