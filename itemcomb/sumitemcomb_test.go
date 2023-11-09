package itemcomb

import (
	"testing"
	"reflect"
	"fmt"
)
func Test_createStoreCatalog(t *testing.T){
	tests := []struct {
        name   string
        storeconf StoreConfType
        want   []Store
	}{
        {
		 name : "テスト1",
		 storeconf : StoreConfType{
			"静岡本店": {
				{ "boundary": "0<=","postage":"300" },
			},
			"駿河屋": {
				{"boundary": "1000>", "postage":"440"},
				{"boundary": "1000<=:1500>", "postage":"385"},
				{"boundary": "5000>", "postage":"240" },
			},
		 }, 
		 want : []Store{
			{"静岡本店", []StoreOperator{&SingleStoreOperator{0, "<=", 300} }},
			{"駿河屋", []StoreOperator{
				&SingleStoreOperator{1000, ">", 440},
				&RangeStoreOperator{[]int{1000, 1500}, []string{"<=", ">"}, 385},
				&SingleStoreOperator{5000,">", 240},
			}},
		 },
		},
	}
	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
				if sc := createStoreCatalog(tt.storeconf); !equalStoreList(tt.want, sc) {
					t.Errorf("createStoreCatalog() = %v, want %v", sc, tt.want)
				}
        })
	}

	
}
func equalStoreList(w []Store, result []Store) bool{
	wmap := map[string][]StoreOperator{}
	fmt.Println("w len=", len(w))
	for _,v := range w {
		if len(wmap[v.name]) == 0 {
			wmap[v.name] = []StoreOperator{}
		}
		wmap[v.name] = append(wmap[v.name], v.postageopes...)
	}
	fmt.Println("result len=", len(result))
	for _,v := range result {
		if want, ok := wmap[v.name]; !ok{
			fmt.Println("not found wmap key name=",v.name)
			return false
		} else {
			if len(v.postageopes) != len(want) {
				fmt.Printf("length err name=%s, vlen=%d, wlen=%d\n",v.name,len(v.postageopes), len(want))
				return false
			}
			for i, ww := range want {
				so := v.postageopes[i]
				vtype := reflect.ValueOf(so)
				wtype := reflect.ValueOf(ww)
				if vtype.Type() != wtype.Type() {
					fmt.Printf("not equal type v=%v, w=%v\n",vtype.Type(), wtype.Type())
					return false
				}
				fmt.Println("wtype=",wtype.Type())
				if !so.Equal(ww) {
					fmt.Println("not equal object")
					return false
				}
			}
		}
	}
	return true
}
func Test_ParseStoreConf(t *testing.T){
	tests := []struct {
        name   string
        storeconfjson string
        want   StoreConfType
	}{
		{
			"テスト1",
			`{
				"静岡本店": [
					{ "boundary": "0<=","postage":"300" }
				],
				"駿河屋": [
					{"boundary": "1000>", "postage":"440"},
					{"boundary": "1000<=:1500>", "postage":"385"},
					{"boundary": "5000>", "postage":"240" }
				]
			}`,
			StoreConfType{
				"静岡本店": {
					{ "boundary": "0<=","postage":"300" },
				},
				"駿河屋": {
					{"boundary": "1000>", "postage":"440"},
					{"boundary": "1000<=:1500>", "postage":"385"},
					{"boundary": "5000>", "postage":"240" },
				},
			 },
		},
	}
	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
				if r, ok := ParseStoreConf(tt.storeconfjson); !ok || !equalStoreConf(tt.want, r) {
					t.Errorf("ParseStoreConf() = %v, want %v", r, tt.want)
				}
        })
	}
}
func equalStoreConf(want StoreConfType, ret StoreConfType) bool{
	fmt.Printf("ret= %v\n",ret)
	if len(want) != len(ret){
		fmt.Printf("not equal len w=%d, r=%d\n", len(want), len(ret))
		return false
	}
	bKey := "boundary"
	pKey := "postage"
	for k,wml := range want {
		rml, ok := ret[k]
		if !ok {
			fmt.Println("not exist key =", k)
			return false
		}
		for wi, wm := range wml {
			if wm[bKey] != rml[wi][bKey] {
				fmt.Printf("not equal boundary w=%s, r=%s\n", wm[bKey], rml[wi][bKey])
				return false
			}
			if wm[pKey] != rml[wi][pKey] {
				fmt.Printf("not equal postage w=%s, r=%s\n", wm[pKey], rml[wi][pKey])
				return false
			}
		}
	}
	return true
}
func Test_ParseItemList(t *testing.T){
	tests := []struct {
        name   string
        itemliststr string
        want   []SelectItem
	}{
		{
			"テスト1",
			`[
				{"itemname":"itemA", "storename":"storeA", "price":500},
				{"itemname":"itemA", "storename":"storeB", "price":600},
				{"itemname":"itemB", "storename":"storeA", "price":1500}
			]`,
			[]SelectItem{
				{"itemA", "storeA", 500},
				{"itemA", "storeB", 600},
				{"itemB", "storeA", 1500},
			},
		},
	}
	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
				if r, ok := ParseItemList(tt.itemliststr); !ok || !equalSelectItemList(tt.want, r) {
					t.Errorf("ParseItemList() = %v, want %v", r, tt.want)
				}
        })
	}
}
func equalSelectItemList(want []SelectItem, ret []SelectItem) bool {
	if len(want) != len(ret) {
		fmt.Println("not equal length = ", len(want))
		return false
	}
	fmt.Printf("ret= %v\n",ret)
	for i,_ := range want {
		if want[i].Name != ret[i].Name {
			fmt.Printf("not equal [%d].Name w=%s, r=%s\n", i, want[i].Name, ret[i].Name)
			return false
		}
		if want[i].Storename != ret[i].Storename {
			fmt.Printf("not equal [%d].Storename w=%s, r=%s\n", i, want[i].Storename, ret[i].Storename)
			return false
		}
		if want[i].Price != ret[i].Price {
			fmt.Printf("not equal [%d].Price w=%d, r=%d\n", i, want[i].Price, ret[i].Price)
			return false
		}
	}
	return true
}
