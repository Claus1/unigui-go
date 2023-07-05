package main

import (
	. "github.com/Claus1/unigui-go"
)

func screenTest(user *User) *Screen_ {
	table := Table("Videos", 0, nil, []string{"Video", "Duration", "Links", "Mine"},
		SeqSeq(Seq("opt_sync1_3_0.mp4", "30 seconds", "@Refer to signal1", true),
			Seq("opt_sync1_3_0.mp4", "37 seconds", "@Refer to signal8", false)))

	cleanTable := func(v Any) Any {
		table.Rows = SeqSeq()
		return table
	}
	cleanButton := Button("Clean table", cleanTable, "")

	selector := Select("Select", "All", nil, []string{"All", "Based", "Group"})

	listRefs := List("Detail ref list signals", "", nil, []string{"Select reference"})

	image := Image("logo", Fname2url("images/unigui.png"), func(v Any) Any { return Info(F("%v logo selected!")) })

	replaceImage := func(val Any) Any {
		image.Image = Fname2url(F("%s/%v", UploadDir, val))
		return image
	}

	replaceButton := UploadButton("Replace the logo", replaceImage, "")

	block := Block("X Block", Seq(cleanButton, selector, Button("Happy signal",
		func(v Any) Any { return Signal{replaceButton, "make everyone happy"} }, "")), Seq(table, listRefs))
	block.Icon = "api"

	bottomBlock := Block("Bottom block", Seq(replaceButton))	

	scr := Screen(Seq(block, bottomBlock), user.SharedBlock("Audios"))

	scr.Handle("Select", "X Block", "Changed", func(v Any) Any {
		if v == "Based" {
			return UpdateError(selector, "Select can not be Based!")
		}
		return false
	})
	return scr
}

func main() {
	//register shared blocks
	ShareBlock(sharedAudios, "Audios")
	//register screens
	Register(screenTest, "Main", 0, "insights")
	Start()
}
