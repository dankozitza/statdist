package statdist

import (
	"encoding/json"
	"fmt"
	"github.com/dankozitza/logdist"
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
var id_cnt int = 0
var access_log string

// Handle
//
// Sets Stat objects in stat_map.
//
func Handle(s Stat, quiet bool) {

	if !quiet {
		logdist.Message("", true, "["+s.ShortStack+"]["+
			s.Status+"]["+strconv.Itoa(s.Id)+"] "+s.Message+"\n")
	}
	stat_map[strconv.Itoa(s.Id)] = s
}

// RmHandle
//
// deletes a Stat object from stat_map.
//
func RmHandle(s Stat) {
	delete(stat_map, strconv.Itoa(s.Id))
}

// GetId
//
// used to give each Stat object a unique id number
//
func GetId() int {
	giveid := id_cnt
	id_cnt += 1
	return giveid
}

// HTTPHandler
//
// Handler used to write stat_map to http.ResponseWriter.
// Add to a http object with:
//
// 	var jsm statdist.JSONStatMap
// 	http.Handle("/stat", jsm)
//
type HTTPHandler string

func (j HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m_map, err := json.MarshalIndent(stat_map, "", "   ")
	if err != nil {
		panic(ErrStatDistGeneric(err.Error()))
	}

	fmt.Fprint(w, string(m_map))

	if access_log != "" {
		m_request, err := json.Marshal(r)
		if err != nil {
			panic(ErrStatDistGeneric(err.Error()))
		}
		logdist.Message(access_log, false, string(m_request)+"\n")
	}
}

// SetAccessLog
//
// Sets the file path for logging http.Request objects
//
func SetAccessLog(f string) {
	access_log = f
}
