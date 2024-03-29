package unigui

import (
	"sort"
)

type (
	Oper interface {
		do()
	}
	screenGen = func(*User) *Screen_
	blockGen  = func(*User) Any

	scrInfo struct {
		name  string
		order int
		icon  string
		gen   screenGen
	}

	User struct {
		UndoBuffer, RedoBuffer, History []Oper
		HistoryCurrent                  int
		activeDialog                    *Dialog_
		screen                          *Screen_
		screens                         map[string]*Screen_
		sharedBlocks                    map[string]Any
		Toolbar                         []*Gui		
		Save, Back, Forward, Undo, Redo func(*User) Any
		Extension                       map[string]Any
		client                          *Client  
	}
	menuItem = [3]Any
)

var (
	users     = make(map[string]User)
	genBlocks = make(map[string]blockGen)
	screens   = make(map[string]scrInfo)

	menu          []menuItem
	sign2funcName = map[string]string{
		"=": "Changed", "->": "Update", "?": "Complete", "+": "Append",
		"-": "Delete", "!": "Editing", "#": "Modify", "$": "Params"}
)

func Register(sgen screenGen, scrName string, order int, icon string) {
	_, found := screens[scrName]
	if found {
		panic(F("Dublicated screen name is found: %s", scrName))
	}
	screens[scrName] = scrInfo{scrName, order, icon, sgen}

	menu = append(menu, menuItem{scrName, icon, order})
	sort.SliceStable(menu, func(i, j int) bool {
		return menu[i][2].(int) < menu[j][2].(int)
	})
}
func ShareBlock(bg blockGen, name string) {
	_, ok := genBlocks[name]
	if ok {
		panic(F("Shared block with repeated name %s!", name))
	}
	genBlocks[name] = bg
}
func call(u *User, f func(*User) Any) Handler {
	return func(value Any) Any {
		if f == nil {
			return nil
		}
		return f(u)
	}
}
func (user *User) Init(cl * Client) {
	user.Toolbar = []*Gui{
		Button("_Back", call(user, user.Back), "arrow_back"),
		Button("_Forward", call(user, user.Forward), "arrow_forward"),
		Button("_Undo", call(user, user.Undo), "undo"),
		Button("_Redo", call(user, user.Redo), "redo")}

	user.sharedBlocks = map[string]Any{}

	user.screens = map[string]*Screen_{}

	user.SetScreen("") //0 order

	user.client = cl
}

var UserConstructor = func(cl* Client) *User {
	user := User{}
	user.Init(cl)
	return &user
}

func (user *User) SetScreen(name string) bool {
	if name == "" {
		name = menu[0][0].(string)
	}
	if user.screen != nil && user.screen.Name == name {
		return false
	}
	scr, ok := user.screens[name]
	if !ok {
		info := screens[name]
		scr = info.gen(user)
		scr.Name = info.name
		scr.Order = info.order
		scr.Icon = info.icon

		if scr.Toolbar == nil {
			scr.Toolbar = user.Toolbar
		}
		user.screens[name] = scr
	}
	user.screen = scr
	if scr.Prepare != nil {
		scr.Prepare()
	}
	return true
}

func (user *User) SharedBlock(name string) Any {	
	if val, ok := user.sharedBlocks[name]; ok {
		return val
	}
	bl := genBlocks[name](user)
	user.sharedBlocks[name] = bl
	return bl
}
func (u *User) appendChange(op Oper) {
	u.UndoBuffer = append(u.UndoBuffer, op)
	u.RedoBuffer = nil
}

func (u *User) handleMessage(msg []Any) Any {
	var result Any
	if u.activeDialog != nil {
		if msg[0] == "root" && msg[1] == nil {
			u.activeDialog = nil
			return nil
		} else if len(msg) == 2 { //button pressed
			result = u.activeDialog.Callback(msg[1])
			u.activeDialog = nil
		} else {
			el := u.findElement(msg)
			if el != nil {
				result = u.processElement(el, msg)
			}
		}
	} else {
		result = u.processMessage(msg)
	}
	if result != nil {
		if dialog, ok := result.(*Dialog_); ok {
			u.activeDialog = dialog
		}
		result = u.prepareResult(result)
	}
	return result
}

