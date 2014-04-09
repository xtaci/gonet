package stats

var Fields = []string{
	"name",
	"lang",
	"country",
	"domain",
	"device_type",
	"login_cnt",
	"play_time",
	"sign_time",
	"last_login_time",
	"this_login_time",
	"is_login",
	"chat_cnt",
	"email_cnt",
	"vip_score",
	"vip_level",
	"cup_cnt",
	"pve_progress", // pve进度
	"main_task",    // 主线任务完成数
	"achievement",  // 成就任务完成数
	"level",
	"total_pvp_cnt",
	"total_pve_cnt",
	"soldier1_trainning",
	"soldier2_trainning",
	"soldier3_trainning",
	"soldier4_trainning",
	"hero1_pvp_total",
	"hero2_pvp_total",
	"hero3_pvp_total",
	"hero1_pvp_win",
	"hero2_pvp_win",
	"hero3_pvp_win",
	"hero1_pvp_lost",
	"hero2_pvp_lost",
	"hero3_pvp_lost",
	"hero1_pve_total",
	"hero2_pve_total",
	"hero3_pve_total",
	"hero1_pve_win",
	"hero2_pve_win",
	"hero3_pve_win",
	"hero1_pve_lost",
	"hero2_pve_lost",
	"hero3_pve_lost",
	"pve_get_food",
	"pve_get_gold",
	"food_for_build",
	"food_for_soldier",
	"food_for_build_level",
	"nature_get_food",
	"pvp_get_food",
	"pvp_get_gold",
	"pvp_lost_food",
	"pvp_lost_gold",
	"gold_for_build",
	"gold_for_build_level",
	"gold_for_search",
	"gold_for_strength",
	"nature_get_gold",
	"task_get_food",
	"task_get_gold",

	// 装备
	"hero1_armour",
	"hero2_armour",
	"hero3_armour",
	"hero1_weapon",
	"hero2_weapon",
	"hero3_weapon",
	"hero1_banner",
	"hero2_banner",
	"hero3_banner",
	"soldier1_armour",
	"soldier2_armour",
	"soldier3_armour",
	"soldier4_armour",
	"soldier1_weapon",
	"soldier2_weapon",
	"soldier3_weapon",
	"soldier4_weapon",

	// 建筑
	"farm_1",
	"farm_2",
	"farm_3",
	"farm_4",

	"mine_1",
	"mine_2",
	"mine_3",
	"mine_4",

	"granary_1",
	"granary_2",
	"granary_3",

	"treasury_1",
	"treasury_2",
	"treasury_3",

	"pub_1",

	"blacksmith_1",

	"barrack_1",

	"cottage_1",
	"cottage_2",
	"cottage_3",
	"cottage_4",
	"cottage_5",

	"pedrero_1",
	"pedrero_2",
	"pedrero_3",
	"pedrero_4",
	"pedrero_5",
	"pedrero_6",

	"catapult_1",
	"catapult_2",
	"catapult_3",

	"freezingTower_1",
	"freezingTower_2",
	"freezingTower_3",

	"magicTower_1",
	"magicTower_2",
	"magicTower_3",

	"archerTower_1",
	"archerTower_2",
	"archerTower_3",
	"archerTower_4",
	"archerTower_5",

	// 珍珠相关
	"gems_cnt", // 当前拥有的珍珠数
	"gems_consume",
	"gems_get",
	"gems_get_from_task",
	"gems_get_from_pve",
	"gems_get_from_rank_1",
	"gems_get_from_gacha",
	"gems_get_from_buy", // 购买获得珍珠

	"gems_buy_10%_food",
	"gems_buy_10%_food_extra", // 税收系统多出的统计量

	"gems_buy_50%_food",
	"gems_buy_50%_food_extra", // 税收系统多出的统计量

	"gems_buy_100%_food",
	"gems_buy_100%_food_extra", // 税收系统多出的统计量

	"gems_buy_10%_gold",
	"gems_buy_10%_gold_extra", // 税收系统多出的统计量

	"gems_buy_50%_gold",
	"gems_buy_50%_gold_extra", // 税收系统多出的统计量

	"gems_buy_100%_gold",
	"gems_buy_100%_gold_extra", // 税收系统多出的统计量

	"gems_buy_10%_popu",
	"gems_buy_10%_popu_extra", // 税收系统多出的统计量

	"gems_buy_50%_popu",
	"gems_buy_50%_popu_extra", // 税收系统多出的统计量

	"gems_buy_100%_popu",
	"gems_buy_100%_popu_extra", // 税收系统多出的统计量

	"gems_buy_soldier",
	"gems_buy_soldier_extra",

	"gems_buy_10%_protect",
	"gems_buy_10%_protect_extra",

	"gems_buy_50%_protect",
	"gems_buy_50%_protect_extra",

	"gems_buy_100%_protect",
	"gems_buy_100%_protect_extra",

	"gems_buy_name",
	"gems_buy_name_extra",

	"gems_buy_buff",
	"gems_buy_buff_extra",

	"gems_buy_ornament", // 购买装饰物
	"gems_buy_ornament_extra",
	"gems_create_alliance",
	// 补全
	"complement_food",
	"complement_food_extra",

	"complement_gold",
	"complement_gold_extra",

	"complement_popu",
	"complement_popu_extra",

	"uuid",

	"gold_add_from_buy",
	"food_add_from_buy",
	"popu_add_from_buy",

	"task_get_popu",

	"gold_add_from_item",
	"food_add_from_item",
	"popu_add_from_item",

	"popu_for_build",
	"popu_for_build_level",
}

