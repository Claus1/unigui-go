# unigui-go
Universal GUI and App Browser for Go

### Purpose ###
Provide a programming technology that does not require front-end programming, for a server written in any language, for displaying on any device, in any resolution, without any tuning. 

### Import ###
```
import . "unigui"
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
func handler_x( value_x interface{}) interface{}
```
where value_x is a value for the event.

All Gui objects except Button have a field ‘Value’. 
For an edit field the value is a string or number, for a switch or check button the value is boolean, for table is row id or index, e.t.c.
When a user changes the value of the Gui object or presses Button, the server calls the ‘Changed’ function handler.

```
cleanTable := func(v Any) Any{
	table.Rows = SeqSeq()
	return nil
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
‘Changed’ handlers have to return Gui object or array(Seq) of Gui objects that were changed by the handler and Unigui has to redraw or does nothing if all visible elements have the same state. Unigui will do all other jobs for synchronizing automatically. If a Gui object doesn't have 'Changed' handler the object accepts incoming value automatically to the 'Value' variable of gui object.

If 'Value' is not acceptable instead of returning an object possible to return Info, Error or Warning or UpdateError pop-up window. The last function has a object parameter, which has to be synchronized simultaneously with informing about the Error.

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
#### They are called before the inner element handler call. They cancel the call of inner element handler but you can call it as shown below.
For example above interception of select_mode changed event will be:
```
screen.Handle(selector, "Changed", func(v Any) Any{
    if v == "Based"{
        return UpdateError(selector, "Select can not be Based!")
    }
    return nil
})	
```

#### Layout of blocks. #### 
If the blocks are simply listed Unigui draws them from left to right or from top to bottom depending on the orientation setting. If a different layout is needed, it can be set according to the following rule: if the vertical area must contain more than one block, then the enumeration in the array will arrange the elements vertically one after another. If such an element enumeration is an array of blocks, then they will be drawn horizontally in the corresponding area.

#### Example ####
screen = Screen(Seq(b1,b2), Seq(b3, Seq(b4, b5)))
#[b1,b2] - the first vertical area, [b3, [b4, b5]] - the second one.

![alt text](https://github.com/Claus1/unigui/blob/main/tests/multiscreen.png?raw=true)

### Basic gui elements ###
Class names are used only for programmer convenience and do not used by Unigui.
#### If the element name starts from _ , Unigui will not show its name on the screen. ####
if we need to paint an icon somewhere in the element, set the element 'Icon' to 'any MD icon name'.


Common form for element constructors:
```
Gui(name string, value Any, changed Handler)
```
calling the function by default record value to Value field.
Changed(value) 

#### Button ####
Normal button.
```
Button("Push me", push_callback, "") 
```
Icon button 
```
Button("_Check", push_callback, "check") //icon == "check" in MD icon list
```
#### Load to server Button ####
Special button provides file loading from user device or computer to the Unigui server.
```
UploadButton("Load", handler_when_loading_finish)
```
handler_when_loading_finish(button_, the_loaded_file_filename) where the_loaded_file_filename is a file name in upload server folder. This folder name is optional UploadDir parameter in unigui.start.

#### Camera Button ####
Special button provides to make a photo on the user mobile device. 
```
CameraButton('Make a photo', handler_when_shooting_finish)
```
handler_when_loading_finish(button_, name_of_loaded_file) where name_of_loaded_file is the made photo name in the server folder. This folder name is an
var UploadDir in unigui.

#### Edit and Text field. ####
```
Edit('Some field', '') #for string value
Edit('Number field', 0.9) #for numbers
```
If set edit = false it will be readonly field or text label.
```
Edit('Some field', '', edit = false) 
#is equal to
Text('Some field')
```
complete handler is optional function which accepts the current field value and returns a string list for autocomplete.
```
Edit('Edit me', value = '', complete = get_complete_list) #value has to be string or number

def get_complete_list(gui_element, current_value):
    return [s for s in vocab if current_value in s]    
```
Can contain optional 'update' handler which is called when the user press Enter in the field.
It can return None or objects for updating as usual handler.


#### Radio button ####
```
Switch('Radio button', value = True[,changed = ..]) #value has to be boolean, changed - optional
```

#### Select group. Contains options field. ####
```
Select('Select something', "choice1", selection_is_changed, options = ["choice1","choice2", "choice3"]) 
```
can be such type 'toggles','list','dropdown'. Unigui automatically chooses between toogles and dropdown,
but the user can set type = 'list' then Unigui build it as vertical select list.

#### Image. #### 
width,changed and height are optional, changed is called if the user click or touch the image.
```
Image("Image", image = "some url", changed = show_image_info, width = .., height = ..)
or short version
Image("Image", "some url", show_image_info, width = .., height = ..)

```

#### Tree. The element for tree-like data. ####
```
Tree(name, selected_item_key, changed_handler, [unique_elems = .., elems = ..])
```
unique_elems for the data without repeating names, it is dictionary {item_name:parent_name}. If 'unique_elems' defined then 'elems' is redundant.
'elems' for data which can contain repeating names. it is array of arrays [item_name,item_key,parent_key].
parent_name and parent_key are None for root items. changed_handler gets the tree object and item key as value which is the item name for unique items. 

### Table. ###
Tables is common structure for presenting 2D data and charts. Can contain append, delete, update handlers, multimode parameter is True if allowed single and multi select mode. True by default. All of them are optional. When you add a handler for such action Unigui will draw an appropriate action icon button in the table header automatically.
```
table = Table('Videos', [0], row_changed, headers = ['Video', 'Duration', 'Owner', 'Status'],  
  rows = [
    ['opt_sync1_3_0.mp4', '30 seconds', 'Admin', 'Processed'],
    ['opt_sync1_3_0.mp4', '37 seconds', 'Admin', 'Processed']
  ], 
  multimode = false, update = update)
