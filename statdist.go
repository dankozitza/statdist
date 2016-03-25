package statdist

import (
	"encoding/json"
	"fmt"
   "github.com/dankozitza/dkutils"
	"github.com/dankozitza/logdist"
   "github.com/nelsam/requests"
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

var prog_stat_map map[string]map[string]Stat = make(map[string]map[string]Stat)
var pname string = "main";

//var stat_map map[string]Stat = make(map[string]Stat)
var id_cnt int = 0
var access_log string

// Handle
//
// Sets Stat objects in default stat_map.
//
func Handle(s Stat, quiet bool) {
	if !quiet {
		logdist.Message("", true, "["+s.ShortStack+"]["+
			s.Status+"]["+strconv.Itoa(s.Id)+"] "+s.Message+"\n")
	}
   if (prog_stat_map[pname] == nil) {
      prog_stat_map[pname] = map[string]Stat{}
   }
   prog_stat_map[pname][strconv.Itoa(s.Id)] = s
}

// RmHandle
//
// deletes a Stat object from default stat_map.
//
func RmHandle(s Stat) {
	delete(prog_stat_map[pname], strconv.Itoa(s.Id))
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
// Handler used to write prog_stat_map to http.ResponseWriter.
// Add to a http object with:
//
// 	var jsm statdist.JSONStatMap
// 	http.Handle("/stat", jsm)
//
type HTTPHandler string

func (j HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m_map, err := json.MarshalIndent(prog_stat_map, "", "   ")
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

// post to this handler by calling:
// curl curl -d "program=pname&id=0&message=&short_stack=&stack=&status=PASS" \
// dankozitza.com/post_stat
type HTTPPostHandler string

func (j HTTPPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   stat_template := map[string]interface{}{
      "id":          int(0),
      "status":      string(""),
      "short_stack": string(""),
      "message":     string(""),
      "stack":       string(""),
      "program":     string("")}

   params, err := requests.New(r).Params()

   if err != nil { // send the template
      r_map, err := json.MarshalIndent(stat_template, "", "   ")
      if err != nil {
         fmt.Fprint(w, "statdist: could not marshal stat_template!\n")
         return
      }
      fmt.Fprint(w, string(r_map))
      return;
   }

   result, err := dkutils.DeepTypePersuade(stat_template, params);
   if err != nil {
      fmt.Fprint(w, "statdist: could not persuade input parameters to " +
            "conform to template\n")
      return
   }

   s := Stat{
      result.(map[string]interface{})["id"].(int),
      result.(map[string]interface{})["status"].(string),
      result.(map[string]interface{})["short_stack"].(string),
      result.(map[string]interface{})["message"].(string),
      result.(map[string]interface{})["stack"].(string)}

   // send the result
   result.(map[string]interface{})["links"] = []interface{}{
      &map[string]interface{}{
         "href": "/statdist",
         "rel":  "index"},
      &map[string]interface{}{
         "href": "/post_stat",
         "rel":  "self"}}
   r_map, err := json.MarshalIndent(result.(map[string]interface{}), "", "   ")
   if err != nil {
      fmt.Fprint(w, "statdist: could not marshal result!\n")
      return
   }

   fmt.Fprint(w, string(r_map))

   program := result.(map[string]interface{})["program"].(string)
   if (prog_stat_map[program] == nil) {
      prog_stat_map[program] = map[string]Stat{}
   }
   // set the result
   prog_stat_map[program][strconv.Itoa(s.Id)] = s
}

// SetAccessLog
//
// Sets the file path for logging http.Request objects
//
func SetAccessLog(f string) {
	access_log = f
}