func (u *User) Progress(str string) {
	v := &Popwindow{Progress : str}
	u.client.send <- ToJson(v)
}

func (u *User) processMessage(arr []Any) Any {
	if arr[0] == "root" {
		nameScr := arr[1].(string)
		u.SetScreen(nameScr)
		return u.screen
	}
	elem := u.findElement(arr)
	//recursive for Signals
	for {
		if elem == nil{
			return Error(F("Invalid element name in the client message %v", arr))
		}
		res := u.processElement(elem, arr)
		sig, ok := res.(Signal)
		if !ok {
			return res
		}
		elem = sig.Maker		
		arr = Seq("", elem.name(), "@", sig.Value)
	}
}

func (u *User) processElement(elem IGui, msg []Any) Any {
	id := 0
	if len(msg) == 5 {
		id = ToInt(msg[4])
		msg = msg[:4]
	}
	sign := msg[2].(string)
	if sign == "$" {
		return nil
	}
	name := msg[1].(string)
	val := msg[3]
	blockName := msg[0]
	funcName, ok := sign2funcName[sign]
	var res Any
	if ok {
		for _, eh := range u.screen.handlers {
			if eh.elemName == name && blockName == eh.blockName && eh.nameFunc == funcName {
				hres := eh.handler(val)
				if hres != false {
					return hres
				} else {
					break //causes call the inner element handler
				}
			}
		}
		
		if funcName == "Changed"{
			res = elem.changed()(val)
		} else {
			hi, ok := getFieldValue(elem, funcName)
			if !ok {
				panic(F("%s doesn't contain field %s!", name, funcName))
			}
			if h, ok := hi.(Handler); ok {
				//it's allowed
				if h == nil && funcName == "Editing" {
					return nil
				}
				res = h(val)

			} else if ch, isCH := hi.(CellHandler); isCH {

				tCell := any2cellVal(val)
				res = ch(*tCell)
				if res == false {
					AcceptRowValue(elem.(*Table_), tCell)
					return nil
				}
			} else {
				panic(F("%s.%s has unknown type!", name, funcName))
			}
			if id != 0 {
				res = Answer{res, nil, id}
			}
		}
	} 	
	return res
}

func (u *User) blockElem(elem Any) *Block_ {
	for _, blAny := range flatten(u.screen.Blocks) {
		block := blAny.(*Block_)
		if block == elem {
			return block
		}
		for _, c := range append(block.Top_childs, block.Childs...) {
			sq, ok := c.([]Any)
			if ok {
				for _, e := range sq {
					if e == elem {
						return block
					}
				}
			} else {
				if c == elem {
					return block
				}
			}
		}
	}
	return nil
}

func (u *User) findPath(elem Any) []string {
	block := u.blockElem(elem)
	name := ""
	g, ok := elem.(IGui)	
	if ok{
		name = g.name()
	} else {
		n, _ := getFieldValue(elem, "Name")
		name = n.(string)
	}
	return []string{block.Name, name}
}

func (u *User) prepareResult(val Any) Any {
	if val == true {
		val = u.screen		
	} else {
		switch val.(type) {
		case Answer, *Dialog_:
			{
			}
		case []Any:
			arr := val.([]Any)
			path := []Any{}
			for _, e := range arr {
				path = append(path, u.findPath(e))
			}
			val = Updater{path, val, true}
		case *Popwindow:
			popwin := val.(*Popwindow)
			if popwin.Data != nil {
				popwin.Update = u.findPath(popwin.Data)
			}
		default:
			val = Updater{u.findPath(val), val, false}
		}
	}
	return val
}

func (u *User) findElement(arr []Any) IGui {
	if arr[0] == "toolbar" {
		for _, b := range u.Toolbar {
			if b.Name == arr[1] {
				return b
			}
		}
	}
	for _, blAny := range flatten(u.screen.Blocks) {
		bl := blAny.(*Block_)
		if bl.Name == arr[0] {
			for _, c := range append(bl.Top_childs, bl.Childs...) {
				sq, ok := c.([]Any)
				if ok {
					for _, e := range sq {
						g := e.(IGui)
						if g.name() == arr[1] {
							return g
						}
					}
				} else {
					g := c.(IGui)
					if g.name() == arr[1] {
						return g
					}
				}
			}
		}
	}
	return nil
}
