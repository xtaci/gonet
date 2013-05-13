GONET 总体设计
=============

gonet是一个游戏服务器，满足类似于与clash of clans这种游戏类型的应用。

架构：
====

Game Server:

    玩家直接连接GS， 处理玩家逻辑，玩家逻辑是可以并行的（打副本），即不存在交互。交互部分由HUB完成。
  
Hub Server:

    若干个GS 连接到一个HUB, 理论上只存在一个HUB，维护玩家全局的信息。以及 GS<--->GS 的交互.
    

数据同步原则
====
一个数据只能由某一个goroutine管理。该goroutine负责数据的读写以及持久化。
    
