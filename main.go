package main

import (
	"flag"
	"os"
	"math"
	"os/exec"
	"strings"
	"time"
	"fmt"
)

const INFO = 
`  ____ __  __ ____    ____   ____ _   _ _____ ____  _   _ _     _____ ____  
 / ___|  \/  |  _ \  / ___| / ___| | | | ____|  _ \| | | | |   | ____|  _ \ 
| |   | |\/| | | | | \___ \| |   | |_| |  _| | | | | | | | |   |  _| | |_) |
| |___| |  | | |_| |  ___) | |___|  _  | |___| |_| | |_| | |___| |___|  _ < 
 \____|_|  |_|____/  |____/ \____|_| |_|_____|____/ \___/|_____|_____|_| \_\
 ===========================================================================
 
Welcome to cmd scheduler!`

const DIVIDER = 
`===========================================================================`

type UserCommands []string

func (arr *UserCommands) String() string {

	str := ""
	for i := 0; i < len(*arr); i++ {
		val := fmt.Sprintf("%d", i  + 1)
		str += val + ". " + (*arr)[i] + " " + " \n"
	}

	return str
}

func (arr *UserCommands) Set(value string) error {
	*arr = append(*arr, strings.TrimSpace(value))
	return nil
}

func main() {
	zeroDur, erz := time.ParseDuration("0s")
	if erz != nil {
		println("ERROR: Can not initialize application. Exiting.")
		return
	}

	duration := flag.Duration("d", time.Second, "duration of sleep")
	commands := UserCommands{}
	flag.Var(&commands, "c", "The list of commands to run in order")
	continueOnError := flag.Bool("e", false, "continue on error of any one command")
	intraExecutionDuration := flag.Duration("r",zeroDur , "duration of sleep between commands")
	maxIter := flag.Int("m", math.MaxInt, "maximum number of iterations to execute")

	if (*duration) < zeroDur {
		println("Duration flag -d cannot be negative. Exiting application.")
		os.Exit(1)
	} else if (*intraExecutionDuration) < zeroDur {
		println("Intra-execution sleep flag -r cannot be negative. Exiting application.")
		os.Exit(1)
	}

	flag.Parse()

	println(INFO)
	fmt.Printf("Cmd Scheduler will run all commands in the current directory\n")
	fmt.Printf("The following command(s) will be run every %v second(s) in order:\n", duration)
	println(commands.String())
	println(DIVIDER)

	for q=0; q < maxIter; q++ {
		println("Sleeping...")
		time.Sleep(*duration)
		println("Start iteration " + fmt.Sprintf("%d", q  + 1))
		for i := 0; i < len(commands); i++ {
			
			println("Executing command " + commands[i])
			var err = runCmd(commands[i])

			if err != nil {
				print("ERR: Cannot run command. " + commands[i])
				if !(*continueOnError) {
					os.Exit(1)
				}
			}
			println("Successfully executed command " + commands[i])
			fmt.Printf("Sleeping %v second(s) in between commands:\n", intraExecutionDuration)
			time.Sleep(*intraExecutionDuration)
		}
		println("End iteration " + fmt.Sprintf("%d", q + 1))
		println(DIVIDER)
	}
}

func runCmd(cmd string) (err error) {
	tCmd := exec.Command("cmd.exe", "/c", cmd)
	//TODO add a commoand for linux/macOs
	out , err := tCmd.Output()
	strArr := strings.Split(string(out), ("\n"))
	for  i := 0; i < len(strArr); i++  {
		if i == len(strArr) - 1 {continue}
		println("RESULT:        " + strArr[i])
	}
	
	return err
}
