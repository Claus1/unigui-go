package main

import (
	. "unigui/unigui"		
)

func screenTest(user* User)* Screen_{	
	table := Table("Videos",0, nil, []string{"Video", "Duration",  "Links", "Mine"},
	SeqSeq(Seq("opt_sync1_3_0.mp4", "30 seconds",  "@Refer to signal1", true),
		Seq("opt_sync1_3_0.mp4", "37 seconds",  "@Refer to signal8", false)))
		
	cleanTable := func(v Any) Any{
		table.Rows = SeqSeq()
		return nil
	}
	cleanButton := Button("Clean table", cleanTable, "")

	selector := Select("Select", "All", nil, []string{"All","Based","Group"})

	listRefs := List("Detail ref list signals", "", nil, []string{"Select reference"})

	blockDispatch := func(value Any) Any{
		for i := 0; i <10; i++{
			listRefs.Options = append(listRefs.Options, F("#%d %v", i, value))
		}
		return listRefs		
	}

	image := Image("logo", Fname2url("images/unigui.png"), func(v Any)Any{return Info(F("%v logo selected!"))})

	replaceImage := func(val Any) Any{
		image.Image = Fname2url(F("%s/%v", UploadDir, val))
		return image
	}
	block := Block("X Block", Seq(cleanButton, selector), Seq(table, listRefs))
	block.Icon = "api"

	replaceButton := UploadButton("Replace the logo", replaceImage, "")	

	bottomBlock := Block("Bottom block", Seq(replaceButton, Button("Happy signal", 
		func(v Any)Any{return Signal{replaceButton, "make everyone happy"}}, "")), image)

	bottomBlock.Dispatch = blockDispatch

	scr :=  Screen(Seq(block, bottomBlock), user.SharedBlock("Audios"))	

	scr.Handle(selector, "Changed", func(v Any) Any{
		if v == "Based"{
			return UpdateError(selector, "Select can not be Based!")
		}
		return nil
	})	
	return scr
}

func main(){		
	//register shared blocks
	ShareBlock(sharedAudios, "Audios")
	//register screens
	Register(screenTest, "Main", 0, "insights")	
	Start()
}