package main

import (
	"db/stats_tbl"
	"encoding/json"
	"fmt"
	. "helper"
	"net/http"
	"strconv"
	"strings"
	"types/stats"
)

const (
	DATE_PREM = 19 // 日期前缀
)

//------------------------------------------------------------ 生成结果
func _gen_result(respones, code, message, status interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"status":   status,
		"code":     code,
		"message":  message,
		"response": respones,
	}
	return result
}

//------------------------------------------------------------ 得到userid，查询开始结束时间字段
func _get_query_field(w http.ResponseWriter, r *http.Request, query string) (uid, start, end, Err int) {
	r.ParseForm()
	if r.Method != query {
		Err = 1
		fmt.Fprintf(w, `{"status" : "failed","code":"404","message":"请求方式不对","response":null}`)
		return uid, start, end, Err
	}
	uid, err := strconv.Atoi(r.FormValue("userid"))
	if err != nil {
		Err = 2
		fmt.Fprintf(w, `{"status" : "failed","code":"404","message":"userid 错误","response":null}`)
		return uid, start, end, Err
	}

	// 查询开始时间
	start, err = strconv.Atoi(r.FormValue("start_date"))
	if err != nil {
		Err = 3
		fmt.Fprintf(w, `{"status" : "failed","code":"404","message":"start_date err","response":null}`)
		return uid, start, end, Err
	}

	// 查询结束时间
	end, err = strconv.Atoi(r.FormValue("end_date"))
	if err != nil {
		Err = 4
		fmt.Fprintf(w, `{"status" : "failed","code":"404","message":"end_date err","response":null}`)
		return uid, start, end, Err
	}
	return uid, start, end, Err
}

//------------------------------------------------------------ 英雄解锁记录
func GetHeroUnlockInfo(w http.ResponseWriter, r *http.Request) {
	uid, start, end, errno := _get_query_field(w, r, "GET")
	if errno != 0 {
		return
	}

	heros := make(map[string][]stats.IntGameInfo)
	heros["hero1"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("hero1_unlock"), start, end)
	heros["hero2"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("hero2_unlock"), start, end)
	heros["hero3"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("hero3_unlock"), start, end)
	type HeroUnlockInfo struct {
		Name  string
		Date  string
		Level string
	}
	infos := make([]HeroUnlockInfo, 0)
	for kval, val := range heros {
		for v := range val {
			str := fmt.Sprint(val[v].Time)
			array_str := []byte(str)
			date := string(array_str[:DATE_PREM])
			tmp := HeroUnlockInfo{
				Name:  kval,
				Level: "1",
				Date:  date,
			}
			infos = append(infos, tmp)
		}
	}

	result := _gen_result(infos, 0, "成功", "success")

	b, erro := json.Marshal(result)
	if erro != nil {
		ERR(erro)
		return
	}
	w.Write(b)
}

//------------------------------------------------------------ 士兵解锁记录
func GetSoldierUnlockInfo(w http.ResponseWriter, r *http.Request) {
	uid, start, end, errno := _get_query_field(w, r, "GET")
	if errno != 0 {
		return
	}
	soldiers := make(map[string][]stats.IntGameInfo, 0)
	soldiers["soldier1"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("soldier1_unlock"), start, end)
	soldiers["soldier2"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("soldier2_unlock"), start, end)
	soldiers["soldier3"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("soldier3_unlock"), start, end)
	soldiers["soldier4"] = stats_tbl.GetPlayLog(uid, fmt.Sprint("soldier4_unlock"), start, end)
	type SoldierUnlockInfo struct {
		Name  string
		Date  string
		Level string
	}
	infos := make([]SoldierUnlockInfo, 0)
	for kval, val := range soldiers {
		for v := range val {
			str := fmt.Sprint(val[v].Time)
			array_str := []byte(str)
			date := string(array_str[:DATE_PREM])
			tmp := SoldierUnlockInfo{
				Name:  kval,
				Level: "1",
				Date:  date,
			}
			infos = append(infos, tmp)
		}
	}

	data := _gen_result(infos, 0, "get success", "success")
	result, _ := json.Marshal(data)
	w.Write(result)
}

//------------------------------------------------------------ vip积分记录
func GetVipScore(w http.ResponseWriter, r *http.Request) {
	uid, start, end, err := _get_query_field(w, r, "GET")
	if err != 0 {
		return
	}
	type ScoreInfo struct {
		Date string
		Nums string
	}
	infos := make([]ScoreInfo, 0)
	scores := stats_tbl.GetPlayLog(uid, "vip_score", start, end)
	for k := range scores {
		str := fmt.Sprint(scores[k].Time)
		array_str := []byte(str)
		date := string(array_str[:DATE_PREM])
		tmp := ScoreInfo{
			Date: date,
			Nums: fmt.Sprint(scores[k].IntValue),
		}
		infos = append(infos, tmp)
	}
	data := _gen_result(infos, 0, "get success", "success")
	result, _ := json.Marshal(data)
	w.Write(result)
}

