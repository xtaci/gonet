Game Logic

you can test packet with 

 echo "00 02 00 05" | xxd -r -ps |nc 127.0.0.1 8888
 size:2 proto:5

数据同步规则:

1. Read尽可能在内存中进行
2. Write必须立即Sync到DB
