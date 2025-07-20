package main

import(
	"sync"
)

//实现一个简单的异步读写安全队列

type SyncQueue[T any] struct{
	lock	sync.RWMutex
	queue	[]T
}

func (queue *SyncQueue[T] )push(ele T){
	queue.hold()
	defer queue.release()
	queue.queue=append(queue.queue,ele)
}

func (queue *SyncQueue[T])pop()T{
	queue.hold()
	defer queue.release()
	ret:=queue.queue[0]
	queue.queue=queue.queue[1:]
	return ret
}
//unSafe
func (queue *SyncQueue[T])clear(){
	queue.queue=make([]T,0)
}
//
func (queue *SyncQueue[T])hold(){
	queue.lock.Lock()
}

func (queue *SyncQueue[T])release(){
	queue.lock.Unlock()
}

