package itemcomb

import (
	"encoding/json"
	"regexp"
	"fmt"
	"strconv"
)
const (
	MarginPrice = 250
)
type StoreConfType map[string][]map[string]string

type MakeCombResult struct{
	LowestSum  SumItemResult
}
func (m *MakeCombResult) GetResult() SumItemResult {
	return m.LowestSum
}

func ParseStoreConf(jsontext string) (StoreConfType, bool) {
	var data StoreConfType
    if err := json.Unmarshal([]byte(jsontext), &data); err == nil {
        return data, true
    }
    return data, false
}
func ParseItemList(jsontext string) ([]SelectItem, bool) {
	var data []SelectItem
    if err := json.Unmarshal([]byte(jsontext), &data); err == nil {
        return data, true
    }
    return data, false
}
func createStoreCatalog(storeconf map[string][]map[string]string) []Store{
	storecatalog := []Store{}
	pattern := `([0-9]*)([<|>]=?)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return storecatalog
	}
	for k,v := range storeconf {
		sobj := Store {
			name : k,
			postageopes : []StoreOperator{},
		}
		for _, boundary := range v {
			bv := re.FindAllStringSubmatch(boundary["boundary"], -1)
			var ope StoreOperator
			switch len(bv) {
			case 1 :
				bint, _ := strconv.Atoi(bv[0][1])
				pos, _ := strconv.Atoi(boundary["postage"])
				ope = &SingleStoreOperator {
					bval : bint,
					ope_str : bv[0][2],
					postage : pos,
				}
			case 2 :
				bint1, _ := strconv.Atoi(bv[0][1])
				bint2, _ := strconv.Atoi(bv[1][1])
				pos, _ := strconv.Atoi(boundary["postage"])
				ope = &RangeStoreOperator {
					bvals : []int{bint1, bint2},
					ope_strs : []string{bv[0][2], bv[1][2]},
					postage : pos,
				}
			default:
				fmt.Println("Error not support format : ", err)
				continue
			}
			sobj.AddTerms(&ope)
		}
		storecatalog = append(storecatalog, sobj)
	}
	return storecatalog
}
func createItemPtn(itemlist []SelectItem) [][]int {
	ngrp := map[string][]int{}
	for i,item := range itemlist {
		if _,ok := ngrp[item.Name]; ok {
			ngrp[item.Name] = append(ngrp[item.Name], i)
		} else {
			ngrp[item.Name] = []int{i}
		}
	}
	r := [][]int{}
	cnt := 0
	for _, v := range ngrp {
		r = append(r, []int{})
		r[cnt] = append(r[cnt], v...)
		cnt += 1
	}
	return r
}
func getStoreInCatalog(stca []Store, name string) (Store, bool){
	for _, s := range stca {
		if s.name == name {
			return s, true
		}
	}
	return Store{}, false
}
func removeHighPriceItem(stca []Store,
						 itemptn [][]int,
						 itemlist []SelectItem,
						) [][]int {
	type minValueData struct {
		midx int
		mval int
	}
	posin_mingrp := []minValueData{}
    for _, iptn := range itemptn {
		minv := -1
		minidx := -1
		for _, i := range iptn {
			s,ok := getStoreInCatalog(stca, itemlist[i].Storename);
			if !ok {
				continue
			}
			sump := itemlist[i].Price + s.GetPostage(itemlist[i].Price)
			if minv == -1 || minv > sump {
				minv = sump
				minidx = i
			}
		}
		posin_mingrp = append(posin_mingrp, minValueData{minidx, minv})
	}
	newptn := [][]int{}
	for i, _ := range itemptn {
		ary := []int{}
		for _, ptn := range itemptn[i] {
			if posin_mingrp[i].midx == ptn {
				ary = append(ary, ptn)
				continue
			}
			if posin_mingrp[i].mval + MarginPrice >= itemlist[ptn].Price {
				ary = append(ary, ptn)
			}
		}
		newptn = append(newptn, ary)
	}
	return newptn
}
func arys_comb(ary1 [][]int, ary2 []int) ([][]int){
	res_len := len(ary1) * len(ary2)
	result := make([][]int, res_len)
	count := 0
	for _, v := range ary1{
		for _, v2 := range ary2{
			tmp := make([]int, len(v)+1)
			copy(tmp, v)
			tmp[len(v)] = v2
			result[count] = tmp
			count += 1
		}
	}
	return result
}
func ary_comb(ary1 []int, ary2 []int) ([][]int){
	res_len := len(ary1) * len(ary2)
	result := make([][]int, res_len)
	count := 0
	for _, v := range ary1{
		for _, v2 := range ary2{
			tmp := make([]int, 2)
			tmp[0] = v
			tmp[1] = v2
			//fmt.Printf("index i=%d, j=%d\n",i, j)
			result[count] = tmp
			count += 1
		}
	}
	return result
}

func makeComb(itemlist [][]int) ([][]int){
	//res_len := get_result_length(itemlist)
	//fmt.Println("result_len=",res_len)
	//result := make([][]int, res_len)
	if len(itemlist) == 0 {
		return make([][]int,0, 0)
	}
	if len(itemlist) == 1 {
		return itemlist
	}
	cur_res := ary_comb(itemlist[0], itemlist[1])
	if len(itemlist) == 2 {
		return cur_res
	}
	for _, v := range itemlist[2:]{
		cur_res = arys_comb(cur_res, v)
	}
	return cur_res
}
func createBulkBuy(itemlist []SelectItem,
					storeconf map[string][]map[string]string,
					) []SumItem {
	bulk := []SumItem{}
	stca := createStoreCatalog(storeconf)
	itemptn := createItemPtn(itemlist)
	argary := removeHighPriceItem(stca, itemptn, itemlist)
	mc := makeComb(argary)
	for _, comb := range mc {
		si := SumItem{
			items : []SelectItem{},
			stores : []Store{},
			storecatalog : stca,
			sums : map[string]map[string]int{},
		}
		for _, ind := range comb {
			si.AddItem(&itemlist[ind])
		}
		bulk = append(bulk, si)
	}
	return bulk
}
func SaitekiPrice(itemlist []SelectItem,
				  storeconf StoreConfType,
				  ) MakeCombResult {
	result := MakeCombResult{}
	bulk := createBulkBuy(itemlist, storeconf)
	var cheapest SumItem
	bestprice := -1
	for _, b := range bulk {
		sumprice := b.GetSum()
		if bestprice == -1 || bestprice > sumprice {
			cheapest = b
			bestprice = sumprice
		}
	}
	result.LowestSum = cheapest.GetResult()
	return result
}

func getTextResult(outf string, res MakeCombResult) string{
	if outf == "json" {
		j, err := json.Marshal(res.GetResult())
		if err != nil {
			return ""
		}
		return string(j)
	}
	return ""
}
func SearchComb(storeconf StoreConfType,
				itemlist []SelectItem,
				outf  string,
				) string{
	res := SaitekiPrice(itemlist, storeconf)
	return getTextResult(outf, res)
}