package statdist

import (
	"fmt"
	"github.com/dankozitza/logdist"
	"net/http"
	"syscall"
	"testing"
)

func TestHandler(t *testing.T) {
	var jsm JSONStatMap
	http.Handle("/statdist", jsm)
}

func TestHandle(t *testing.T) {
	s := Stat{GetId(), "status", "shortstack", "message", "stack"}
	Handle(s)

	if _, ok := stat_map["0"]; !ok {
		fmt.Println("TestHandle: stat_map does not have correct contents!")
		t.Fail()
	}
}

func TestRmHandle(t *testing.T) {
	s := Stat{0, "", "", "", ""}
	RmHandle(s)

	if _, ok := stat_map["0"]; ok {
		fmt.Println("TestRmHandle: stat_map does not have correct contents!")
		t.Fail()
	}
}

func TestSetAccessLog(t *testing.T) {
	var ldh logdist.LogDistHandler = "statdist_access.log"
	http.Handle("/statdist_access.log", ldh)

	SetAccessLog("statdist_access.log")
	//http.ListenAndServe("localhost:9000", nil)
}

func TestClean(t *testing.T) {
	fmt.Println("TestClean: removing statdist_access.log")
	syscall.Exec("/usr/bin/rm",
		[]string{"rm", "-f", "statdist_access.log"}, nil)
}
