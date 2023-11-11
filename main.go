package main

import (
	"fmt"
	"os"
    "strings"
	"encoding/json"
	"github.com/gkjg8787/Go-ItemCombSum/itemcomb"
)
func createErrorMsg(errtype string,
					errname string,
					errvalue string,
					) string{
	errmsg := errtype
	if len(errname) != 0 {
		errmsg += "," + errname + "=" + errvalue
	}
	ret := itemcomb.SumItemResult{
		Errormsg : errmsg,
	}
	j, err := json.Marshal(ret)
	if err != nil {
		return ""
	}
	return string(j)
}

func trimParamString(param string) string{
	ret := strings.Replace(param,`\"`, `"`, -1)
	ret = strings.TrimLeft(ret,"\"")
	ret = strings.TrimRight(ret,"\"")
    return ret
}

func main(){
	if len(os.Args) != 3 {
		fmt.Println(createErrorMsg("parameter", "", ""))
		return
	}
    //fmt.Println("storeconf = ", os.Args[1])
    //fmt.Println("itemlist = ", os.Args[2])
    //fmt.Println("")
	rstconf := trimParamString(os.Args[1])
	storeconf, ok := itemcomb.ParseStoreConf(rstconf)
	if !ok {
		fmt.Println(createErrorMsg("encoding", "storeconf", string(rstconf)))
		return
	}
	ril := trimParamString(os.Args[2])
	itemlist, ok := itemcomb.ParseItemList(ril)
	if !ok {
		fmt.Println(createErrorMsg("encoding", "itemlist", string(ril)))
		return
	}
	outf := "json"
	result := itemcomb.SearchComb(storeconf, itemlist, outf)
    /*
    if len(result) == 0 {
        fmt.Println("no result")
        return
    }
    */
	fmt.Println(result)
	return
}
