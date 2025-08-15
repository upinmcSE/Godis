## What is Redis? (REmote DIctionary Server)

- open source, in-memory, NoSql Key-Value

### Why Redis?
 - free
 - performant
   - low latency (nhanh)
   - high throughput
 - simple (chỉ là đọc/ghi key-value)
 - support multiple data structures

### Why not Redis?
 - Potential for data loss (lưu trong ram nên tắt máy thì mất hết dữ liệu)
 - No complex query language (không query phức tạp được như sql)
 - Single-threaded nature (chỉ dùng 1 thread)
 - Not ideal for complex data relationships

### What for?
- Caching
    - 80/20 rule (80% time of users usually 20% popular data)
- Session storage (lưu phiên đăng nhập)
- Distributed lock (trường hợp có nhều server và nhiều task - > nếu task có server đang xử lý thì sẽ khóa task lại không cho server khác xử lý)


## What is going on in a server?

- Hardware
    - CPU
    - Memory (RAM)
    - I/O
    - Input -> | CPU <--> RAM | -> Output

- Software 
    - Users -> Applications

- Hardware <-- OS --> Applications
- Application -> System calls(API) -> OS
  - System call: open(), read(), write(), kill(), fork(), exec(), ...

- Program (chương trình):
    - Excutable code(binary file) -> 1 đoạn code
    - A static entity
- Process (tiến trình):
    - An instance of a program in execution
    - Each process runs one program
    - One program can run in multiple processes
    - Act like a mini virtual server
        - CPU: time slicing of processor
        - RAM: address space

- Program ~ Class, Process ~ Object trong Java

- OS Thread (Luồng)
    - An object of activity within the process (mỗi process có nhiều thread)
    - smallest sequence of programmed instructions that can be managed independently by OS scheduler
    - managed by OS

- User space thread
    - Threads managed by thread library in user space
    - e.g JVM, golang runtime
    - Lightweight and faster than kernel thread (nhẹ và nhanh)

- style 1: One-to-One : user space thread <-> kernel space thread
- style 2: Many-to-Many : user space threads <-> kernel space threads

- Context-switch
  - vì ..

#### context switching trong mô hình many-to-many có tốn resource hơn so với one-to-one không?

- Multi-threading programming
    - pros
        - Better performance : concurrency/parallelism
        - Better resource utilization 
    - cons
        - complex, hard to debug (phức tạp, khó code)
        - race condition

- Race condition
    - locking
        - lock contention bottleneck (khi nhiều threads concurrently mà gặp lock thì bị nghẽn cổ chai)
        - context switching overhead (khi thread bị từ chối task do đang lock thì thread sẽ chuyển sang sleep)
        - complex, debugging difficulties (Heisenbug)
    - atomic variable ( là một biến mà các thao tác trên nó luôn được thực hiện một cách nguyên tử, tức là không thể bị chia cắt bởi các luồng khác )

### Mutex lock nó cũng là 1 resource, tại sao lúc nhiều thread gọi đến nó thì không bị race condition?
- tại 1 thời điểm nó chỉ acquire 1 thằng nên không bị race condition

## Communication

- How do servers talk to each other?
    - Server can use different software and hardware ()
    - A Protocol defines the format and the order of messages exchanged between two or more communicating entities
    - Network protocols are like a common language for computers.

### Protocol
- TCP
  - Là 1 Protocol để 2 máy tính giao tiếp với nhau 1 cách reliable
  
- UDP
  - fast but unreliable

- HTTP/HTTPS (build on top of TCP)
- ...

### Redis Protocol
- RESP (build on top of TCP)
  - why doesn't Redis use HTTP protocol?
      - Minimize overhead: verbose text-based headers, various status codes, content nagotiation, cookies, etc
      - Efficient Parsing: simple, much easier to parse than JSON/XML
      - No web browser access : a backend service (end-user của redis là 1 server chứ không phải client)

### Client-Server Application model
- Client <--- request|response ---> Server

  - 1. thread per connection model (mỗi request gửi lên server thì server sẽ tạo ra 1 thread riêng để xử lý)
      - Pros:
          - simple
          - leverages Multi-Core Process ( tận dụng được hết core-cpu)
          - handles Blocking I/O
      - Cons:
          - Memory overhead: each thread requires its own stack space (2mb/thread)
          - CPU overhead: context switching
          - Rick of race condition
  - 2. thread pool ( tạo sẵn các threads ngồi đợi task)
      - Pros
          - avoid overload the hardware ( kiểm soát được số lượng threads)
        
      - Cons
          - hard to choose the pool's size
          - a bit more complex
          - risk of race condition
          - overhead for very short tasks (các task nhẹ pha đợi task nặng)
  - 3. Event driven
      - A request comes in, and the event loop dispatches it
      - If the request requires an I/O operator(e.g, database query, API call), the operation is initiated, but the thread doesn't wait. It register a "callback" to be executed when the I/O completes
      - The event loop immediately moves on to process other requests or events
      - When an I/O operation finishes, the OS notifies the event loop, which then takes the completed data and executes the associated callback

      - requests -> [Event Queue] ->[Event Loop(thread)] -> [Thread Pool]
    
      - Pros:
          - Scalable for I/O bound application (vì nó không bao giờ bị block)
          - Efficient Resource Usage: avoid context switching
          - Reduce race condition risk
      - Cons: 
          - Complex
          - CPU-Bound Operations Block Everything

      - Redis use a variety of Event Driven model

## TCP Server

- What is TCP Server/
    - 
    - 
    - (bắt tay 3 bước)
    - 