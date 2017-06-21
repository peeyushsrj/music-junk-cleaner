/*Program to remove junk data in mp3 files*/
package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ResponseMsg struct {
	Count   int
	Context string
}

var (
	musicList []string
	junkList  []string
	upgrader  = websocket.Upgrader{}
)

func LinesFromFile(path string) ([]string, error) {
	var arr []string

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return arr, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr = append(arr, strings.TrimSpace(scanner.Text()))
	}
	if scanner.Err() != nil {
		return arr, scanner.Err()
	}
	return arr, nil
}

func BrowseXFiles(x string, root string) ([]string, error) {
	var arr []string
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), x) { //.mp3
				arr = append(arr, path)
			}
		}
		return nil
	})
	if err != nil {
		return arr, err
	}
	return arr, nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [path to mp3s]\n", os.Args[0])
		return
	}

	//load local junk list
	l, err := LinesFromFile("./junk.txt")
	if err != nil {
		log.Fatal("Error loading in junklist", err)
	}
	junkList = l

	// musicList = []string{}
	musicList, err = BrowseXFiles(".mp3", os.Args[1:][0])
	if err != nil {
		log.Fatal("Error in walking over files", err)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", compute)
	log.Println("Running on :7899")
	err = http.ListenAndServe(":7899", nil)
	if err != nil {
		log.Fatal("listenAndServe", err)
	}
}

func home(rw http.ResponseWriter, req *http.Request) {
	// fmt.Println("Client connected", req.RemoteAddr)
	var v = struct {
		Host  string
		Count int
	}{
		req.Host,
		len(musicList),
	}
	t := template.Must(template.ParseFiles("main.html"))
	t.Execute(rw, &v)
}

func compute(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade websocket error:", err)
		return
	}
	// defer c.Close()

	//Regex for websites junk (possible junk)
	webRex := "(www.|)[a-zA-Z0-9_\\-]+\\.[a-zA-Z]{2,4}"
	rx, _ := regexp.Compile(webRex)
	var iter = 0
	var cleaned_fi string

	//Scan music list
	for _, fi := range musicList {
		iter = iter + 1

		//Only base names, mp3 extension exclude
		cleaned_fi = filepath.Base(strings.TrimSuffix(fi, ".mp3"))

		//Possible junk with websites name, other exclude
		if rx.MatchString(cleaned_fi) {

			//junk List Match
			junk := stringInSlice(cleaned_fi, junkList)
			if junk == "" { //junk not found
				c.WriteJSON(&ResponseMsg{iter, cleaned_fi})

				var v = struct{ junk string }{}
				c.ReadJSON(&v)
				junk = v.junk

				//New junk from user added to local junk list
				junkList = append(junkList, junk)
				appendTojunkDB(junk)
			}
			// os.Rename(fi, strings.Replace(fi, junk, "", 1))
		}
	}
	c.WriteJSON(&ResponseMsg{iter, cleaned_fi})
	c.Close()
}

func appendTojunkDB(sp string) {
	if sp != "" {
		file, _ := os.OpenFile("junk.txt", os.O_RDWR|os.O_APPEND, 0666)
		defer file.Close()
		b := make([]byte, 1000) //this can be efficient
		file.Read(b)
		if !strings.Contains(string(b), sp) {
			file.WriteString(sp + "\n")
		}
	}
}

func stringInSlice(a string, b []string) string {
	for _, el := range b {
		if strings.Contains(a, el) {
			return el
		}
	}
	return ""
}
