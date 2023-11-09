package main

import (
	"fmt"
	"os"
	"github.com/gkjg8787/Go-ItemCombSum/itemcomb"
)

func main(){
	if len(os.Args) != 3 {
		fmt.Println("")
		return
	}
	storeconf, ok := itemcomb.ParseStoreConf(os.Args[1])
	if !ok {
		return
	}
	itemlist, ok := itemcomb.ParseItemList(os.Args[2])
	if !ok {
		return
	}
	outf := "json"
	result := itemcomb.SearchComb(storeconf, itemlist, outf)
	fmt.Println(result)
	return
}