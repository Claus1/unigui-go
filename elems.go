package unigui

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
		Type string 
		Icon string 
	}

	IGui interface{
		name() string
		changed() Handler
	}

	Edit_ struct {
		Gui
		Complete, Update Handler
		Edit             bool
	}

	Image_ struct {
		Gui
		Image   string
		Scroll bool
		Width   int
		Height  int
	}

	Select_ struct {
		Gui
		Options    []string
	}
	
	Tree_ struct {
		Gui		
		Options  map[string]string 
	}

	Signal struct {
		Maker IGui
		Value string
	}

	Popwindow struct {
		Error, Warning, Info, Progress string
		Data, Update         Any
	}
)

func (gui *Gui) name() string{
	return gui.Name
}

func (gui *Gui) changed() Handler{
	return gui.Changed
}

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
	switch len(wh) {
	case 0 :
		w, h = 500, 350	
	case 1 :
		w, h = wh[0], 350	
	default: 
		w, h = wh[0], wh[1]		
	}
	g := &Image_{Gui{Name: name}, image, false, w, h}
	if click == nil {
		g.Changed = func(value Any) Any {
			g.Value = value
			return nil
		}
	}
	return g
}

func Tree(name string, value string, selected Handler, fields map[string]string) *Tree_ {
	t := &Tree_{Gui{Name: name, Value: value, Changed: selected, Type: "tree"}, fields}
	if selected == nil {
		t.Changed = func(value Any) Any {
			t.Value = value
			return nil
		}
	}	
	return t
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
		Scroll             bool
		Width, Height      int		
	}

	Dialog_ struct {
		Name, Text, Type string
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
	return &Dialog_{name, text,"dialog", nil, buttons, callback}
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
	return &Screen_{Blocks: blocks, Type: "screen", Header: AppName}
}
