package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "io/ioutil"
    "os"
    "os/exec"
	"encoding/json"
	"strconv"	
    "time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func check(e error) {
    if e != nil {
        log.Print(e)
    }
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func sendphoto(photo string) {
		bot, err := tgbotapi.NewBotAPI(BOTID+":"+TOKEN)
			if err != nil {
				log.Panic(err)
			}
		
		log.Printf("Authorized on account %s", bot.Self.UserName)
		
		for i, chatid := range CHATIDS {
			fmt.Println(i, "-->", chatid)
			effe, err := strconv.ParseInt(chatid, 10, 64)
			check(err)
			msg := tgbotapi.NewPhotoUpload(effe, photo)
			//msg.Caption = "Van nieuwe server"
			_, err2 := bot.Send(msg)

			if err2 != nil {
				log.Panic(err2)
			}
		}

}

func sendvideo(video string) {
    // convert video to telegram compatible format
    err := os.Remove("/tmp/voordeur.mp4") 
    check(err)
    lsCmd := exec.Command("/usr/bin/ffmpeg", "-i", video, "-c:a", "copy", "-s", "640x360", "/tmp/voordeur.mp4")
    lsOut, err := lsCmd.Output()
    check(err)
    fmt.Println(string(lsOut))


		fmt.Println(CHATIDS)
		bot, err := tgbotapi.NewBotAPI(BOTID+":"+TOKEN)
	    check(err)
		
		log.Printf("Authorized on account %s", bot.Self.UserName)
		
		for i, chatid := range CHATIDS {
			fmt.Println(i, "sending video to -->", chatid)
			effe, err := strconv.ParseInt(chatid, 10, 64)
			check(err)
			msg := tgbotapi.NewVideoUpload(effe, "/tmp/voordeur.mp4")
			// msg.Caption = "Van nieuwe server"
			_, err2 := bot.Send(msg)

			if err2 != nil {
				log.Panic(err2)
			}
		}

}


// Get the pathToVideo from the kerberos.io webhook
func kerberosio(rw http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
    //d.DisallowUnknownFields() // error if user sends extra data

	// anonymous struct type: handy for one-time use
	t := struct {
		VideoPath *string `json:"pathToVideo"` // pointer to string, so we can test for field absence
	}{}

	err := d.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if t.VideoPath == nil {
		http.Error(rw, "missing field 'pathToVideo' from JSON object", http.StatusBadRequest)
		return
	}

	// optional check
	if d.More() {
		http.Error(rw, "extraneous data after JSON object", http.StatusBadRequest)
		return
	}

	// got all fields we expected: no more, no less

	log.Println("/etc/opt/kerberosio/capture/" + *t.VideoPath)

	// if less than 30 seconds passed since doorbell ring send video by telegram
	Duration := time.Since(Ringtime)
	if Duration.Seconds() < 30 {
		fmt.Println(Ringtime)
    	sendvideo("/etc/opt/kerberosio/capture/" + *t.VideoPath)
		Ringtime = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
}

// Print complete body
func debug(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadAll(r.Body)
    log.Printf("incoming message - %s", body)
	// now := time.Now()
	Duration := time.Since(Ringtime)
	if Duration.Seconds() < 30 {
		fmt.Println(Ringtime)
		Ringtime = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}


    defer r.Body.Close()
}

func setringtime(w http.ResponseWriter, r *http.Request) {
	// Set ringtime when doorbell is pressed
	Ringtime = time.Now()
	fmt.Println(Ringtime)

    // get jpg from mjpeg stream
    err := os.Remove("/tmp/voordeur.jpg") 
    check(err)
    lsCmd := exec.Command("/usr/bin/ffmpeg", "-i", "http://localhost:8889", "-frames", "1", "/tmp/voordeur.jpg")
    lsOut, err := lsCmd.Output()
    check(err)
    fmt.Println(string(lsOut))
	sendphoto("/tmp/voordeur.jpg")
}

//Global vars
var Ringtime time.Time
var BOTID string
var TOKEN string
var CHATIDS []string

func main() {
	CHATIDS = strings.Fields(os.Getenv("CHAT_IDS"))
	BOTID = os.Getenv("BOT_ID")
	TOKEN = os.Getenv("TOKEN")

    http.HandleFunc("/ring", setringtime)
    http.HandleFunc("/debug/", debug)
    http.HandleFunc("/kerberosio/", kerberosio)
    http.HandleFunc("/headers", headers)

    http.ListenAndServe(":8090", nil)

}