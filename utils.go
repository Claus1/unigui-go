package unigui

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var (
	F       = fmt.Sprintf
	flatten func(v []Any) []Any
)

func init() {
	flatten = func(v []Any) []Any {
		var res []Any
		for _, e := range v {
			arr, ok := e.([]Any)
			if ok {
				res = append(res, flatten(arr)...)
			} else {
				res = append(res, e)
			}
		}
		return res
	}
}

type Answer struct {
	Answer Any `json:"answer"`
	Param  Any `json:"param"`
	Id     int `json:"id"`
}

func ToJson(o Any) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return b
}

func Fname2url(fn string) string {
	return F("%s/%s", ResourcePort, strings.ReplaceAll(fn, " ", "%20"))
}

func Url2fname(url string) string {
	return strings.ReplaceAll(url[strings.Index(url, "/")+1:], "%20", " ")
}

func Seq(arr ...Any) []Any {
	return arr
}

func SeqSeq(arr ...[]Any) [][]Any {
	return arr
}

func UpdateError(elem2update Any, str string) *Popwindow {
	return &Popwindow{Error: str, Data: elem2update}
}

func Error(str string) *Popwindow {
	return &Popwindow{Error: str}
}

func Warning(str string) *Popwindow {
	return &Popwindow{Warning: str}
}

func Info(str string) *Popwindow {
	return &Popwindow{Info: str}
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func filterGui(name string, value Any) bool {
	return true
}

func filterWidthHeight(name string, value Any) bool {
	return value != nil && !(name == "width" && value == 0) && !(name == "height" && value == 0)
}

func filterTable(name string, value Any) bool {
	return value != ""
}

func filterPopwindow(name string, value Any) bool {
	return value != nil && value != ""
}

func filterUpdater(name string, value Any) bool {
	return value != nil && value != false
}

func (s Updater) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterUpdater)
}
func (s Gui) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterGui)
}

func (s Edit_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterGui)
}

func (s Select_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterGui)
}

func (s Image_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterWidthHeight)
}

func (s Tree_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterGui)
}

func (s Block_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterWidthHeight)
}

func (s Dialog_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterGui)
}

func (s Table_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterTable)
}

func (s Screen_) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterGui)
}

func (s Popwindow) MarshalJSON() ([]byte, error) {
	return serialize(&s, filterPopwindow)
}

func getFieldValue(s Any, fname string) (Any, bool) {

	pe := reflect.ValueOf(s).Elem()
	e := reflect.Indirect(pe)
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		if name == fname {
			return e.Field(i).Interface(), true
		}
	}
	return nil, false
}

func ToInt(val Any) int {
	return int(val.(float64))
}

func any2cellVal(v Any) *TableCell {
	arr := v.([]Any)
	ints := arr[1].([]Any)
	return &TableCell{arr[0], [2]int{ToInt(ints[0]), ToInt(ints[1])}}
}

func serialize(s Any, filter func(string, Any) bool) ([]byte, error) {
	e := reflect.ValueOf(s).Elem()
	mp := map[string]Any{}

	for i := 0; i < e.NumField(); i++ {
		rawname := e.Type().Field(i).Name
		if rawname[0] <= 'Z' {
			name := strings.ToLower(rawname)
			value := e.Field(i).Interface()
			//only for public fields
			if !(value == "" && (name == "type" || name == "icon")) && filter(name, value) {
				vr := reflect.ValueOf(value).Kind()
				if vr != reflect.Func {
					mp[name] = value
				} else if !reflect.ValueOf(value).IsNil() {
					mp[name] = nil
				}
			}
		}
	}
	return json.Marshal(mp)
}
