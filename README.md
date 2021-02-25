# unigui-go
Universal GUI and App Browser for Go

### Purpose ###
Provide a programming technology that does not require front-end programming, for a server written in any language, for displaying on any device, in any resolution, without any tuning. 

### Import ###
```
import . "github.com/Claus1/unigui-go"
```

### How does work inside ###
The exchange protocol for the solution is JSON as the most universally accessible, comprehensible, readable, and popular format compatible with all programming languages.  The server sends JSON data to Unigui which has built-in tools (autodesigner) and automatically builds a standart Google Material Design GUI for user data. No markup, drawing instructions and the other dull job are required. Just the simplest description what you want. From the constructed Unigui screen the server receives a JSON message flow which fully describes what the user did. The message format is ["Block", "Elem", "type of action", value], where "Block"and "Elem"are the names of the block and its element, "value" is the JSON value of the action/event that has happened. The server can either accept the change or roll it back by sending an info window about an inconsistency. The server can open a dialog box, send popup Warning, Error,.. or an entirely new screen. Unigui instantly and automatically displays actual server state. 

### Programming ###
Unigui is the language and platform independent technology. This repo explains how to work with Unigui using Go and the tiny but optimal framework for that.
Unigui web version is included in this library. Unigui for Python is accessible in https://github.com/Claus1/unigui

### High level - Screen ###
The program declares screen builder function which has to be registered.

Screen example examples/HelloUnigui.go
The block example with a table, button and selector
```
package main

import . "unigui"		

func screenTest(user* User)* Screen_{	
	table := Table("Videos",0, nil, []string{"Video", "Duration",  "Links", "Mine"},
	    SeqSeq(Seq("opt_sync1_3_0.mp4", "30 seconds",  "@Refer to signal1", true),
		       Seq("opt_sync1_3_0.mp4", "37 seconds",  "@Refer to signal8", false)))
			
	cleanButton := Button("Clean table", nil, "")
	selector := Select("Select", "All", nil, []string{"All","Based","Group"})
	block := Block("X Block", Seq(cleanButton, selector), table)
	block.Icon = "api"
	return Screen(block)	
}
func main(){			
	//register screens
	Register(screenTest, "Main", 0, "api")	
	Start()
}
```

### Server start ###
Unigui builds the interactive app for the code above.
Connect a browser to localhast:8080 and will see:

