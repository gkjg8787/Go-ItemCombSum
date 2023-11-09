package itemcomb



type SumItemKey struct {
	Sum  string
	Postage  string
	SumPosIn  string
	SumPosOut  string
	SumPostage  string
	Items  string
	ItemName  string
	Price  string
}
func CreateSumItemKey() SumItemKey{
	return SumItemKey{
		Sum : "sum",
		Postage : "postage",
		SumPosIn : "sumposin",
		SumPosOut : "sumposout",
		SumPostage : "sumpostage",
		Items : "items",
		ItemName : "itemname",
		Price : "price",
	}
}

type SumItem struct {
	items []SelectItem
	stores []Store
	storecatalog []Store
	sums map[string]map[string]int
}
func (s *SumItem) AddItem(si *SelectItem) {
	s.items = append(s.items, *si)
	if !s.existStoreName(si.Storename) {
		if store, ok := s.getStore(si.Storename); ok {
			s.stores = append(s.stores, store)
		}
	}
}
func (s *SumItem) existStoreName(storename string) bool {
	for _, store := range s.stores {
		if store.name == storename {
			return true
		}
	}
	return false
}
func (s *SumItem) getStore(storename string) (Store, bool) {
	for _, store := range s.storecatalog {
		if store.name == storename {
			return store, true
		}
	}
	return Store{}, false
}
func (s *SumItem) CreateSums(){
	s.sums = map[string]map[string]int{}
	SMK := CreateSumItemKey()
	for _, item := range s.items {
		if _, ok := s.sums[item.Storename]; !ok {
			s.sums[item.Storename] = map[string]int{SMK.Sum:0, SMK.Postage:0}
		}
		s.sums[item.Storename][SMK.Sum] += item.Price
	}
	for _, store := range s.stores {
		if len(s.sums[store.name]) == 0 {
			s.sums[store.name] = map[string]int{}
		}
		s.sums[store.name][SMK.Postage] = store.GetPostage(s.sums[store.name][SMK.Sum])
	}
}
func (s *SumItem) GetSum() int{
	if len(s.sums) == 0 {
		s.CreateSums()
	}
	SMK := CreateSumItemKey()
	allsum := 0
	for _, value := range s.sums {
		allsum += value[SMK.Sum] + value[SMK.Postage]
	}
	return allsum
}
func (s *SumItem) GetResult() SumItemResult {
	if len(s.sums) == 0 {
		s.CreateSums()
	}
	SMK := CreateSumItemKey()
	result := SumItemResult{
		Errormsg : "", 
		SumPosIn : s.GetSum(),
		SumPostage : 0,
		StoreSums : []StoreSumResult{},
	}
	sumPostage := 0
	for storename, value := range s.sums {
		sumPostage += value[SMK.Postage]
		ssr := StoreSumResult{
			StoreName : storename,
			Postage : value[SMK.Postage],
			SumPosOut : value[SMK.Sum],
			Items : []ItemResult{},
		}
		for _, item := range s.items {
			if item.Storename != storename {
				continue
			}
			ir := ItemResult{
				ItemName : item.Name,
				Price : item.Price,
			}
			ssr.Items = append(ssr.Items, ir)
		}
		result.StoreSums = append(result.StoreSums, ssr)
	}
	result.SumPostage = sumPostage
	return result
}
type SumItemResult struct {
	Errormsg  string `json:"errormsg"`
	SumPosIn  int `json:"sumposin"`
	SumPostage  int `json:"sumpostage"`
	StoreSums  []StoreSumResult `json:"storesums"`
}
type StoreSumResult struct {
	StoreName string `json:"storename"`
	Postage  int `json:"postage"`
	SumPosOut  int `json:"sumposout"`
	Items  []ItemResult `json:"items"`
}
type ItemResult struct {
	ItemName  string `json:"itemname"`
	Price  int `json:"price"`
}

type Store struct {
	name string
	postageopes []StoreOperator
}
func (s *Store) AddTerms(so *StoreOperator){
	s.postageopes = append(s.postageopes, *so)
}
func (s *Store) GetPostage(price int) int{
	sumPostage := 0
	for _, v := range s.postageopes {
		sumPostage += v.Calc(price)
	}
	return sumPostage
}


type StoreOperator interface {
	Calc(int) int
	Equal(any) bool
}
func InRange(ope_str string, bval int, price int) bool{
	if "<" == ope_str {
		return (bval < price)
	} else if "<=" == ope_str {
		return (bval <= price)
	} else if ">" == ope_str {
		return (bval > price)
	} else if ">=" == ope_str {
		return (bval >= price)
	}
	return false
}
type RangeStoreOperator struct {
	bvals []int
	ope_strs []string
	postage int
}
func (ro *RangeStoreOperator) Calc(price int) int{
	if len(ro.bvals) != 2 ||
		len(ro.ope_strs) != 2{
		return 0
	}
	res1 := InRange(ro.ope_strs[0], ro.bvals[0], price)
	res2 := InRange(ro.ope_strs[1], ro.bvals[1], price)
	if res1 && res2 {
		return ro.postage
	}
	return 0
}
func (ro *RangeStoreOperator) Equal(sp any) bool{
	var t RangeStoreOperator
	switch sp.(type) {
	case RangeStoreOperator :
		t = sp.(RangeStoreOperator)
	case *RangeStoreOperator :
		t = *(sp.(*RangeStoreOperator))
	default :
		return false
	}
	if len(ro.bvals) != len(t.bvals) {
		return false
	}
	if len(ro.ope_strs) != len(t.ope_strs){
		return false
	}
	if len(t.bvals) != len(t.ope_strs){
		return false
	}
	if ro.postage != t.postage {
		return false
	}
	for i, bval := range t.bvals {
		if ro.bvals[i] != bval { return false }
		if ro.ope_strs[i] != t.ope_strs[i] { return false }
	}
	return true
}

type SingleStoreOperator struct {
	bval int
	ope_str string
	postage int
}
func (so *SingleStoreOperator) Calc(price int) int{
	if InRange(so.ope_str, so.bval, price) {
		return so.postage
	}
	return 0
}
func (ro *SingleStoreOperator) Equal(sp any) bool{
	var t SingleStoreOperator
	switch sp.(type) {
	case SingleStoreOperator :
		t = sp.(SingleStoreOperator)
	case *SingleStoreOperator :
		t = *(sp.(*SingleStoreOperator))
	default :
		return false
	}
	if ro.bval != t.bval {
		return false
	}
	if ro.ope_str != t.ope_str {
		return false
	}
	if ro.postage != t.postage {
		return false
	}
	return true
}

type SelectItem struct {
	Name string `json:"itemname"`
	Storename string `json:"storename"`
	Price int	`json:"price"`
}
