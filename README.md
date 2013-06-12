###A game server implemented with golang.
![Architecture](/doc/arch.png)

* Game Server:  
玩家直接连接GS， 处理玩家逻辑，逻辑结果直接存入数据库，升级交给Event服务器处理，玩家交互由HUB完成。     
（GS的操作均为同步操作。）
  
* Hub Server:  
若干个GS 连接到一个HUB, 理论上只存在一个HUB，维护基础的玩家全局的信息。以及 GS<--->GS 的交互.  
基本信息持久保持在内存中。   
    
* Event Server:  
事件服务器，完成例如建筑升级，招募，研究等需要升级等待的操作。    
玩家事件在触发后会读取数据库，计算变更后，写入。 （ES的操作均为异步操作。）

* Stats Server:     
统计服务器，根据玩家的行为，记录策划需要的数据，以便于后期统计。     
统计属于事后分析，数据量较大，性能需求不同, 故单独列为一个服务器。
统计服务器会访问到两个mongodb数据库，玩家数据库(只读)和统计数据库，最终会汇入统计数据库。

####条件
0. 确保安装好bzr, graphviz
1. 确保安装好mongodb
2. 确保config.ini中的mongo_xxxx配置正确
3. export GOPATH='当前目录'

####安装
* xtaci@ubuntu:~/gonet$ go get labix.org/v2/mgo      
* xtaci@ubuntu:~/gonet$ make    
* xtaci@ubuntu:~/gonet$ ./startall  

####启动
1. 启动支持多config配置，如 hub --config=/path/to/config.ini
2. 放入后台运行的方式，直接加 '&', 如  hub &
3. 详细配置请直接阅读和修改config.ini
