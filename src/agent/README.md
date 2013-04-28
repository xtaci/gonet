Game Logic

you can test packet with 

 echo "000D 0001 0001 41 00000001 01 0001 42" | xxd -r -ps |nc 127.0.0.1 8888 -q 1|hexdump

数据同步规则:

1. Read尽可能在内存中进行
2. Write必须立即Sync到DB
