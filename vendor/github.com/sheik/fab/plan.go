package fab

import (
	"flag"
	"fmt"
	"github.com/sheik/fab/pkg/log"
	"os"
	"sort"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "more verbose output")
)

func (steps Plan) Execute(name string) (err error) {
	if step, ok := steps[name]; ok {

		if step.Gate != nil {
			err = fmt.Errorf("target \"%s\" did not pass gate: %s",
				name,
				step.Gate.Error(),
			)
			return
		}

		if steps[name].Depends != "" {
			dependencies := strings.Split(steps[name].Depends, " ")

			for _, stepName := range dependencies {
				if !steps[name].executed {
					steps.ProcessTarget(stepName)
				}
			}
		}

		if !step.executed {
			if step.Check.Func != nil {
				if step.Check.Func(step.Check.Args) {
					log.Warn("skipping %s", name)
					step.executed = true
					steps[name] = step
					return
				}
			}
			log.Info("executing %s", name)
			if *verbose {
				fmt.Println(step.Command)
			}
			if step.Function != nil {
				step.Function.(func(...interface{}) error)()
			} else {
				if step.Interactive {
					err = InteractiveCommand(step.Command)
				} else {
					err = Exec(step.Command)
				}
			}
			if err != nil {
				return
			}
			step.executed = true
			steps[name] = step
		}

		return
	}
	return fmt.Errorf("build target \"%s\" not found", name)
}

type Step struct {
	Command      string
	Function     interface{}
	Precondition string
	Check        Check
	Gate         error
	Fail         string
	Help         string
	Depends      string
	Default      bool
	Interactive  bool
	executed     bool
}

type Check struct {
	Func func(any) bool
	Args any
}

type Plan map[string]Step

func Complete(args ...string) []string {
	return args
}

var UpdateStep = Step{
	Command: `
		GOPRIVATE=github.com/sheik go install github.com/sheik/fab/cmd/fab@latest
		GOPRIVATE=github.com/sheik go get github.com/sheik/fab@latest 
		if [[ -d ./vendor ]]; then
			echo "fab: updating vendor directory"
			go mod vendor
		fi
		`,
	Help: "update fab",
}

var HelpStep = Step{
	Help: "print help message for createfile",
}

func (steps Plan) PrintHelp(args ...interface{}) error {
	var items []string
	for name := range steps {
		items = append(items, name)
	}
	sort.Strings(items)
	for _, item := range items {
		if steps[item].Help != "" {
			fmt.Printf("%30s : %s\n", log.Green(item), steps[item].Help)
		}
	}
	return nil
}

func Run(steps Plan) {
	var err error
	flag.Parse()

	// populate Plan map with auto targets
	steps["update"] = UpdateStep
	HelpStep.Function = steps.PrintHelp
	steps["help"] = HelpStep

	target := flag.Arg(0)
	if target == "" {
		if target, err = steps.DefaultTarget(); err != nil {
			log.Error(err.Error())
			os.Exit(3)
		}
	}
	steps.ProcessTarget(target)
	os.Exit(0)
}

func (steps Plan) DefaultTarget() (string, error) {
	for target, step := range steps {
		if step.Default {
			return target, nil
		}
	}
	return "", fmt.Errorf("no default target found in createfile")
}

func (steps Plan) ProcessTarget(name string) {
	var err error
	preconditionFailed := false
	step := steps[name]
	if strings.Contains(step.Command, ":INPUT:") {
		step.Command = strings.ReplaceAll(step.Command, ":INPUT:", strings.Join(os.Args[2:], " "))
		steps[name] = step
	}

	if step.Precondition != "" {
		err = Exec(step.Precondition)
		if err != nil {
			preconditionFailed = true
			log.Error("failed precondition for %s", name)
			os.Exit(1)
		}
	}
	if !preconditionFailed {
		err = steps.Execute(name)
	}
	if err != nil || preconditionFailed {
		if step.Fail != "" {
			log.Warn("error running target \"%s\": failing over to \"%s\"", name, step.Fail)
			err = steps.Execute(step.Fail)
		}
		if err != nil {
			log.Error("error running target \"%s\": %s", name, err)
			os.Exit(2)
		}
	}
}
