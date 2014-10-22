package statdist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ErrStatDistGeneric string

func (e ErrStatDistGeneric) Error() string {
	return "statdist error: " + string(e)
}

type Stat struct {
	Id         int
	Status     string
	ShortStack string
	Message    string
	Stack      string
}

var stat_map map[string]Stat = make(map[string]Stat)

// Handle
//
// Sets Stat objects in stat_map.
//
func Handle(s Stat) {
	stat_map[strconv.Itoa(s.Id)] = s
}

// RmHandle
//
// deletes a Stat object from stat_map.
//
func RmHandle(s Stat) {
	delete(stat_map, strconv.Itoa(s.Id))
}

// JSONStatMap
//
// Handler used to write stat_map to http.ResponseWriter.
// Add to a http object with:
//
// 	var jsm statdist.JSONStatMap
//		http.Handle("/stat", jsm)
//
type JSONStatMap string

func (j JSONStatMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m_map, err := json.MarshalIndent(stat_map, "", "   ")
	if err != nil {
		panic(ErrStatDistGeneric(err.Error()))
	}

	fmt.Fprint(w, string(m_map))

	// may want to keep this somewhere along with logs
	// use syslog.New(priority Priority, tag string) (w *Writer, err error) maybe
	// in a different package.
	//
	// will have to call logdist manually
	//
	//fmt.Println(seestack.Short(), "r:")
	//fmt.Println(r)

}