package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var conf *config

// Wheel struct
type Wheel struct {
	alphabet string
	offset   int
}

// Machine struct
type Machine struct {
	wheelOrder []string
	keyphrase  string
	wheels     [3]Wheel
}

// NewMachine generator
func NewMachine(wheelOrder []string, keyphrase string) *Machine {
	keyphrase = strings.ToUpper(keyphrase)
	var wheels = [3]Wheel{}
	if len(wheelOrder) == 3 {
		i := 0
		for index := range wheelOrder {
			wheel := new(Wheel)
			wheel.alphabet, _ = getAlphabet(index, wheelOrder)
			wheels[i] = *wheel
			wheels[i].offset = strings.Index(wheels[i].alphabet, string(keyphrase[i]))
			i++
		}
	}
	return &Machine{wheelOrder, keyphrase, wheels}
}

// Return the a wheel alphabet,true or false.
func getAlphabet(index int, wheelOrder []string) (string, bool) {
	alphabet, ok := conf.Wheels[string(wheelOrder[index][0])][string(wheelOrder[index][1])]
	return alphabet, ok
}

// Decode handler for web requests
func decodeHandler(w http.ResponseWriter, r *http.Request) {

}

// Encode handler for web requests
func encodeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, `<title>Enigma IV Encoder</title>`)
	if r.FormValue("text") == "" || r.FormValue("keyphrase") == "" {
		if r.FormValue("keyphrase") == "" {
			fmt.Fprintf(w, "<p>Enter a three letter keyphrase.</p>")
		}
		if r.FormValue("text") == "" {
			fmt.Fprintf(w, "<p>Enter text to encode.</p>")
		}
	} else {
		wheelOrder := []string{"1A", "3A", "5A"}
		keyphrase := r.FormValue("keyphrase")
		// Validate wheel selections.
		for index, wheelName := range wheelOrder {
			if _, ok := getAlphabet(index, wheelOrder); !ok {
				log.Fatal("Error in wheel selection: ", wheelName)
			}
		}

		message := r.FormValue("text")

		// Create a new machine.
		machine := *NewMachine(wheelOrder, keyphrase)

		encodedText := machine.encodeMessage(message)
		fmt.Fprintf(w, "Encoded text:</br>%s", encodedText)
		fmt.Println("Encoded")
	}
	fmt.Fprintf(w, `<form action="/encode" method="POST">
	Text to encode:</br>
	<textarea rows="5" cols="40" name="text">%s</textarea></br>
	Keyphrase: <input name="keyphrase" size="3" maxlength="3" value="%s" type="text"></input></br>
	<input type="submit" value="Encode">
	</form>`, r.FormValue("text"), r.FormValue("keyphrase"))
}

func main() {

	conf = getConfig()

	http.HandleFunc("/encode", encodeHandler)
	http.HandleFunc("/decode", encodeHandler)
	http.ListenAndServe(":9001", nil)
}
