package main

import . "github.com/Claus1/unigui-go"		

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
	Register(screenTest, "Main", 0, "insights")	
	Start()
}