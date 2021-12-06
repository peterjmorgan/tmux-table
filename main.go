package main

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func checkCommandExists(cmd string) (result bool) {
	checkCmd := exec.Command("command", "-v", cmd)
	output, _ := checkCmd.Output()
	if len(output) == 0 {
		return false // fail
	}
	return true // success
}

type tmuxWin struct {
	Name string
	Count int
	Datetime time.Time
}

func main() {
	cmd := "tmux"
	if !checkCommandExists(cmd) {
		fmt.Print("tmux command doesn't exist")
	}
	tmuxLsCmd := exec.Command("tmux", "ls")
	tmuxLsOut, err := tmuxLsCmd.Output()
	if err != nil {
		fmt.Printf("exec failed: %s\n", tmuxLsCmd.String())
	}
	linePat := regexp.MustCompile(`(.*?): (\d+) windows \(created (.*?)\)`)
	windows := []tmuxWin{}
	tz,err := time.LoadLocation("America/Phoenix")
	if err != nil {
		panic(err.Error())
	}
	time.Local = tz

	for _, line := range (strings.Split(string(tmuxLsOut), "\n")) {
		if linePat.MatchString(line) {
			lineMatch := linePat.FindAllStringSubmatch(line, -1)
			count, _ := strconv.Atoi(lineMatch[0][2])
			datetime,err := dateparse.ParseLocal(lineMatch[0][3])
			if err != nil {
				panic(err.Error())
			}
			windows = append(windows, tmuxWin{
				lineMatch[0][1],
				count,
				datetime,
			})
		}
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Session Name", "Num windows", "Started On"})
	for _, row := range windows {
		t.AppendRow(table.Row{
			row.Name,
			row.Count,
			row.Datetime,
		})
	}
	t.Render()
}