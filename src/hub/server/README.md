GS <---> HUB 

通信协议

|LENGTH|SEQNUM|PROTO|PAYLOAD|

|16|32|16|...|

such as...

echo "000D 0001 0001 41 00000001 01 0001 42" | xxd -r -ps |nc 127.0.0.1 8888 -q 1|hexdump -C
