###A game server implemented with golang.

####条件
0. 确保安装好bzr
1. 确保安装好mongodb
2. 确保config.ini中的mongo_xxxx配置正确
3. export GOPATH='当前目录'

####安装
xtaci@ubuntu:~/gonet$ go get labix.org/v2/mgo      
xtaci@ubuntu:~/gonet$ make    
xtaci@ubuntu:~/gonet$ ./startall  
