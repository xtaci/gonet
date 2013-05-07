server implemented with Go language.

mysql -uroot -p game < db.sql

sh$go get github.com/ziutek/mymysql/mysql

sh$go get github.com/ziutek/mymysql/native

sh$cd src/scripts;./proto_gen.sh

sh$go install gate

sh$gate


gonet只保证表级的数据操作原子性,

跨表操作的原子性不保证
