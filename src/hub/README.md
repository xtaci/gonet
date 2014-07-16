### HUB 服务器

HUB 作为一个独立的服务运行，各个逻辑服务器(GS)启动时会连接到HUB.      
HUB 管理玩家基础、重要的状态信息的存取变更，如：    

* 排名    
* 状态
* 联盟消息管理   

HUB也承担玩家之间消息转发（不在同一个服务器登陆的情况）。     
HUB只处理来自GS的两类消息：   

1. 来自Game Server 的Call请求(request & ack)     
2. Game Server 间的消息Forward      

### 设计考虑
1. 玩家动态排名,基于动态有序统计      
2. 玩家状态机是行级锁实现       
3. 联盟管理保留最大group_msg_max这么多条消息。玩家登陆后，联盟的消息会转发过来。

![状态机](https://github.com/xtaci/gonet/raw/develop/doc/fsm.png)
