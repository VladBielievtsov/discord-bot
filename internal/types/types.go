package types

import "sync"

type QueueItem struct {
	ID       string
	VideoURL string
	User     QueueUser
}

type QueueUser struct {
	Name      string
	AvatarURL string
}

type Queue struct {
	Items []QueueItem
	mu    sync.Mutex
}

func (q *Queue) Enqueue(item QueueItem) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Items = append(q.Items, item)
}

func (q *Queue) Dequeue() *QueueItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return &item
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.Items) == 0
}
