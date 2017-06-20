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
	_ "time"
)

var (
	musicList []string
	spamList  []string //Loaded from database
	upgrader  = websocket.Upgrader{}
	new_fi    string
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [path to mp3s]\n", os.Args[0])
		return
	}

	//load spam list
	file, err := os.OpenFile("spam.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		spamList = append(spamList, strings.TrimSpace(scanner.Text()))
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
	fmt.Println(spamList)

	// musicList = []string{}
	err = filepath.Walk(os.Args[1:][0], func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), ".mp3") {
				musicList = append(musicList, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error in walking: ", err)
	}

	http.HandleFunc("/", Hhandler)
	http.HandleFunc("/ws", compute)
	http.HandleFunc("/bs", bs)

	log.Println("Running on :7899")
	err = http.ListenAndServe(":7899", nil)
	if err != nil {
		log.Fatal("listenAndServe", err)
	}
}

func bs(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "bootstrap.min.css")
}
func Hhandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Client connected", req.RemoteAddr)
	var v = struct {
		Host  string
		Count int
	}{
		req.Host,
		len(musicList),
	}
	t := template.Must(template.ParseFiles("socketed.html"))
	t.Execute(rw, &v)
}

func compute(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	webRex := "(www.|)[a-zA-Z0-9_\\-]+\\.[a-zA-Z]{2,4}"
	rx, _ := regexp.Compile(webRex)
	var fi_mod string
	iter := 0
	for _, fi := range musicList {
		iter = iter + 1
		//Cut mp3 suffix & we need to check for mp3 file only not folder
		fi_mod = filepath.Base(strings.TrimSuffix(fi, ".mp3"))
		//Only Website named songs - maybe spam websites
		if rx.MatchString(fi_mod) {
			//Spam List Match
			var flag = false
			var spam = ""
			for _, sp := range spamList {
				if strings.Contains(fi_mod, sp) {
					spam = sp
					flag = true
				}
			}
			//User need to input spam
			if spam == "" {
				var u = struct {
					Count   int
					Context string
				}{
					iter,
					fi_mod,
				}
				c.WriteJSON(&u)
				var v = struct{ Spam string }{}
				c.ReadJSON(&v)
				spam = v.Spam
			}
			// fmt.Println("SPAM", spam)
			//Here we'll have spam variable
			if flag == false {
				spamList = append(spamList, spam)
				appendToSpamDB(spam)
			}
			os.Rename(fi, strings.Replace(fi, spam, "", 1))
		}
	}
	var u = struct {
		Count   int
		Context string
	}{
		iter,
		fi_mod,
	}
	c.WriteJSON(&u)
	c.Close()
}

func appendToSpamDB(sp string) {
	file, _ := os.OpenFile("spam.txt", os.O_RDWR|os.O_APPEND, 0666)
	defer file.Close()
	file.WriteString(sp + "\n")
}
