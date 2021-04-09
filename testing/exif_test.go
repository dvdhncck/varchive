package main

import ( 
	"fmt"
	"testing"
	"davidhancock.com/varchive"
)

func Test_parseDimensions(t *testing.T) {
	
	w,h,e := varchive.GetVideoInfo(`test-data/one/sample 002.MTS`) // path relative to this .go file

	fmt.Printf("%d, %d, %s", w, h, e)
	
}
