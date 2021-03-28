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
		Gui
		Complete, Update Handler
		Edit             bool
	}

	Image_ struct {
		Gui
		Image   string
		Width   int
		Height  int
	}

	Select_ struct {
		Gui
		Options    []string
	}

	arr3str = [][3]string

	optLevel struct {
		s3    [3]string
		level int
	}

	Tree_ struct {
		Gui
		Options    [][2]string
		elems      arr3str
	}

	Signal struct {
		Maker Any
		Value string
	}

	Popwindow struct {
		Error, Warning, Info string
		Data, Update         Any
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

func UploadButton(name string, clicked Handler, icon string) *Gui {
	bp := Button(name, clicked, icon)
	bp.Type = "gallery"
	return bp
}

func CameraButton(name string, clicked Handler, icon string) *Gui {
	bp := Button(name, clicked, icon)
	bp.Type = "camera"
	return bp
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
	return &Edit_{Gui{Name: str, Value: ""}, nil, nil, false}
}

func Edit(name string, value Any, changed Handler) *Edit_ {
	g := &Edit_{Gui{Name: name, Value: value, Changed: changed},nil, nil, true}
	if changed == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}

func Select(name string, value Any, changed Handler, options []string) *Select_ {
	g := &Select_{Gui{Name: name, Value: value, Changed: changed}, options}
	if changed == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}

func List(name string, value Any, changed Handler, options []string) *Select_ {
	s := Select(name, value, changed, options)
	s.Type = list
	return s
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
	g := &Image_{Gui{Name: name}, image, w, h}
	if click == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}

func Tree(name string, value string, selected Handler, fields *map[string]string) *Tree_ {
	t := &Tree_{Gui{Name: name, Value: value, Changed: selected}, nil, nil}
	if selected == nil {
		t.Changed = func(value Any) Any {
			t.Value = value
			return nil
		}
	}
	t.SetMapFields(fields)
	return t
}

func (s *Tree_) SetFields(elems arr3str) {
	s.elems = elems
	s.Type = list
	filterSoft := func(par2equal string) arr3str {
		var arr arr3str
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
	elems := make(arr3str, len(*fields))
	i := 0
	for k, v := range *fields {
		elems[i] = [3]string{k, k, v}
		i++
	}
	s.SetFields(elems)
}

type (
	CellHandler = func(cellValue TableCell) Any

	TableCell struct {
		Value Any
		Where [2]int
	}

	Table_ struct {
		Gui
		View string
		Headers          []string
		Rows             [][]Any

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
	t := &Table_{Gui :Gui{Name: name, Value: value}, Headers: headers, Rows: rows}
	if selected == nil {
		t.Changed = func(value Any) Any {
			t.Value = value
			return nil
		}
	} else {
		t.Changed = selected
	}
	t.Tools = true
	t.Modify = func(cellValue TableCell) Any{
		AcceptRowValue(t, &cellValue)
		return nil
	}
	return t
}

type (
	Block_ struct {
		Name, Icon         string
		Logo               *Image_
		Top_childs, Childs []Any
		Scrool             bool
		Width, Height      int
		Dispatch           Handler
	}

	Dialog_ struct {
		Name, Text string
		Content    *Block_
		Buttons    []string
		Callback   Handler
	}
	Screen_ struct {
		Name, Icon, Header, Type string
		Blocks             []Any
		Order              int
		Prepare, Save      func()
		Toolbar            []*Gui
		Dispatch           Handler
		handlers           []elemHandle
	}
)

type elemHandle struct {
	elemName string
	blockName string
	nameFunc string
	handler  Handler
}

func Dialog(name string, text string, callback Handler, buttons ...string) *Dialog_ {
	return &Dialog_{name, text, nil, buttons, callback}
}

func (s *Screen_) Handle(elemName string, blockName string, nameFunc string, handler Handler) {
	s.handlers = append(s.handlers, elemHandle{elemName, blockName, nameFunc, handler})
}

func Block(name string, top_childs []Any, childs ...Any) *Block_ {
	if top_childs == nil {
		top_childs = make([]Any, 0)
	}
	if childs == nil {
		childs = make([]Any, 0)
	}
	return &Block_{Name: name, Top_childs: top_childs, Childs: childs}
}

func Screen(blocks ...Any) *Screen_ {
	return &Screen_{Blocks: blocks, Type: "Screen"}
}
