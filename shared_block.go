package main

import (
	. "unigui/unigui"
	"math/rand"
	. "math"	
)

func updated( value Any)Any{
	return &Popwindow{Warning : F("Element is updated to %v", value)}
}

func completeTable(tvalue TableCell) Any{	
	switch s := tvalue.Value.(type){
	case string:
		if len(s) > 2{
			return &[]string{"aaa", "bbbb", "cccccc"}
		}
	}
	return nil
}

func completeEdit(value Any) Any{
	return &[]string{"aaa1", "bbbb2", "cccccc3"}
}

func callDialog(value Any) Any{
	return Dialog{"Dialog", "Answer pls..", nil, []string{"Yes", "No"}, dialogCallback}
}

func dialogCallback(pressedButton Any) Any{
	return &Popwindow{Warning : F("The user pressed the button %v", pressedButton)} 
}

func deleteRow(t* Table_, value Any) Any{	
	index := t.Value.(int)
	t.Rows = append(t.Rows[:index], t.Rows[index+1:]...) 
	return t
}

func genRows() [][]Any{
	rows := [][]Any{}
	for i:= 0; i < 100; i++{
		rows = append(rows, []Any{F("sync%v.mp3", i),
			Floor(rand.Float64() * 15000)/100, rand.Intn(100)})
	}
	return rows
}

func sharedAudios(user* User) Any{

	table := Table("Audios",0, nil, []string{"Audio", "Duration,sec", "Stars"}, genRows())
	table.Complete = completeTable
	table.View = "i-1,2"

	table.Modify = func(tvalue TableCell) Any{
		return Error(F("%v is not modified to %v",table.Name, tvalue.Value))
	}			
	table.Update = func(tvalue TableCell) Any{
		AcceptRowValue(table, &tvalue)
		return Warning(F("%v is updated to %v",table.Name, tvalue.Value))
	}			
	table.Changed =	func(value Any) Any{
		table.Value = value
		return Warning(F("%v is changed to %v",table.Name, value))	
	}		
	treeData := map[string]string{
		"Animals" : "",
		"Brushtail Possum" : "Animals",
		"Genet" : "Animals",
		"Silky Anteater" : "Animals",
		"Greater Glider" : "Animals",
		"Tarsier" : "Animals",
		"Kinkajou" : "Animals",
		"Tree Kangaroo" : "Animals",
		"Sunda Flying Lemur" : "Animals",
		"Green Tree Python" : "Animals",
		"Fruit Bat" : "Animals",
		"Tree Porcupines" : "Animals",
		"Small Tarsier" : "Tarsier",
		"Very small Tarsier": "Small Tarsier"}

	treeSelected := func(value Any) Any{
		return Info(F("%v selected!", value))
	}
	tree := Tree("Inheritance","Animals",treeSelected, &treeData)	
	return Block("New block", Seq(Button("Dialog", callDialog, ""), 
		Edit("Simple Enter update", "cherokkee", updated)), Seq(tree, table))
}
