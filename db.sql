DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(20) DEFAULT NULL,
  `mac`   char(11) DEFAULT NULL,
  `score` int(11) DEFAULT '0',
  `archives` varchar(10240) DEFAULT NULL,
  `last_save_time` datetime DEFAULT NULL,
  `protect_time` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uc_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP PROCEDURE IF EXISTS gen_users;
DELIMITER //
CREATE PROCEDURE gen_users(MAX INT)
BEGIN
	DECLARE x INT;
	DECLARE name VARCHAR(10);
	
	SET @x = 0;
	
	REPEAT
		SET @x = @x + 1;
		SET @name = CONCAT('user', @x);
		INSERT INTO users(name, password) VALUES(@name, MD5(@name));
	UNTIL @x > MAX END REPEAT;
END
//