![alt text](https://github.com/Claus1/unigui/blob/main/tests/screen1.png?raw=true)

### Handling events ###
All handlers are functions which have a signature
```
func handlerX( valueX interface{}) interface{}
```
where valueX is a value for the event.

All Gui objects have a field ‘Value’. 
For an edit field the value is a string or number, for a switch or check button the value is boolean, for table is row id or int index, e.t.c.
When a user changes the value of the Gui object or presses Button, the server calls the ‘Changed’ function handler.

```
cleanTable := func(v Any) Any{
	table.Rows = SeqSeq() //empty [][]Any
    
	return table          //table is changed, just return it for updating on the screen
}
cleanButton := Button("Clean table", cleanTable, "")
```
where Any, Seq, SeqSeq just short names defined as
```
type Any = interface{}
func Seq(arr ...Any) []Any{
	return arr
}
func SeqSeq(arr ...[]Any) [][]Any{
	return arr
}
```
‘Changed’ handlers have to return Gui object or array(Seq) of Gui objects that were changed by the handler and Unigui has to redraw or does nothing if all visible elements have the same state. Unigui will do all other jobs for synchronizing automatically. If a Gui object doesn"t have "Changed" handler the object accepts incoming value automatically to the "Value" variable of gui object.

If "Value" is not acceptable instead of returning an object possible to return Info, Error or Warning or UpdateError pop-up window. The last function has a object parameter, which has to be synchronized simultaneously with informing about the Error.

```
selector := Select("Select", "All", nil, []string{"All","Based","Group"})

selector.Changed = func(v Any) Any{
    if v == "Based"{
        return UpdateError(selector, "Select can not be Based!")
    }
    return nil
})	
```
#### If a handler returns true the whole screen will be redrawn. Also it causes calling Screen function Prepare() which used for syncronizing GUI elements one to another and with the program/system data. Prepare() is also automatically called when the screen loaded. Prepare() is optional.

If the handler returns nil Unigui considers it as Ok and does nothing.

### Block details ###
The width and height of blocks is calculated automatically depending on their children. It is possible to set the block width and make it scrollable in height, for example for images list. Possible to add MD icon to the header, if required. Width, scroll, .. are optional. Block helper is
```
func Block(name string, top_childs []Any, childs ...Any) *Block_ 
```
 
The top_childs parameter of the Block constructor is an array of widgets which has to be in the header just after the block name.
Blocks can be shared between the user screens with its states. Such a block has to be declared inside Block making function and registered by call ShareBlock(myBlock).
Examples of such block examples/shared_block.go:
```

func sharedAudios(user* User) Any{
	table := Table("Audios",0, nil, []string{"Audio", "Duration,sec", "Stars"}, genRows())
    tableBlock := Block("Table chart, Button("Press me", nil), table)
	tableBlock.Icon = "insights"
    return tableBlock
}
func main(){		
	//register shared blocks
	ShareBlock(sharedAudios, "Audios")
    ...
	//register screens
	...
}
```


If some elements are enumerated inside an array, Unigui will display them on a line, otherwise everyone will be displayed on a new own line(s).
 
Using a shared block in some screen:
```
scr :=  Screen(Seq(block, bottomBlock), user.SharedBlock("Audios"))	//user is always accesible in screen making function

```

#### Events interception of shared blocks ####
Interception handlers have the same in/out format as usual handlers.
#### They overrides inner element handler call. If such method returns false the  inner element handler will be calld. 
For example above interception of select_mode changed event will be:
```
screen.Handle(elemName string, blockName string, handlerName string, func(v Any) Any{
    if v == "Based"{
        return UpdateError(selector, "Select can not be Based!")
    }
    return false //call the inner element handlers.
})	
```

#### Layout of blocks. #### 
If the blocks are simply listed Unigui draws them from left to right or from top to bottom depending on the orientation setting. If a different layout is needed, it can be set according to the following rule: if the vertical area must contain more than one block, then the enumeration in the array will arrange the elements vertically one after another. If such an element enumeration is an array of blocks, then they will be drawn horizontally in the corresponding area.

#### Example ####
screen = Screen(Seq(b1,b2), Seq(b3, Seq(b4, b5)))
#[b1,b2] - the first vertical area, [b3, [b4, b5]] - the second one.

![alt text](https://github.com/Claus1/unigui/blob/main/tests/multiscreen.png?raw=true)

### Basic gui elements ###

#### If the element name starts from _ , Unigui will not show its name on the screen. ####
if we need to paint an icon in the element, set the element "Icon" to "any MD icon name".

Common form for element constructors:
```
gui := Gui(name string, value Any, changed Handler)
```
changed Handler normally is not used directly.

#### Button ####
Normal button.
```
Button(name string,push_callback Handler, icon string) 
```
Icon button : name is started form _ for hiding
```
Button("_Check", push_callback, "check") //icon == "check" in MD icon list
```
#### Load to server Button ####
Special button provides file uploading from user device or computer to the Unigui server.
```
UploadButton(name string, handler_when_loading_finish Handler, icon string)
```
handler_when_loading_finish(the_loaded_file_filename) where the_loaded_file_filename is a file name in upload server folder. This folder name is global UploadDir parameter in unigui which can be changed before Start().

#### Camera Button ####
Special button provides to make a photo on the user mobile device. 
```
CameraButton("Make a photo", handler_when_shooting_finish)
```
handler_when_loading_finish(button_, name_of_loaded_file) where name_of_loaded_file is the made photo name in the server folder. This folder name is global UploadDir parameter in unigui which can be changed before Start().

#### Edit and Text field. ####
```
Edit("Name", "value", changed_handler) #for string value
numEdit := Edit("Number field", 0.9, nil) #for edit number
```
If set Edit = false it will be readonly field or text label.
```
numEdit.Edit = false

//Text field
Text("Some text")
```
Complete handler is optional function which accepts the current field value and returns a string list for autocomplete.
```
edit := Edit("Edit me",  "")
edit.Complete = getCompleteList 

def getCompleteList(current_value Any):
    return []string{"option1","option2","option3"}    
```
Possible to set "Update" handler which is called when the user press Enter in the field.
It can return nil or Gui object(s) for updating as usual handler.


#### Radio button ####
```
Switch("Radio button", value bool, changed Handler) #changed - optional
```

#### Select group. Contains options field. ####
```
//build select field
Select(name String, value Any, selectionChanged Handler, options []string)

//build as vertical list
List(name String, value Any, selectionChanged Handler, options []string)
```
Select can be such type "toggles","dropdown". Unigui automatically chooses between toogles and dropdown,
but the user can set preferrable type then Unigui build it as the user want.

#### Image. #### 
Width, and Height are optional, click is called if the user click or touch the image, can be nil as all Hadlers
```
func Image(name string, image string, click Handler, wh ...int) *Image_

```

#### Tree. The element for tree-like data. ####
```
func Tree(name string, value string, selected Handler, fields *map[string]string) *Tree_
```
fields for the tree data {item_name:parent_name}.

parent_name is "" for root items. selected is called when the user clicks on a tree item. 

### Table. ###
Tables is common structure for presenting 2D data and charts. Can contain Append, Delete, Update handlers, Multimode parameter is True if allowed single and multi select mode. True by default. All of them are optional. When you assign a handler for such action Unigui will draw the appropriate action icon button in the table header automatically.
If Modify and Update are not defined, unigui will not draw Edit button and user can not edit the table data. The same rule for Delete, Append, e.t.c.
```
func Table(name string, value Any, selected Handler, headers []string, rows [][]Any) *Table_
```
If "headers" length is equal "rows" length Unigui counts rows id as an index in rows array.
If "rows" length is "headers" length + 1, Unigui counts rows id as the last row field.
So it is possible to use some keys as row ids just by adding it to the row as the last element.
value == [0] means 0 row is selected in multiselect mode (in array). value = 1 means rows at 1 index is selected in sinlge mode selection.

By default Table has toolbar with search field and icon action buttons. It is possible to hide it if set "Tools" table variable to false.

Table shows a paginator if all rows can not be drawn on the screen. Otherwise a table paginator is redundant and omitted.

If the selected row is not on the currently visible page then setting "show = True" table parameter causes Unigui to switch to the page with the selected row. 

### Table handlers. ###
Complete, Modify and Update are CellHandlers where CellHandler = func(cellValue TableCell) Any .
cellValue is consisted from the cell value and its position in the table.
```
TableCell struct {
		Value Any
		Where [2]int
	}
```
"Update" is called when the user presses the Enter, "Modify" when the cell value is changed by the user. By default it has standart modify method which updates rows data, it can be locked by setting the table variale "Edit" to false.
They can return Error or Warning if the value is not accepted, othewise false for accepting the value (false means continue the standart process).
```
table.Update = func(tvalue TableCell) Any{
    AcceptRowValue(table, &tvalue)
    return Warning(F("%v is updated to %v",table.Name, tvalue.Value))
}
```
The "Changed" table handler accepts the selected row number or id as a value.

"Editing" handler is called when the user switches the table edit mode. it is optional and has standart signature where the parameter says the table is being edited or not.

### Chart ###
Chart is a table with additional Table constructor parameter "View" which explaines unigui how to draw a chart. The format is "{x index}-{y index1},{y index2}[,..]". "0-1,2,3" means that x axis values will be taken from 0 column, and y values from 1,2,3 columns of row data.
"i-3,5" means that x axis values will be equal the row indexes in rows, and y values from 3,5 columns of rows data. If a table constructor got View = ".." parameter then unigui displays a chart icon at the table header, pushing that switches table mode to the chart mode. If a table constructor set Type to "view" in addition to "View" parameter the table will be displayed as chart on start. In the chart mode pushing the icon button on the top right switches back to table row mode.

### Signals ###
Unigui supports a dedicated signal event handling mechanism. They are useful in table fields and shared blocks when the containing blocks and screens must respond to their elements without program linking. If a string in a table field started from @ then it considered as a signal. If the user clicks such field in non-edit mode then Unigui generates a signal event, which pop-up to dispatch functions of its containters. First Unigui look at the element block, if not found than at the screen, if not found User.Dispatch will be called, which can be redefined for such cases. Any handler can return Signal(element_that_generated_the_event, "the_event_value") which will be processed.


### Dialog ###
```
func Dialog(name string, text string, callback Handler, buttons ...string) *Dialog_
```
where buttons is a list of the dialog buttons like ["Yes","No", "Cancel"].
Dialog callback has the signature as other with value == pushed button name
```
func dialogCallback(pressedButton Any) Any{
	return Warning(F("The user pressed the button %v", pressedButton))
}
```
Content dialog field can be filled by any Block for additional dialog functionality.

### Popup windows ###
They are intended for non-blocking displaying of error messages and informing about some events, for example, incorrect user input and the completion of a long process on the server.
```
Info(info_message)
Warning(warning_message)
Error(error_message)
UpdateError(updated_element, error_nessage)
```
They are returned by handlers and cause appearing on the top screen colored rectangles window for 3 second. UpdateError also says Unigui to update updated_element.

### Other subtle benefits of a Unigui protocol and technology. ###
1. Possible to works with any set of resource process servers as a single system, within the same GUI user space, carries out any available operations, including cross, on the fly, without programming.
2. Reproduces and saves sequences of the user interaction with the system without programming. It can be used for complex testing, supporting of security protocols and more.
3. Possible to mirror a session to other users, works simultaneously in one session for many users. 


### Milti-user programming? You don"t need it! ###
Unigui automatically creates and serves an environment for every user.
The management class is User which contains all required methods for processing and handling the user activity. A programmer can assign methods 
```
Dispatch                        func(*User, Signal) Any
Save, Back, Forward, Undo, Redo func(User) Any
//also store and use any data in User.Extention which is defined as 
Extension  map[string]Any
```
Such methods suit for history navigation, undo/redo and initial operations.

For constructing custom User use UserConstuctor variable which return new User for any new session.
```
UserConstuctor = func() *User{
    user := User
    //assign custom function
    ..
    return &user
}

The code and modules outside that folders are common for all users as usual. By default Unigui UserConstuctor creates a user with empty behavior function and fields. ```


In screen and shared block functions User automatically acccesible as a function argument.


More info about User class and methods you can find in manager.go in the root dir.

Examples are in examples folder.

The articles about Unigui and its protocol:

in English https://docs.google.com/document/d/1G_9Ejt9ETDoXpTCD3YkR8CW508Idk9BaMlD72tlx8bc/edit?usp=sharing

in Russian https://docs.google.com/document/d/1EleilkEX-m5XOZK5S9WytIGpImAzOhz7kW3EUaeow7Q/edit?usp=sharing

