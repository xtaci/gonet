###A game server implemented with golang.
![Architecture](doc/arch.png)

#### 部署:     
* Game Server(GS):  
玩家直接连接GS， 处理玩家逻辑，并与HUB/EventServer/StatsServer通信，GS存在若干个。     
  
* Hub Server(HUB):  
若干个GS 连接到一个HUB, 只存在一个HUB，维护基础的全局信息，以及 GS<--->GS 的消息转发.  
    
* Event Server(ES):  
事件服务器，完成例如建筑升级，招募，研究等需要升级等待的操作。
事件被触发后会直接访问数据库，计算变更后，尝试性写入，如果玩家在线，结果以GS的为准。

* Stats Server(SS):     
统计服务器，根据玩家的行为，记录策划需要的数据，以便于后期统计。     
统计属于事后分析，数据量较大，性能需求不同, 故单独列为一个服务器。
统计服务器会访问到两个mongodb数据库，玩家数据库(只读)和统计数据库，最终会汇入统计数据库。

#### 通信原则:     
1.  GS到HUB/SS/ES的通信，都是Call同步调用，即GS必须等待ACK。         
2.  HUB到GS的通信，只有forward数据包。       
3.  单播消息在玩家离线时会存入db, 登录后的启动过程 ___GS___ 直接读取db，并forward给玩家goroutine。
4.  多播消息在只保留一个固定长度的FIFO, 登录后，多播消息的未读部分 ___HUB___ 会直接forward给玩家goroutine

#### 服务器状态一致性
1.  GS节点可以单独重启    
2.  HUB/ES/SS 重启后，GS必须全部重启。    

#### 安装先决条件:
0. 确保安装好bzr, graphviz, gawk
1. 确保安装好mongodb
2. 确保config.ini中的mongo_xxxx配置正确
3. export GOPATH='当前目录'

#### 安装:
* xtaci@ubuntu:~/gonet$ go get labix.org/v2/mgo      
* xtaci@ubuntu:~/gonet$ go get github.com/stevedonovan/luar
* xtaci@ubuntu:~/gonet$ go get github.com/aarzilli/golua/lua
* xtaci@ubuntu:~/gonet$ make    
* xtaci@ubuntu:~/gonet$ ./startall  

#### 启动:
1. 启动支持多config配置，如 hub --config=/path/to/config.ini
2. 放入后台运行的方式，直接加 '&', 如  hub &
3. 详细配置请直接阅读和修改config.ini
