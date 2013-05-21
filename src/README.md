GONET 总体设计
=============

gonet是一个游戏服务器，满足类似于与clash of clans这种游戏类型的应用。

架构：
====

Game Server:

    玩家直接连接GS， 处理玩家逻辑，玩家逻辑是可以并行的（打副本），即不存在交互。交互部分由HUB完成。
  
Hub Server:

    若干个GS 连接到一个HUB, 理论上只存在一个HUB，维护玩家全局的信息。以及 GS<--->GS 的交互.
    
CoolDown Server:

    冷却服务器，完成建筑升级，招募，研究等需要升级等待的操作。
    
