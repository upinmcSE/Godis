## Implement I/O Multiplexing

- Client connect đến port 3000, epoll notify server accept. Sau khi server accept thì socket được tạo và đăng kí với poll để lắng nghe. Khi client gửi dữ liệu, dữ liệu được ghi vào socket -> epoll notify server -> server xử lý rồi ghi vào socket để trả về. 

- Khi ngắt kết nối -> socket bị xóa -> epoll remove ? 


## Benmark

./redis/src/redis-benchmark -n 10000 -t ping_mbulk -c 200 -h localhost -p 3000

TCP Server with thread-per-connection + blocking IO:

Summary:
  throughput summary: 69444.45 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        2.335     0.048     1.863     4.695    10.183    15.247



IO Multiplexing:

Summary:
  throughput summary: 84745.77 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        1.437     0.344     1.063     3.591     5.959     6.823



## RESP
- RESP là protocol phục vụ cho Redis client giao tiếp với Redis server
- Reliable (vì là build on top of TCP)
- Simple to implement
- Fast to parse
- Human-readable

- Support different types of data:
  - string
  - array
  - integer, double
  - boolean
  - error
  - ...

### Simple String


### Bulk String


### Integer
- 

### Array

### Error
-
-
-

## Q/A
- RESP balance giữa readable với high performance vậy ạ. Tại sao json cũng readable nhưng mà performance lại kém vậy ạ

```
"hello" =>
{\r\n
    type: 0,
    value: "hello"
}
```
```
"hello" => +hello\r\n
```

=> trông rất cồng cành và tốn resource(nhiều byte hơn so với RESP) 
- 