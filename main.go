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
	storeconf := itemcomb.ParseMapAny(os.Args[1])
	itemlist := itemcomb.ParseItemList(os.Args[2])
	outf := "json"
	result := itemcomb.SearchComb(storeconf, itemlist, outf)
	fmt.Println(result)
	return
}