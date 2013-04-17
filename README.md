server implemented with Go language.

mysql -uroot -p game < db.sql

mysql>use game;

mysql>call gen_users(10000);

mysql>select count(*) from users;

sh$go get github.com/ziutek/mymysql/mysql

sh$go get github.com/ziutek/mymysql/native

sh$go install gate

sh$gate
