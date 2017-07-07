package main

import (
	"fmt"
	"log"
	"net/http"
    "path"
	"strings"
	"encoding/json"
    "os"
    "os/exec"
	"bytes"
)


type CameraStatus struct {
  CameraConnected    bool
}

func main() {
	http.HandleFunc("/", handleMainPage)

	log.Println("Starting webserver on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Println("failed to start server")
		log.Fatal("http.ListendAndServe() failed with %s\n", err)
	}
}

func pathToSlice(path string) []string {
	a := strings.SplitAfter(path, "/")
	var newSlice []string

	for _, s := range a {
		s = strings.Trim(s, "/")
		s = strings.ToLower(s)
		if(s != "") {
			newSlice = append(newSlice, s)
		}
	}	

	return newSlice
}

func handleApiCall(pathSlice []string, w http.ResponseWriter, r *http.Request) {
	f := pathSlice[1]
	
	switch f {
        case "status":
            camConnected := false
            if _, err := os.Stat("/dev/video0"); err == nil {
              camConnected = true
            }

            status := CameraStatus{camConnected}
			js, _ := json.Marshal(status)
			w.Header().Set("Content-Type", "application/json")
  			w.Write(js)

		case "photo":
            appDir := os.Getenv("SNAP_APP_PATH")
            appPath := path.Join(appDir, "/bin/shoot")

            cmd := exec.Command(appPath)
            cmd.Stdin = strings.NewReader("some input")
	        var out bytes.Buffer
	        cmd.Stdout = &out
	        err := cmd.Run()
	
	        if err != nil {
	            fmt.Println(err)
	        }

            dir := os.Getenv("SNAP_APP_DATA_PATH")
		    fp := path.Join(dir, "shot.jpg")
            fmt.Println("retrieving image from: " + fp)
            http.ServeFile(w, r, fp)
	}
}

func sendOkMessage(w http.ResponseWriter, message string) {
			w.Header().Set("Server", "REST Button")
			w.WriteHeader(200)
			w.Write([]byte(message))
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if( r.URL.Path == "/" ) {
		fp := path.Join( os.Getenv("$SNAP_APP_PATH"), "templates", "index.html")
        http.ServeFile(w, r, fp)
	} else {
		d := pathToSlice(r.URL.Path)
		m := d[0]
		switch m {
			case "api":
				handleApiCall(d, w, r)
		}
	}	
}
