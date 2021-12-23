package main

// TODO clipboard stuff
/*
Could use clip.h like i did with C++ but check out this go package:
https://github.com/atotto/clipboard
*/
import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/phillvancejr/webview"
)

const (
	sidewallFile = "sidewall.json"
	// NOTE, port is also used in web_gen, so you need to update it if it is changed here
	port   = "8080"
	width  = 502
	height = 380
)

// TODO have default json to serve if for some reason you can't get the local cached file
var (
	//defaultJSon string
	//	defaultMatches =
	//go:embed web
	webFS embed.FS
)

func main() {
	// check for updated matches
	//go checkUpdate()
	// launch server
	go serve()
	// webview stuff on main thread
	runtime.LockOSThread()
	w := webview.New(true)
	defer w.Destroy()
	w.SetSize(width, height, webview.HintFixed)
	w.Center()
	w.Topmost(true)
	w.SetTitle("Vomela Sidewall Searcher")
	w.Navigate("http://localhost:" + port + "/app")
	w.Run()
}

type sidewallMap map[string]map[string]string

// creats a sidewallMap from map[string]interface{} s
func sidewallMapFromUnmarshal(in map[string]interface{}) (m sidewallMap, count int) {

	result := make(sidewallMap)
	// I have to convert the map[string]interface{} into a map[string]string
	for vomela, sidewalls := range in {
		if _, exists := result[vomela]; !exists {
			result[vomela] = make(map[string]string)
		}
		temp := make(map[string]string)
		// first I have to cast sidewalls from interface{} (due to json.UnMarshal) to map[string]interface{}
		for brand, match := range sidewalls.(map[string]interface{}) {
			// then I can cast the inner interface{} to string
			temp[brand] = match.(string)
		}
		result[vomela] = temp
	}
	count = countEntries(&result)

	return result, count
}

/*
Checks if the result sidewall.json file is out of sync with the version on github, if so then use the version from github
*/
func checkUpdate() {
	if matches, fetchedCount, e := fetchMatches(); e != nil {
		log.Fatal(e)
	} else {
		data, e := ioutil.ReadFile(sidewallFile)
		if e != nil {
			log.Fatal(e)
		}

		var temp map[string]interface{}

		json.Unmarshal(data, &temp)

		_, count := sidewallMapFromUnmarshal(temp)

		// if they are out of sync, used the fetched version
		if fetchedCount != count {
			jsonFile, _ := os.Create(sidewallFile)

			encoder := json.NewEncoder(jsonFile)
			encoder.Encode(matches)
			// TODO  - update result json string?
			// TODO - notify user there was an update and prompt restart/refresh
			// or just do it, webview should allow me to execute JS from the Go side
		}
	}

}

// counts entries in map to get a value that can be used to determine if two sidewallMaps are the same
func countEntries(m *sidewallMap) int {
	count := 0
	for _, match := range *m {
		count++

		for range match {
			count++
		}
	}
	return count
}

// fetches the sidewall.json from github and converst to a sidewallMap
func fetchMatches() (sidewallMap, int, error) {
	resp, e := http.Get("https://raw.githubusercontent.com/phillvancejr/Vomela-Sidewall-Searcher/master/sidewall_data/sidewall.json")
	if e != nil {
		return nil, 0, fmt.Errorf("sidewall.json get request failed: %v", e)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var rawMatches map[string]interface{}

	if e := json.Unmarshal(body, &rawMatches); e != nil {
		return nil, 0, fmt.Errorf("error unmarshalling json!: %v", e)
	}

	result, count := sidewallMapFromUnmarshal(rawMatches)

	return result, count, nil
}

func serve() {
	fmt.Println("Listening on port: ", port)
	// send back sidewall.json
	http.HandleFunc("/matches", func(w http.ResponseWriter, req *http.Request) {
		if sidewallJSON, e := ioutil.ReadFile(sidewallFile); e != nil {
			// TODO navigate webview to error page or show error window by calling JS from webview dispatch to show error
			log.Fatal(e)
		} else {
			fmt.Fprint(w, string(sidewallJSON))
		}
	})
	// serve web folder
	fsys, _ := fs.Sub(webFS, "web")
	http.Handle("/", http.FileServer(http.FS(fsys)))

	http.HandleFunc("/app", func(w http.ResponseWriter, req *http.Request) {
		t, e := template.ParseFS(webFS, "web/template.html")

		if e != nil {
			log.Fatal("template parsing error!: ", e)
		}

		e = t.Execute(w, port)

		if e != nil {
			log.Fatal("template execution error!: ", e)
		}

	})

	http.ListenAndServe(":"+port, nil)
}

//generate sidewall json file
//go:generate go run sidewall_data/sidewall_gen.go
