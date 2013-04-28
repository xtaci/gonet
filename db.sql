DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(20) DEFAULT NULL,
  `mac`   char(11) DEFAULT NULL,
  `score` int(11) DEFAULT '0',
  `state` int(11) DEFAULT '0',
  `last_save_time` datetime DEFAULT NULL,
  `protect_time` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uc_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `cities`;
CREATE TABLE `cities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `owner_id` int(11) DEFAULT NULL,
  `is_capital` tinyint(1) DEFAULT '0',
  `type` varchar(255) COLLATE utf8_unicode_ci DEFAULT 'City',
  `achievement` int(11) DEFAULT '0',
  `x` int(11) DEFAULT NULL,
  `y` int(11) DEFAULT NULL,
  `gold` int(11) DEFAULT '0',
  `wood` int(11) DEFAULT '0',
  `food` int(11) DEFAULT '0',
  `iron` int(11) DEFAULT '0',
  `stone` int(11) DEFAULT '0',
  `workers` text COLLATE utf8_unicode_ci,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `protected_until` datetime DEFAULT NULL,
  `scout` int(11) DEFAULT '0',
  `swordsman` int(11) DEFAULT '0',
  `crossbow_archer` int(11) DEFAULT '0',
  `items_effect` text COLLATE utf8_unicode_ci,
  `updated_resource_at` datetime DEFAULT NULL,
  `squire` int(11) DEFAULT '0',
  `templar` int(11) DEFAULT '0',
  `paladin` int(11) DEFAULT '0',
  `archer_cavalry` int(11) DEFAULT '0',
  `royal_knight` int(11) DEFAULT '0',
  `leadership` float DEFAULT '1000',
  `skills_effect` text COLLATE utf8_unicode_ci,
  `action_events_count` int(11) DEFAULT '0',
  `recruit_events_count` int(11) DEFAULT '0',
  `deals_count` int(11) DEFAULT '0',
  `lock_version` int(11) DEFAULT '0',
  `last_move_time` int(11) DEFAULT '0',
  `durability` int(11) DEFAULT '0',
  `arcane_mage` int(11) DEFAULT '0',
  `battle_mage` int(11) DEFAULT '0',
  `holy_mage` int(11) DEFAULT '0',
  `is_auto_fix` tinyint(1) DEFAULT '0',
  `revive_time` int(11) DEFAULT '0',
  `item_warehouse_lv` int(11) DEFAULT '0',
  `item_transport_lv` int(11) DEFAULT '0',
  `skeleton` int(11) DEFAULT '0',
  `ghost_rider` int(11) DEFAULT '0',
  `ram` int(11) DEFAULT '0',
  `zeppelin` int(11) DEFAULT '0',
  `steel_golem` int(11) DEFAULT '0',
  `cruiser` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_cities_on_name` (`name`),
  KEY `index_cities_on_owner_id` (`owner_id`),
  KEY `index_cities_on_type` (`type`),
  KEY `index_cities_on_x` (`x`),
  KEY `index_cities_on_y` (`y`)
) ENGINE=InnoDB CHARSET=utf8;

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