var Conver map[string]int

func init() {
	Conver = make(map[string]int)
	for k := range Fields {
		Conver[Fields[k]] = k
	}
}

/*
var Conver = map[string]int{
	"name":                 0,
	"lang":                 1,
	"country":              2,
	"domain":               3,
	"device_type":          4,
	"login_cnt":            5,
	"play_time":            6,
	"sign_time":            7,
	"last_login_time":      8,
	"this_login_time":      9,
	"is_login":             10,
	"chat_cnt":             11,
	"email_cnt":            12,
	"vip_score":            13,
	"vip_level":            14,
	"cup_cnt":              15,
	"pve_progress":         16, // pve进度
	"main_task":            17, // 主线任务完成数
	"achievement":          18, // 成就任务完成数
	"level":                19,
	"total_pvp_cnt":        20,
	"total_pve_cnt":        21,
	"soldier1_trainning":   22,
	"soldier2_trainning":   23,
	"soldier3_trainning":   24,
	"soldier4_trainning":   25,
	"hero1_pvp_total":      26,
	"hero2_pvp_total":      27,
	"hero3_pvp_total":      28,
	"hero1_pvp_win":        29,
	"hero2_pvp_win":        30,
	"hero3_pvp_win":        31,
	"hero1_pvp_lost":       32,
	"hero2_pvp_lost":       33,
	"hero3_pvp_lost":       34,
	"hero1_pve_total":      35,
	"hero2_pve_total":      36,
	"hero3_pve_total":      37,
	"hero1_pve_win":        38,
	"hero2_pve_win":        39,
	"hero3_pve_win":        40,
	"hero1_pve_lost":       41,
	"hero2_pve_lost":       42,
	"hero3_pve_lost":       43,
	"pve_get_food":         44,
	"pve_get_gold":         45,
	"food_for_build":       46,
	"food_for_soldier":     47,
	"food_for_build_level": 48,
	"nature_get_food":      49,
	"pvp_get_food":         50,
	"pvp_get_gold":         51,
	"pvp_lost_food":        52,
	"pvp_lost_gold":        53,
	"gold_for_build":       54,
	"gold_for_build_level": 55,
	"gold_for_search":      56,
	"gold_for_strength":    57,
	"nature_get_gold":      58,
	"task_get_food":        59,
	"task_get_gold":        60,

	// 装备
	"hero1_armour":    61,
	"hero2_armour":    62,
	"hero3_armour":    63,
	"hero1_weapon":    64,
	"hero2_weapon":    65,
	"hero3_weapon":    66,
	"hero1_banner":    67,
	"hero2_banner":    68,
	"hero3_banner":    69,
	"soldier1_armour": 70,
	"soldier2_armour": 71,
	"soldier3_armour": 72,
	"soldier4_armour": 73,
	"soldier1_weapon": 74,
	"soldier2_weapon": 75,
	"soldier3_weapon": 76,
	"soldier4_weapon": 77,

	// 建筑
	"farm_1": 78,
	"farm_2": 79,
	"farm_3": 80,
	"farm_4": 81,

	"mine_1": 82,
	"mine_2": 83,
	"mine_3": 84,
	"mine_4": 85,

	"granary_1": 86,
	"granary_2": 87,
	"granary_3": 88,

	"treasury_1": 89,
	"treasury_2": 90,
	"treasury_3": 91,

	"pub_1": 92,

	"blacksmith_1": 93,

	"barrack_1": 94,

	"cottage_1": 95,
	"cottage_2": 96,
	"cottage_3": 97,
	"cottage_3": 97,

	"pedrero_1": 98,
	"pedrero_2": 99,
	"pedrero_3": 100,
	"pedrero_4": 101,
	"pedrero_5": 102,
	"pedrero_6": 103,

	"catapult_1": 104,
	"catapult_2": 105,
	"catapult_3": 106,

	"freezingTower_1": 107,
	"freezingTower_2": 108,
	"freezingTower_3": 109,

	"magicTower_1": 110,
	"magicTower_2": 111,
	"magicTower_3": 112,

	"archerTower_1": 113,
	"archerTower_2": 114,
	"archerTower_3": 115,
	"archerTower_4": 116,
	"archerTower_5": 117,

	// 珍珠相关
	"gems_cnt":             118, // 当前拥有的珍珠数
	"gems_consume":         119,
	"gems_get":             120,
	"gems_get_from_task":   121,
	"gems_get_from_pve":    122,
	"gems_get_from_rank_1": 123,
	"gems_get_from_gacha":  124,
	"gems_get_from_buy":    125, // 购买获得珍珠

	"gems_buy_10%_food":       126,
	"gems_buy_10%_food_extra": 127, // 税收系统多出的统计量

	"gems_buy_50%_food":       128,
	"gems_buy_50%_food_extra": 129, // 税收系统多出的统计量

	"gems_buy_100%_food":       130,
	"gems_buy_100%_food_extra": 131, // 税收系统多出的统计量

	"gems_buy_10%_gold":       132,
	"gems_buy_10%_gold_extra": 133, // 税收系统多出的统计量

	"gems_buy_50%_gold":       134,
	"gems_buy_50%_gold_extra": 135, // 税收系统多出的统计量

	"gems_buy_100%_gold":       136,
	"gems_buy_100%_gold_extra": 137, // 税收系统多出的统计量

	"gems_buy_10%_popu":       138,
	"gems_buy_10%_popu_extra": 139, // 税收系统多出的统计量

	"gems_buy_50%_popu":       140,
	"gems_buy_50%_popu_extra": 141, // 税收系统多出的统计量

	"gems_buy_100%_popu":       142,
	"gems_buy_100%_popu_extra": 143, // 税收系统多出的统计量

	"gems_buy_soldier":       144,
	"gems_buy_soldier_extra": 145,

	"gems_buy_10%_protect":       146,
	"gems_buy_10%_protect_extra": 147,

	"gems_buy_50%_protect":       148,
	"gems_buy_50%_protect_extra": 149,

	"gems_buy_100%_protect":       150,
	"gems_buy_100%_protect_extra": 151,

	"gems_buy_name":       152,
	"gems_buy_name_extra": 153,

	"gems_buy_buff":       154,
	"gems_buy_buff_extra": 155,

	"gems_buy_ornament":       156, // 购买装饰物
	"gems_buy_ornament_extra": 157,
	// 补全
	"complement_food":       158,
	"complement_food_extra": 159,

	"complement_gold":       160,
	"complement_gold_extra": 161,

	"complement_popu":       162,
	"complement_popu_extra": 163,

	"uuid":              164,
	"gold_add_from_buy": 165,
	"food_add_from_buy": 166,
	"popu_add_from_buy": 167,

	"task_get_popu": 168,

	"gold_add_from_item": 169,
	"food_add_from_item": 170,
	"popu_add_from_item": 171,

	"popu_for_build":       172,
	"popu_for_build_level": 173,
}
*/
