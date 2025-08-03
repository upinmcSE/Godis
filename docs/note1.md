## What is Redis?
- open source, in-memory, NoSql Key-Value

 ## Why Redis?
 - free
 - performant
 - simple
 - support multiple data structures

 ## Why not Redis?
 - Potential for data loss
 - No complex query language
 - Single-threaded nature
 - Not ideal for complex data relationships

 ## What for?
- Caching
    - 80/20 rule
- Session storage
- Distributed lock


## What is going on in a server?

- Hardware
    - CPU
    - Memory (RAM)
    - I/O
    - Input -> | CPU <--> RAM | -> Output

- Software 
    - Users -> Applications

- Hardware <-- OS --> Applications
- Application -> System calls -> OS

- Program
    - Excutable code(binary file)
    - A static entity
- Process:
    - An instance of a program in execution
    - Each process runs one program
    - One program can run in multiple processes
    - Act like a mini virtual server
        - CPU: time slicing of processor
        - RAM: address space

- Program ~ Class, Process ~ Object trong Java

- OS Thread (Luồng)
    - 
    - 
    - 

- User space thread
    - 
    -
    -

#### context switching trong mô hình many-to-many có tốn resource hơn so với one-to-one không?

- Multi-threading programming
    - pros
        -
        -
    - cons
        - complex, hard to debug
        - race condition

- Race condition
    - locking
        - lock contention bottleneck
        - context switching overhead
        - complex, debugging difficulties (Heisenbug)
    - atomic variable

### Mutex lock nó cũng là 1 resource, tại sao lúc nhiều thread gọi đến nó thì không bị race condition?

## Communication

- How do servers talk to each other?
    - 
    - 
    - 
- Protocol
    - TCP
        - 
    - UDP
        - 

    - HTTP/HTTPS (build on top of TCP)
    - ...

- Redis Protocol
    - RESP (build on top of TCP)
    - why doesn't Redis use HTTP protocol?
        - 
        -
        -
- Client-Server Application model
    - Client <--- request|response ---> Server

    - 1. thread per connection model
        - Pros:
            -
            -
            -
        - Cons:
            -
            - 
            - 
    - 2. thread pool
        - Pros
            -
        
        - Cons
            - 
            -   
            - 
            - 
    - 3. Event driven
        -
        -
        -
        -

        - Pros:
            -
            -
            -
        - Cons: 
            -
            -

## TCP Server

- What is TCP Server/
    - 
    - 
    - (bắt tay 3 bước)
    - 