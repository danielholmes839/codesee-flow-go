package queue_test

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/network/gossip/libp2p/queue"
)

// TestRetrievalByPriority tests that message can be retrieved in priority order
func TestRetrievalByPriority(t *testing.T) {
	// create a map of messages -> priority with messages assigned random priorities
	messages := createMessages(1000, randomPriority)
	testQueue(t, messages)
}

// TestRetrievalByInsertionOrder tests that messages with the same priority can be retrieved in insertion order
func TestRetrievalByInsertionOrder(t *testing.T) {

	// create a map of messages -> priority with messages assigned fixed priorities
	messages := createMessages(1000, fixedPriority)
	testQueue(t, messages)
}

// TestConcurrentQueueAccess tests that the queue can be safely accessed concurrently
func TestConcurrentQueueAccess(t *testing.T) {
	writerCnt := 5
	readerCnt := 5
	messageCnt := 1000

	messages := createMessages(messageCnt, randomPriority)

	var priorityFunc queue.MessagePriorityFunc = func(message interface{}) queue.Priority {
		return messages[message.(string)]
	}

	msgChan := make(chan string, len(messages))
	for k := range messages {
		msgChan <- k
	}
	close(msgChan)

	ctx, cancel := context.WithCancel(context.Background())
	mq := queue.NewMessageQueue(ctx, priorityFunc)

	writeWg := sync.WaitGroup{}
	write := func() {
		defer writeWg.Done()
		for msg := range msgChan {
			err := mq.Insert(msg)
			assert.NoError(t, err)
		}
	}
	var readMsgCnt int64
	done := make(chan struct{})
	defer close(done)
	read := func() {
		for {
			select {
			case <-done:
				return
			default:
				mq.Remove()
				atomic.AddInt64(&readMsgCnt, 1)
			}
		}
	}

	// kick off writers
	for i := 0; i < writerCnt; i++ {
		writeWg.Add(1)
		go write()
	}

	// kick off readers
	for i := 0; i < readerCnt; i++ {
		go read()
	}

	writeWg.Wait()

	assert.Eventually(t, func() bool {
		actualCnt := atomic.LoadInt64(&readMsgCnt)
		return int64(messageCnt) == actualCnt
	}, 5*time.Second, 5*time.Millisecond)

	cancel()

	assert.Equal(t, 0, mq.Len())
}

// TestQueueShutdown tests that Remove unblocks when the context is shutdown
func TestQueueShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mq := queue.NewMessageQueue(ctx, fixedPriority)
	ch := make(chan struct{})

	go func() {
		mq.Remove()
		close(ch)
	}()

	cancel()
	assert.Eventually(t, func() bool {
		select {
		case <-ch:
			return true
		default:
			return false
		}
	}, time.Second, time.Millisecond)
}

func testQueue(t *testing.T, messages map[string]queue.Priority) {

	// create the priority function
	var priorityFunc queue.MessagePriorityFunc = func(message interface{}) queue.Priority {
		return messages[message.(string)]
	}

	// create queues for each priority to check expectations later
	queues := make(map[queue.Priority][]string)
	for p := queue.LowPriority; p <= queue.HighPriority; p++ {
		queues[queue.Priority(p)] = make([]string, 0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create the queue
	mq := queue.NewMessageQueue(ctx, priorityFunc)

	// insert all elements in the queue
	for msg, p := range messages {

		err := mq.Insert(msg)
		assert.NoError(t, err)

		// remember insertion order to check later
		queues[p] = append(queues[p], msg)

		// sleep to make sure two consecutive elements don't get the same timestamp
		time.Sleep(1 * time.Millisecond)
	}

	// create a slice of the expected messages in the order in which they are expected
	var expectedMessages []string
	for p := queue.HighPriority; p >= queue.LowPriority; p-- {
		expectedMessages = append(expectedMessages, queues[queue.Priority(p)]...)
	}

	// check queue length
	assert.Equal(t, len(expectedMessages), mq.Len())

	// check that elements are retrieved in order
	for i := 0; i < len(expectedMessages); i++ {

		item := mq.Remove()

		assert.Equal(t, expectedMessages[i], item.(string))
	}

	assert.Equal(t, 0, mq.Len())
}

func BenchmarkPush(b *testing.B) {
	b.StopTimer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var mq = queue.NewMessageQueue(ctx, randomPriority)
	for i := 0; i < b.N; i++ {
		err := mq.Insert("test")
		if err != nil {
			b.Error(err)
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := mq.Insert("test")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkPop(b *testing.B) {
	b.StopTimer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var mq = queue.NewMessageQueue(ctx, randomPriority)
	for i := 0; i < b.N; i++ {
		err := mq.Insert("test")
		if err != nil {
			b.Error(err)
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mq.Remove()
	}
}

func createMessages(messageCnt int, priorityFunc queue.MessagePriorityFunc) map[string]queue.Priority {
	msgPrefix := "message"
	// create a map of messages -> priority
	messages := make(map[string]queue.Priority, messageCnt)

	for i := 0; i < messageCnt; i++ {
		// choose a random priority
		p := priorityFunc(nil)
		// create a message
		msg := msgPrefix + strconv.Itoa(i)
		messages[msg] = p
	}

	return messages
}

func randomPriority(_ interface{}) queue.Priority {
	rand.Seed(time.Now().UnixNano())
	p := rand.Intn(int(queue.HighPriority-queue.LowPriority+1)) + int(queue.LowPriority)
	return queue.Priority(p)
}

func fixedPriority(_ interface{}) queue.Priority {
	return queue.MediumPriority
}
