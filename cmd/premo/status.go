package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/meshplus/bitxhub-kit/fileutil"
	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/exec"
	gops "github.com/shirou/gopsutil/process"
	"github.com/urfave/cli/v2"
)

var processes = []string{
	"bitxhub",
	"pier-ether start",
	"pier-fabric start",
	"pier-flato start",
}

var statusCMD = &cli.Command{
	Name:   "status",
	Usage:  "List the status of instantiated components",
	Action: showStatus,
}

func showStatus(ctx *cli.Context) error {
	repoRoot, err := repo.PathRoot()
	if err != nil {
		return err
	}
	if !fileutil.Exist(repoRoot) {
		return fmt.Errorf("please run `premo init` first")
	}

	var table [][]string
	table = append(table, []string{"Name", "Component", "PID", "Status", "Created Time", "Args"})

	for _, process := range processes {
		table, err = existProcess(repoRoot, process, table)
		if err != nil {
			return err
		}
	}

	PrintTable(table, true)
	return nil
}

func existProcess(repoPath string, pro string, table [][]string) ([][]string, error) {
	pids, err := getProccessPid(pro)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(pids); i++ {
		status := "TERM"
		pid, err := strconv.Atoi(pids[i])
		if err != nil {
			return table, err
		}
		exist, err := gops.PidExists(int32(pid))
		if err == nil && exist {
			status = "RUNNING"
		}
		process, err := gops.NewProcess(int32(pid))
		if err != nil {
			continue
		}
		createTime, err := process.CreateTime()
		if err != nil {
			continue
		}
		tm := time.Unix(0, createTime*int64(time.Millisecond))
		timeFormat := tm.Format(time.RFC3339)

		component, _ := process.Name()
		slice, _ := process.CmdlineSlice()
		args := strings.Join(slice, " ")

		table = append(table, []string{
			fmt.Sprintf(strings.Split(pro, " ")[0]+"-%d", i),
			component,
			strconv.Itoa(pid),
			status,
			timeFormat,
			args,
		})
	}
	return table, nil
}

func getProccessPid(process string) ([]string, error) {
	arg := fmt.Sprintf("ps aux | grep '%v' | grep -v grep | awk '{print $2}'", process)
	bytes, err := exec.Execute("", arg)
	if err != nil {
		return nil, err
	}
	pids := strings.Split(string(bytes), "\n")
	return pids[:len(pids)-1], nil
}

// PrintTable accepts a matrix of strings and print them as ASCII table to terminal
func PrintTable(rows [][]string, header bool) {
	// Print the table
	t := tabby.New()
	if header {
		addRow(t, rows[0], header)
		rows = rows[1:]
	}
	for _, row := range rows {
		addRow(t, row, false)
	}
	t.Print()
}

func addRow(t *tabby.Tabby, rawLine []string, header bool) {
	// Convert []string to []interface{}
	row := make([]interface{}, len(rawLine))
	for i, v := range rawLine {
		row[i] = v
	}

	// Add line to the table
	if header {
		t.AddHeader(row...)
	} else {
		t.AddLine(row...)
	}
}
