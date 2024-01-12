# redis 配置

## redis主从配置

### redis主节点配置

```
1.配置pidfile（一般本地有多个redis服务的时候改下名字）
2.配置port （不同的redis服务需要有不同的port，默认为6379）
2.配置dir （rdb保存路径）
3.配置dbfilename（rdb文件名）
4.配置save （rdb快照保存频率和条件）
```

### redis从节点配置

```
1.和主节点配置步骤相同
2.配置replicaof （redis从属于哪一个主节点）
```
