package main

import (
	"log"
	"net"
	"time"
)

// element in the queue
type Task struct {
	conn net.Conn
}

// represent the thread in the pool
type Worker struct {
	id       int
	taskChan chan Task
}

// represent the thread pool
type Pool struct {
	taskQueue chan Task
	workers   []*Worker
}

// create a new worker
func NewWorker(id int, taskChan chan Task) *Worker {
	return &Worker{
		id:       id,
		taskChan: taskChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for task := range w.taskChan {
			log.Printf("Worker %d is handling job from %s", w.id, task.conn.RemoteAddr())
			handleConnection(task.conn)
		}
	}()
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		taskQueue: make(chan Task),
		workers:   make([]*Worker, numOfWorker),
	}
}

// push task to queue
func (p *Pool) AddTask(conn net.Conn) {
	p.taskQueue <- Task{conn: conn}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.taskQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1000)
	conn.Read(buf)
	time.Sleep(1 * time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nEngineer Pro\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// 1 pool with 2 threads
	pool := NewPool(2)
	pool.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		pool.AddTask(conn)
	}
}