//------------------------------------------------------------ 资源增加记录
func GetSourceAdd(w http.ResponseWriter, r *http.Request) {
	uid, start, end, err := _get_query_field(w, r, "GET")
	if err != 0 {
		return
	}
	type SourceInfo struct {
		Name   string
		Date   string
		Nums   string
		Origin string
	}

	infos := make([]SourceInfo, 0)
	adds := make(map[string][]stats.IntGameInfo, 0)
	query := map[string]string{
		"gold_add_from_buy":  "gold",
		"food_add_from_buy":  "food",
		"popu_add_from_buy":  "popu",
		"pve_get_food":       "food",
		"pve_get_gold":       "gold",
		"nature_get_food":    "food",
		"nature_get_gold":    "gold",
		"task_get_food":      "food",
		"task_get_gold":      "gold",
		"task_get_popu":      "popu",
		"pvp_get_food":       "food",
		"pvp_get_gold":       "gold",
		"gold_add_from_item": "gold",
		"food_add_from_item": "food",
		"popu_add_from_item": "popu",
	}
	for k, _ := range query {
		adds[k] = stats_tbl.GetPlayLog(uid, k, start, end)
	}

	for kval, val := range adds {
		for v := range val {
			str := fmt.Sprint(val[v].Time)
			array_str := []byte(str)
			date := string(array_str[:DATE_PREM])
			tmp := SourceInfo{
				Name:   query[kval],
				Date:   date,
				Nums:   fmt.Sprint(val[v].IntValue),
				Origin: kval,
			}
			infos = append(infos, tmp)
		}
	}
	data := _gen_result(infos, "0", "get success", "success")
	result, _ := json.Marshal(data)
	w.Write(result)
}

//------------------------------------------------------------ 资源消耗记录
func GetSourceConsume(w http.ResponseWriter, r *http.Request) {
	uid, start, end, err := _get_query_field(w, r, "GET")
	if err != 0 {
		return
	}
	type SourceInfo struct {
		Name   string
		Date   string
		Nums   string
		Origin string
	}

	infos := make([]SourceInfo, 0)
	consume := make(map[string][]stats.IntGameInfo, 0)
	query := map[string]string{
		"food_for_build":       "food",
		"food_for_soldier":     "food",
		"food_for_build_level": "food",
		"pvp_lost_food":        "food",
		"pvp_lost_gold":        "gold",
		"gold_for_build":       "gold",
		"gold_for_build_level": "gold",
		"gold_for_search":      "gold",
		"gold_for_strength":    "gold",
		"popu_for_build":       "popu",
		"popu_for_build_level": "popu",
	}
	for k, _ := range query {
		consume[k] = stats_tbl.GetPlayLog(uid, k, start, end)
	}

	for kval, val := range consume {
		for v := range val {
			str := fmt.Sprint(val[v].Time)
			array_str := []byte(str)
			date := string(array_str[:DATE_PREM])
			tmp := SourceInfo{
				Name:   query[kval],
				Date:   date,
				Nums:   fmt.Sprint(val[v].IntValue),
				Origin: kval,
			}
			infos = append(infos, tmp)
		}
	}
	data := _gen_result(infos, "0", "get success", "success")
	result, _ := json.Marshal(data)
	w.Write(result)
}

//-------------------------------------------------------- 贝壳抽奖记录
func GetGaCha(w http.ResponseWriter, r *http.Request) {
	uid, start, end, errno := _get_query_field(w, r, "GET")
	if errno != 0 {
		return
	}
	type GachaInfo struct {
		Date string
		Item string
		Nums string
	}
	infos := make([]GachaInfo, 0)
	gachas := stats_tbl.GetGachaLog(uid, start, end)
	for k := range gachas {
		str := fmt.Sprint(gachas[k].Time)
		array_str := []byte(str)
		date := string(array_str[:DATE_PREM])
		ins := strings.Split(gachas[k].Key, "#")
		item := ins[2]
		tmp := GachaInfo{
			Date: date,
			Nums: fmt.Sprint(gachas[k].IntValue),
			Item: item,
		}
		infos = append(infos, tmp)
	}
	data := _gen_result(infos, 0, "get success", "success")
	result, _ := json.Marshal(data)
	w.Write(result)

}

func StartGM() {
	http.HandleFunc("/mods/gacha", GetGaCha)
	http.HandleFunc("/mods/source_consume", GetSourceConsume)
	http.HandleFunc("/mods/source_add", GetSourceAdd)
	http.HandleFunc("/mods/vip_score", GetVipScore)
	http.HandleFunc("/mods/hero_unlock", GetHeroUnlockInfo)
	http.HandleFunc("/mods/soldier_unlock", GetSoldierUnlockInfo)
	http.ListenAndServe(":7777", nil)
}
