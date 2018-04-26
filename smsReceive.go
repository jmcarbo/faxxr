package main

import (
	"encoding/xml"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type smsMsg struct {
	XMLName xml.Name `xml:"Message"`
	Body    []string `xml:"Body,omitempty"`
	Media   []string `xml:"Media,omitempty"`
}

type smsML struct {
	XMLName  xml.Name `xml:"Response"`
	Message  *smsMsg  `xml:"Message,omitempty"`
	Redirect string   `xml:"Redirect,omitempty"`
}

func smsReceive(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Print("Unable to parse form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logSmsStatus(r.PostForm)

	msg := ""
	switch strings.Replace(strings.ToLower(r.PostForm.Get("Body")), " ", "", -1) {
	case "help", "options":
		msg = `faxxr options are:
help
options
settings
fax enable
fax disable`
	case "settings":
		msg += "faxxr settings:"
		config.Range(func(k, v interface{}) bool {
			msg += "\n" + k.(string) + " = " + v.(string)
			return true
		})
	case "faxenable":
		config.Store("fax", "enable")
		msg = "Fax enabled."
	case "faxdisable":
		config.Store("fax", "disable")
		msg = "Fax disabled."
	default:
		msgs := []string{
			"Say what?",
			"I don't understand.",
			"That's not something I can do.",
			"Maybe you should try Google.",
			"My vocabulary is limited.",
			"It's all Greek to me.",
		}
		msg = msgs[rand.Intn(len(msgs))]
		msg += " Try \"help\" or \"options\" to see what I can do."
	}

	w.Header().Set("Content-Type", "application/xml")
	data := &smsML{
		Message: &smsMsg{
			Body: []string{msg},
		},
	}

	b, err := xml.Marshal(data)
	w.Write([]byte(xml.Header))
	w.Write(b)
}
