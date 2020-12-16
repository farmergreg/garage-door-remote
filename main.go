package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var mutex = &sync.Mutex{}
var chip *gpiod.Chip
var line *gpiod.Line

func operateGarageDoor(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	t := time.Now()
	fmt.Fprintf(w, t.Format("Jan 2 15:04:05 2006 MST"))
	fmt.Fprintf(w, " Garage Door Opener Button Pressed")
	log.Print("Garage Door Button Pressed")

	line.SetValue(1)
	time.Sleep(50 * time.Millisecond)
	line.SetValue(0)

	mutex.Unlock()
}

func main() {
	var err error
	log.Print("Garage Door Opener")

	chip, err = gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	defer chip.Close()

	v := 0
	line, err = chip.RequestLine(rpi.GPIO21, gpiod.AsOutput(v))
	if err != nil {
		panic(err)
	}
	defer line.Close()

	log.Print("Waiting for commands...")
	http.HandleFunc("/", operateGarageDoor)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
