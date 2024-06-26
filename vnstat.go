package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	args = flag.String("p", "-m", "vnstat args")
	max  = flag.Float64("max", 999.00, "max gb traffic todo")
	loop = flag.Int64("loop", 5, "how many minutes to check loop")
	ver  = flag.Int64("ver", 2, "the version of vnstat")
)

const (
	GiB = "GiB"
	MiB = "MiB"
	TiB = "TiB"
)

func main() {
	flag.Parse()
	// check on startup
	check()
	// check loop
	tk := time.NewTicker(time.Minute * time.Duration(*loop))
	for {
		<-tk.C
		check()
	}
}

func check() {
	if *ver == 1 {
		RunCommand("vnstat", "-u")
	}
	result := RunCommand("vnstat", *args)
	lines := strings.Split(result, "\n")
	switch *args {
	case "-m":
		dataLines := make([]string, 0)
		for _, line := range lines {
			if strings.Contains(line, "bit/s") {
				dataLines = append(dataLines, line)
			}
		}
		// recent month
		recent := dataLines[len(dataLines)-1]
		vvs := strings.Split(recent, "|")
		for i, v := range vvs {
			if i == 1 {
				v = strings.Trim(v, " ")
				total := strings.Split(v, " ")
				if len(total) == 2 {
					amount, _ := strconv.ParseFloat(total[0], 32)
					dw := total[1]
					switch dw {
					case GiB:
						checkTr(amount, recent)
					case TiB:
						amount = amount * 1024
						checkTr(amount, recent)
					default:
						log.Printf("ignore the wd -> %s", dw)
						return
					}
				}
			}
		}
	case "--json":
	default:
		log.Printf("unhandle result:: %s", result)
	}
}

func checkTr(amount float64, recent string) {
	if amount < *max {
		log.Printf("traffic ok, used:: %.2f, left:: %.2f \norigin data::%s", amount, *max-amount, recent)
	} else {
		log.Printf("oops traffic used up, god!!! exec poweroff cmd now.")
		RunCommand("poweroff")
	}
}

func RunCommand(cmdName string, arg ...string) string {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, cmdName, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	// run
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// result
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	cancel()
	cmd.Wait()
	log.Println(string(opBytes))
	return string(opBytes)
}
