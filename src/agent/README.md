### 游戏服务器(GS)设计:      
一个玩家对应一个goroutine(称之为agent), 每个goroutine接收信息类型包含：     

1. 来自客户端的    
2. 来自其他玩家的    
3. 时钟信号

### 持久化设计
1. 玩家在线时agent负责变更内存中玩家数据，并持久化到数据库，   
2. 升级逻辑以agent为准, 即agent直接改变数据库中的数据.    
3. 在线时数据持久化策略：   
    a) 超过一定的操作数量，刷入数据库    
    b) 超过一定的时间，刷入数据库    
    c) 离线时，刷入数据库
4. 事件服务器通过CAS(Compare-and-Set) 的方式变更升级数据, 保证玩家在线时逻辑处理的一致性.     

#### 包结构
<pre>
|LENGTH|TIME_ELAPSED|PROTO|PAYLOAD|

|16|32|16|...|
</pre>

#### 包测试：
<pre>
heart_beat:
echo "0006 00000001 0000" | xxd -r -ps |nc 127.0.0.1 8888 -q 2|hexdump -C
</pre>