```
If 'headers' length is equal 'rows' length Unigui counts rows id as an index in rows array.
If 'rows' length is 'headers' length + 1, Unigui counts rows id as the last row field.
So it is possible to use some keys as row ids just by adding it to the row as the last element.
If table does not contain append, delete arguments, then it will be drawn without add and remove icons.  
value = [0] means 0 row is selected in multiselect mode (in array). multimode is False so switch icon for single select mode will be not drawn and switching to single select mode is not allowed.


By default Table has toolbar with search field and icon action buttons. It is possible to hide it if set 'tools = False' in the Table constructor.

Table shows a paginator if all rows can not be drawn on the screen. Otherwise a table paginator is redundant and omitted.

If the selected row is not on the currently visible page then setting 'show = True' table parameter causes Unigui to switch to the page with the selected row. 

### Table handlers. ###
complete, modify and update have the same format as the others elements, but value is consisted from the cell value and its position in the table.
'update' is called when the user presses the Enter, 'modify' when the cell value is changed by the user. By default it has standart modify method which updates rows data, it can be locked by
setting 'edit = False' in Table constructor.
They can return Error or Warning if the value is not accepted, othewise the handler has to call accept_rowvalue(table, value) for accepting the value.
```
def table_updated(table_, tabval):
    value, position = tabval
    #check value
    ...
    if error_found:
        return Error('Can not accept the value!')
    accept_rowvalue(_, value)
```
The 'changed' table handler accept the selected row number or id as a value.

'editing' handler is called when the user switches the table edit mode. it is optional and has signature editing(table, edit_mode_now) where the second parameter says the table is being edited or not.

### Chart ###
Chart is a table with additional Table constructor parameter 'view' which explaines unigui how to draw a chart. The format is '{x index}-{y index1},{y index2}[,..]'. '0-1,2,3' means that x axis values will be taken from 0 column, and y values from 1,2,3 columns of row data.
'i-3,5' means that x axis values will be equal the row indexes in rows, and y values from 3,5 columns of rows data. If a table constructor got view = '..' parameter then unigui displays a chart icon at the table header, pushing it switches table mode to the chart mode. If a table constructor got type = 'view' in addition to view parameter the table will be displayed as chart on start. In the chart mode pushing the icon button on the top right switches back to table row mode.

### Signals ###
Unigui supports a dedicated signal event handling mechanism. They are useful in table fields and shared blocks when the containing blocks and screens must respond to their elements without program linking. If a string in a table field started from @ then it considered as a signal. If the user clicks such field in non-edit mode then Unigui generates a signal event, which comes to dispatch function of its containters. First Unigui look at the element block, if not found than at the screen, if not found User.dispatch will be called, which can be redefined for such cases. Any handler can return Signal(element_that_generated_the_event, '@the_event_value') which will be processed.


### Dialog ###
```
Dialog(name, text, callback, buttons, content = None)
```
where buttons is a list of the dialog buttons like ['Yes','No', 'Cancel'].
Dialog callback has the signature as other with value = pushed button name
```
def dicallback(current_dialog, bname):
    if bname == 'Yes':
        do_this()
    elif ..
```
content can be filled by any Gui element sequence for additional dialog functionality.

### Popup windows ###
They are intended for non-blocking displaying of error messages and informing about some events, for example, incorrect user input and the completion of a long process on the server.
```
Info(info_message)
Warning(warning_message)
Error(error_message)
UpdateError(updated_element, error_nessage)
```
They are returned by handlers and cause appearing on the top screen colored rectangles window for 3 second. UpdateError also says Unigui to update changed updated_element.

### Other subtle benefits of a Unigui protocol and technology. ###
1. Works with any set of resource process servers as a single system, within the same GUI user space, carries out any available operations, including cross, on the fly, without programming.
2. Reproduces and saves sequences of the user interaction with the system without programming. It can be used for complex testing, supporting of security protocols and more.
3. Saves and restores the state of the unigui session of the user. Mirrors a session to other users, works simultaneously in one session for many users. 


### Milti-user programming? You don't need it! ###
Unigui automatically creates and serves an environment for every user.
The management class is User which contains all required methods for processing and handling the user activity. A programmer can redefine methods in the inherited class, point it as system user class and that is all. Such methods suit for history navigation, undo/redo and initial operations. The screen folder contains screens which are recreated for every user. The same about blocks. The code and modules outside that folders are common for all users as usual. By default Unigui use the system User class and you do not need to point it. If we need special user class logic, we can define own inheritor User.
```
class Hello_user(unigui.User):
    def __init__(self):
        super().__init__()
        print('New Hello user connected and created!')
    def dispatch(self, elem, ref):
        if http_link(ref[1:]):
            open_inbrowser()
        else:
            return Warning(f'What to do with {ref}?') 

unigui.start('Hello app', user_type = Hello_user)
```
In screens and blocks sources we can access the user by call get_user()
```
user = get_user()
print(isinstance(user, Hello_user))
```

More info about User class methods you can find in manager.py in the root dir.

Examples are in tests folder.

The articles about Unigui and its protocol in details:

in English https://docs.google.com/document/d/1G_9Ejt9ETDoXpTCD3YkR8CW508Idk9BaMlD72tlx8bc/edit?usp=sharing

in Russian https://docs.google.com/document/d/1EleilkEX-m5XOZK5S9WytIGpImAzOhz7kW3EUaeow7Q/edit?usp=sharing

