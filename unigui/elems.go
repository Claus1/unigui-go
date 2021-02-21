package unigui

import (
	"sort"
	"strings"
)

const (
	toggles  = "toggles"
	list     = "list"
	dropdown = "dropdown"
)

type (
	Any = interface{}

	Handler = func(value Any) Any

	Gui struct {
		Name       string
		Value      Any
		Changed    Handler
		Type, Icon string
	}
	Edit_ struct {
		Name             string
		Value            Any
		Changed          Handler
		Type, Icon       string
		Complete, Update Handler
		Edit             bool
	}
	Image_ struct {
		Name    string
		Value   Any
		Changed Handler
		Type    string
		Image   string
		Width   int
		Height  int
	}
	Select_ struct {
		Name       string
		Value      Any
		Changed    Handler
		Type, Icon string
		Options    []string
	}
	arrOf3s = [][3]string

	optLevel struct {
		s3    [3]string
		level int
	}
	Tree_ struct {
		Name       string
		Value      Any
		Changed    Handler
		Type, Icon string
		Options    [][2]string
		elems      arrOf3s
	}
	Signal struct {
		Maker Any		
		Value string
	}

	Popwindow struct {
		Error, Warning, Info string
		Data                 Any
	}
)

func Button(name string, clicked Handler, icon string) *Gui {
	b := &Gui{name, nil, clicked, "", icon}
	if clicked == nil {
		b.Changed = func(value Any) Any {
			return nil
		}
	}
	return b
}
func Switch(name string, value bool, changed Handler) *Gui {
	g := &Gui{name, value, changed, "", ""}
	if changed == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}

func Text(str string) *Edit_ {
	return &Edit_{Name: str, Edit: true}
}
func Edit(name string, value Any, changed Handler) *Edit_ {
	g := &Edit_{Name: name, Value: value, Changed: changed}
	if changed == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}

func Image(name string, image string, click Handler, wh ...int) *Image_ {
	var w, h int
	if len(wh) == 0 {
		w, h = 500, 350
	} else if len(wh) == 1 {
		w, h = wh[0], 350
	} else if len(wh) == 2 {
		w, h = wh[0], wh[1]
	}
	g := &Image_{Name: name, Image: image, Width: w, Height: h}
	if click == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}
func Tree(name string, value string, selected Handler, fields *map[string]string) *Tree_ {
	t := &Tree_{Name: name, Value: value, Changed: selected}
	if selected == nil {
		t.Changed = func(value Any) Any {
			t.Value = value
			return nil
		}
	}
	t.SetMapFields(fields)
	return t
}
func (s *Tree_) SetFields(elems arrOf3s) {
	s.elems = elems
	s.Type = list
	filterSoft := func(par2equal string) arrOf3s {
		var arr arrOf3s
		for _, e := range elems {
			if e[2] == par2equal {
				arr = append(arr, e)
			}
		}
		sort.SliceStable(arr, func(i, j int) bool {
			return arr[i][0] < arr[j][0]
		})
		return arr
	}
	var options []optLevel
	//declare recursive
	var make4root func(optLevel)
	make4root = func(ol optLevel) {
		options = append(options, ol)
		childs := filterSoft(ol.s3[1])
		for _, e := range childs {
			make4root(optLevel{e, ol.level + 1})
		}
	}
	roots := filterSoft("")
	for _, e := range roots {
		make4root(optLevel{e, 0})
	}
	s.Options = nil
	for _, e := range options {
		vlines := e.level - 1
		str := ""
		if vlines >= 0 {
			str = strings.Repeat("|", vlines)
			str += "\\"
		}
		str += e.s3[0]
		s.Options = append(s.Options, [2]string{str, e.s3[1]})
	}
}
func (s *Tree_) SetMapFields(fields *map[string]string) {
	elems := make(arrOf3s, len(*fields))
	i := 0
	for k, v := range *fields {
		elems[i] = [3]string{k, k, v}
		i++
	}
	s.SetFields(elems)
}

type (	
	CellHandler  = func(cellValue TableCell) Any

	TableCell struct {
		Value Any
		Where [2]int
	}

	Table_ struct {
		Name       string
		Value      Any
		Changed    Handler
		Type, Icon,View string
		Headers    []string
		Rows       [][]Any

		//Update the user pressed Enter in the field
		Complete, Modify, Update CellHandler
		Append, Delete           Handler //row
		Multimode                bool
		Edit                     bool
		//called when edit mode is changed by the user
		Editing Handler
		//setting to true causes switching to the page with selected row
		Show bool
		//if fasle the table tools are not visible
		Tools bool
	}
)

func AcceptRowValue(t *Table_, tc *TableCell) {
	t.Rows[tc.Where[0]][tc.Where[1]] = tc.Value
}

func Table(name string, value Any, selected Handler, headers []string, rows [][]Any) *Table_ {
	t := &Table_{Name: name, Value: value, Headers: headers, Rows: rows}
	if selected == nil {
		t.Changed = func(value Any) Any {
			t.Value = value
			return nil
		}
	}
	t.Tools = true
	return t
}

type (
	Block_ struct {
		Name               string
		Top_childs, Childs []Any
		Scrool             bool
		Width, Height      int
	}	

	Dialog struct {
		Name, Text string
		Content    *Block_
		Buttons    []string
		Callback   Handler
	}
	Screen_ struct {
		Name, Icon, Header string
		Blocks             []Any
		Order              int
		Prepare, Save      func()
		Toolbar            []*Gui
		handlers           map[string]Handler
	}
)

func elemHandle(nameElem string, nameFunc string) string{
	return F("%s@%s",nameElem, nameFunc)
}

func (s *Screen_) Handle(nameElem string, nameFunc string, handler Handler) {
	s.handlers[elemHandle(nameElem, nameFunc)] = handler
}

func Block(name string, top_childs []Any, childs ...Any) *Block_ {
	return &Block_{Name: name, Top_childs: top_childs, Childs: childs}
}

func Screen(blocks ...Any) *Screen_ {
	return &Screen_{Blocks: blocks, handlers: make(map[string]Handler)}
}
