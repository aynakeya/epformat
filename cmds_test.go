package main

import "testing"

func TestGetAllFiles(t *testing.T) {
	for _, a := range getAllFiles(".") {
		t.Log(a)
	}
	for _, a := range getAllFiles("cmds.go") {
		t.Log(a.info.Name(), a.path)
	}
}
