package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"encoding/csv"
	"io"
	"sync"
	"runtime"
)

var wg sync.WaitGroup

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	wg.Add(3)
	var start, end string
	fmt.Println("Od datum: MM-DD\n=>")
	fmt.Scan(&start)
	fmt.Println("Do datum: MM-DD\n=>")
	fmt.Scan(&end)
	get_start, _ := time.Parse("01-02", start)
	get_end, _ := time.Parse("01-02", end)
	start_date := time.Date(2018, get_start.Month(), get_start.Day(), 0, 0, 0, 0, time.UTC)
	end_date := time.Date(2018, get_end.Month(), get_end.Day(), 0, 0, 0, 0, time.UTC)
	go data_log(start_date, end_date)
	go prv_envlog(start_date, end_date)
	go vtor_envlog(start_date, end_date)
	wg.Wait()
}

func data_log(start_date, end_date time.Time) {
	file, _ := os.Open("data_log.txt")
	defer file.Close()
	citac := bufio.NewScanner(file)
	ss1, _ := os.Create("ss1.csv")
	ss1.WriteString("Date,TempSS1\n")
	defer ss1.Close()
	ss24, _ := os.Create("ss24.csv")
	ss24.WriteString("Date,TempSS24\n")
	defer ss24.Close()
	hum, _ := os.Create("humidity.csv")
	hum.WriteString("Date,Humidity\n")
	defer hum.Close()
	pla, _ := os.Create("plafon.csv")
	pla.WriteString("Date,Plafon\n")
	defer pla.Close()
	for citac.Scan() {
		vreme, _ := time.Parse("02-01-2006\t15:04:05", citac.Text()[:19])
		if vreme.After(start_date) && vreme.Before(end_date) {
			prv_del := vreme.Format("2006-01-02 15:04:05")
			posleden := citac.Text()[len(citac.Text())-1:]
			if posleden == "1" {
				ss1.WriteString(prv_del + `,` + citac.Text()[32:36] + "\n")
				continue
			}
			if posleden == "7" {
				ss24.WriteString(prv_del + `,` + citac.Text()[32:36] + "\n")
				continue
			}
			if posleden == "y" {
				hum.WriteString(prv_del + `,` + citac.Text()[29:33] + "\n")
				continue
			}
			if posleden == "n" {
				pla.WriteString(prv_del + `,` + citac.Text()[29:33] + "\n")
				continue
			}
		}
	}
	wg.Done()
}

func prv_envlog(start_date, end_date time.Time) {
	Env, _ := os.Open("envlog.csv")
	ss2, _ := os.Create("novss2.csv")
	ss2.WriteString("Date,Temperature,Humidity\n")
	defer ss2.Close()
	Env.Seek(30, 0)
	csvreader := csv.NewReader(bufio.NewReader(Env))
	csvreader.Comma = ','
	for {
		linija, err := csvreader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		vreme, _ := time.Parse("01/02/2006 15:04:05", linija[0]+" "+linija[1])
		if vreme.After(start_date) && vreme.Before(end_date) {
			ss2.WriteString(vreme.Format("2006-01-02 15:04:05") +
				"," + linija[2] +
				"," + linija[3] + "\n")
		}
	}
	wg.Done()
}

func vtor_envlog(start_date, end_date time.Time) {
	Env, _ := os.Open("envlog(1).csv")
	oh, _ := os.Create("novoh.csv")
	oh.WriteString("Date,Temperature,Humidity\n")
	defer oh.Close()
	Env.Seek(30, 0)
	csvreader := csv.NewReader(bufio.NewReader(Env))
	csvreader.Comma = '\t'
	for {
		linija, err := csvreader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		vreme, _ := time.Parse("01/02/2006 15:04:05", linija[0]+" "+linija[1])
		if vreme.After(start_date) && vreme.Before(end_date) {
			oh.WriteString(vreme.Format("2006-01-02 15:04:05") +
				"," + linija[2] +
				"," + linija[3] + "\n")
		}
	}
	wg.Done()
}
