package exec

import (
	"fmt"
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

	s := spin.New("\033[36mStart execute command: " + cmd.String() + "\033[m")
	s.Start()
	bytes, err := cmd.Output()
	if err != nil {
		s.Stop()
		if len(bytes) != 0 {
			PrintMessage(string(bytes), Red)
		}
		PrintMessage(err.Error(), Red)
		return nil, fmt.Errorf(err.Error())
	}
	s.Stop()
	PrintMessage(string(bytes), Orange)
	return bytes, nil
}

func Execute(repo string, args ...string) ([]byte, error) {
	var arg []string
	arg = append(arg, "-c")
	arg = append(arg, args...)
	cmd := exec.Command(DefaultShell, arg...)
	cmd.Dir = repo
	bytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return bytes, nil
}

func PrintMessage(str string, color uint64) {
	strs := strings.Split(str, "\n")
	for i := 0; i < len(strs); i++ {
		fmt.Printf("\033[%vm%v \033[m\n", color, strs[i])
	}
}
