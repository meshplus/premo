package exec

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/jiuhuche120/spin"
)

const (
	Red    = 31
	Orange = 33

	DefaultShell = "/bin/bash"
)

func ExecuteShell(repo string, args ...string) ([]byte, error) {
	var arg []string
	arg = append(arg, "-c")
	arg = append(arg, args...)
	cmd := exec.Command(DefaultShell, arg...)
	cmd.Dir = repo
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	s := spin.New("\033[36mStart execute command: " + cmd.String() + "\033[m")
	s.Start()
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	data1, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	data2, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	if len(data2) != 0 {
		s.Stop()
		if len(data1) != 0 {
			PrintMessage(string(data1), Red)
		}
		PrintMessage(string(data2), Red)
		return nil, fmt.Errorf(string(data2))
	}
	s.Stop()
	if len(data1) != 0 {
		PrintMessage(string(data1), Orange)
	}
	return data1, nil
}

func Execute(repo string, args ...string) ([]byte, error) {
	var arg []string
	arg = append(arg, "-c")
	arg = append(arg, args...)
	cmd := exec.Command(DefaultShell, arg...)
	cmd.Dir = repo
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	data1, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	data2, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	if len(data2) != 0 {
		return nil, fmt.Errorf(string(data2))
	}
	return data1, nil
}

func PrintMessage(str string, color uint64) {
	strs := strings.Split(str, "\n")
	for i := 0; i < len(strs); i++ {
		fmt.Printf("\033[%vm%v \033[m\n", color, strs[i])
	}
}
