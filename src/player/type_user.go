package player

type User struct {
	mq chan string;
	id int;
	name string;
    owner_id int;
	is_capital int;
	achievement int;
	x int;
	y int;
	gold int;
	wood int;
	food int;
	iron int;
	stone int;
	scout int;
	swordsman int;
	crossbow_archer int;
	squire int;
	templar int;
	paladin int;
	archer_cavalry int;
	royal_knight int;
	action_events_count int;
	recruit_events_count int;
	deals_count int;
	lock_version int;
	last_move_time int;
	durability int;
	arcane_mage int;
	battle_mage  int;
	holy_mage int;
	is_auto_fix int;
	revive_time int;
	item_warehouse_lv int;
	item_transport_lv int;
	skeleton int;
	ghost_rider int;
	ram int;
	zeppelin int;
	steel_golem int;
	cruiser int
}
