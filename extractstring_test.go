package extractstring

import (
	"bytes"
	"testing"
)

func TestGetStrings(t *testing.T) {
	p, err := New("lua.json")
	if err != nil {
		t.Error(err)
		return
	}
	context := []byte(`t[1] = "this is\' \" s\"tring"
	--t = "this is\' \" s\"tring2"
	--[[
	v = "this is\' \" s\"tring3"
	]]
	s = [["this is\' \" s\"tring4"]]
	print('this is\' \" s\"tring5')`)

	entryBytes, _, _ := p.GetStrings(context, nil)
	if bytes.Equal(entryBytes[0], []byte(`this is\' \" s\"tring`)) == false {
		t.Error("entryBytes[0] is", string(entryBytes[0]))
	}
	if bytes.Equal(entryBytes[1], []byte(`"this is\' \" s\"tring4"`)) == false {
		t.Error("entryBytes[1] is", string(entryBytes[1]))
	}
	if bytes.Equal(entryBytes[2], []byte(`this is\' \" s\"tring5`)) == false {
		t.Error("entryBytes[2] is", string(entryBytes[2]))
	}

}
