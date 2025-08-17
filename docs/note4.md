## PING

- [Redis-CLI] ---send---> [data byte[]] ---RESP Decoder---> [Command: array of string] ----> [Execute Command] ---RESP Encoder---> [data byte[]] ----reply---> [Redis-CLI]


## What is store data?

```
PING hello
{
    Cmd: "PING"
    Args: ["hello"]
}

SET key 10
{
    Cmd: "SET"
    Args: ["key", "10"]
}
```

- using hash table for store data
  - key: String
  - value: interface{}

  - expiredTime = time.Now.UnixMilli() + ttl

- GET:
  - check exist
  - check expire
  - return
- SET:
  - calculate expire time
  - insert to hash tables

## Delete expired key
- Passive mode : delete on access
  - delete key when user calls GET command
- Active mode: 
  - implement idea of redis : https://valkey.io/commands/expire/ ; https://github.com/redis/redis/blob/b9d9d4000b0b45b87c7c6fea23ec9fd8fcac107e/src/expire.c#L187

### If no event occurs to Wait() then the ActiveExpired function does not run?
- Yes