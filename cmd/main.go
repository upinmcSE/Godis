package main

import (
	"log"
	"net"
	"time"
)

// element in the queue
type Job struct {
	conn net.Conn
}

// represent the thread in the pool
type Worker struct {
	id int
	jobChan chan Job
}

// represent the thread pool
type Pool struct {
	jobQueue chan Job
	workers []*Worker
}


// Create a new worker
func NewWorker (id int, jobChan chan Job) *Worker{
	return &Worker{
		id: id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start(){
	go func(){
		for job := range w.jobChan {
			log.Printf("Worker %d is handling job from %s", w.id, job.conn.RemoteAddr())
			process(job.conn)
		}
	}()
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers: make([]*Worker, numOfWorker),
	}
}

func (p *Pool) Start(){
	for i := 0; i< len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

// push job to queue
func (p *Pool) AddJob(conn net.Conn){
	p.jobQueue <- Job{conn: conn}
}

func process(conn net.Conn) {
	defer conn.Close()
	var buf []byte = make([]byte, 1000)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep((time.Second * 1))
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, World\r\n"))
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
		// conn == socket
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		pool.AddJob(conn)


	}

}
