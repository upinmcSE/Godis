
## Thread per connect

```go
func process(conn net.Conn){
	log.Println(conn.RemoteAddr())

	// read data from client
	var buf []byte = make([]byte, 1000)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// process
	time.Sleep((time.Second * 1))

	// reply
	//conn.Write([]byte("Hello, World"))
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, World\r\n"))
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// conn == socket
		conn, err := listener.Accept()
		if err != nil{
			log.Fatal(err)
		}

		//process(conn)

		// create go routine to handle the connection
		go process(conn)
	}
}

```

```go
package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr())

	// Tạo scanner để đọc dữ liệu từ client theo dòng
	scanner := bufio.NewScanner(conn)

	// Vòng lặp để đọc dữ liệu liên tục từ client
	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text()) // Lấy message và loại bỏ khoảng trắng
		log.Printf("Received from %v: %s", conn.RemoteAddr(), message)


		if message == "stop" {
			log.Printf("Client %v sent stop, closing connection", conn.RemoteAddr())
			conn.Write([]byte("Connection closed\r\n"))
			break
		}

		// Xử lý message (ví dụ: giả lập thời gian xử lý)
		time.Sleep(time.Second * 1)

		
		response := "HTTP/1.1 200 OK\r\n\r\nHello, World\r\n"
		_, err := conn.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing to %v: %v", conn.RemoteAddr(), err)
			break
		}
	}

	// Kiểm tra lỗi từ scanner
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from %v: %v", conn.RemoteAddr(), err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server listening on :3000")

	for {
		// Chấp nhận kết nối mới
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Xử lý kết nối trong một goroutine
		go process(conn)
	}
}
```

## Thread pool

```go
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
			log.Println("Worker &d is handling job from %s", w.id, job.conn.RemoteAddr())
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
```

## Everything is a file
- a file  == a stream of bytes
- 

### why is it helpful?
- 
- 
- 
- 

### File descetiptor

## IO models

### Blocking IO
-

### Non-Blocking IO
- 
- need to keep asking

### Async IO
- 

- pros: 
- cons:

### IO Multiplexing
- 

#### Linux System calls

#### MacOS System calls