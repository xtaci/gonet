DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(20) DEFAULT NULL,
  `mac`   char(11) DEFAULT NULL,
  `score` int(11) DEFAULT '0',
  `last_save_time` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `buildings`;
CREATE TABLE `buildings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `list` varchar(40000),
  `map` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),	
  UNIQUE KEY `uniq_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

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
		INSERT INTO users(name, mac) VALUES(@name, 'AABBCCDDEEFF');
	UNTIL @x > MAX END REPEAT;
END
//
