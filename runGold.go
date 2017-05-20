package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var closeChan chan byte

//var runChan chan byte

func run() {
	go startHttp() //run http server
	path := ""
	t := NewTurtle()
	g, err := crawlGoldNow()
	if err != nil {
		log.Println(err)
		return
	}
	timeBefore := time.Now()
	h, m, d := timeBefore.Date()
	path = fmt.Sprintf("log/%d_%d_%d.log", h, m, d)
	file, err := os.Create(path)
	if err != nil {
		log.Println("Log err.")
	} else {
		log.SetOutput(file)
	}
	t.readTurtle()
	t.updateBase(g.price)
	for {
		select {
		case <-closeChan:
			return
		default:
		}
		//runChan <-'1'//send a message to runChan to work
		if time.Now().Weekday() != timeBefore.Weekday() {
			file.Close()
			timeBefore = time.Now()
			h, m, d := timeBefore.Date()
			path = fmt.Sprintf("log/%d_%d_%d.log", h, m, d)
			file, err = os.Create(path)
			if err != nil {
				log.Println("Log err.")
			} else {
				log.SetOutput(file)
			}
		} else {
			file.Close()
			file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				log.Println("Log err.")
			} else {
				log.SetOutput(file)
			}
		}
		go work(t)
		time.Sleep(time.Minute) //sleep one min
		t.saveTurtle()
	}
}

func work(t *Turtle) {
	times := time.Now()
	h, m, s := times.Clock()
	if times.Weekday() == 6 && h >= 4 || times.Weekday() == 0 || times.Weekday() == 1 && h < 8 {
		log.Println(times.Weekday(), ":", h)
		return
	}
	g, err := crawlGoldNow()
	if err != nil {
		log.Println(err)
		return
	}

	err = writeGoldMin(g) //write into db
	if err != nil {
		log.Println(err)
	}
	if h == 23 && m == 59 { //the last minute in one day
		year, mon, day := times.Date()
		g.date = fmt.Sprintf("%d-%d-%d %d:%d:%d", year, mon, day, h, m, s)
		tmpg, _ := readGoldDay(1)
		g.open = tmpg[0].price
		importDailyGold([]Gold{g}) //write into db
		t.updateBase(g.price)
	}

	t.run(g)

}

func test() {
	t := NewTurtle()
	t.readTurtle()
	time.Sleep(5 * time.Second)
}

func main() {
	initDB()
	initMail("./mail.txt")
	run()
}
