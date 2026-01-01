package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/lukpank/go-glpk/glpk"
)

const FILENAME = "../input.txt"

type State []int

func (state State) ToString() string {
	str := "{"
	for index, num := range state {
		str += strconv.Itoa(num)
		if index < len(state)-1 {
			str += ","
		}
	}
	str += "}"
	return str
}

type Action []int

type Machine struct {
	DesiredState State
	Actions      []Action
}

func parseInput() ([]Machine, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return nil, errors.New("there was an issue opening the file " + FILENAME)

	}
	actionsRegexp, err := regexp.Compile(`\((\d+,?)+\)`)
	if err != nil {
		panic("there is an issue with the regexp for light diagrams")
	}
	joltageRegexp, err := regexp.Compile(`\{(\d+,?)+\}`)
	if err != nil {
		panic("there is an issue with the regexp for light diagrams")
	}

	scanner := bufio.NewScanner(file)
	var machines []Machine
	for scanner.Scan() {
		line := scanner.Text()

		actionStrings := actionsRegexp.FindAllString(line, -1)
		if actionStrings == nil {
			return nil, errors.New("this line contained no actions! line: " + line)
		}

		joltageString := joltageRegexp.FindString(line)
		if len(joltageString) == 0 {
			return nil, errors.New("this line contained no joltages! line: " + line)
		}

		var actions []Action
		for _, actionString := range actionStrings {
			actionNumberStrings := strings.Split(actionString[1:len(actionString)-1], ",")
			var action Action
			for _, actionNumberString := range actionNumberStrings {
				actionNumber, err := strconv.Atoi(actionNumberString)
				if err != nil {
					return nil, errors.New("failed to convert numbers in actions to int. action: " + actionString)
				}
				action = append(action, actionNumber)
			}
			actions = append(actions, action)
		}

		var desiredState State
		joltageString = joltageString[1 : len(joltageString)-1]
		for _, joltageStr := range strings.Split(joltageString, ",") {
			joltage, err := strconv.Atoi(joltageStr)
			if err != nil {
				return nil, errors.New("failed to convert joltages to ints. joltages: " + joltageString)
			}
			desiredState = append(desiredState, joltage)
		}
		machines = append(machines, Machine{
			DesiredState: desiredState,
			Actions:      actions,
		})
	}
	return machines, nil
}

func main() {
	machines, err := parseInput()
	if err != nil {
		println("parsing failed! " + err.Error())
		return
	}
	total := 0
	for _, machine := range machines {
		lp := glpk.New()
		defer lp.Delete()
		lp.SetProbName("minimize sum")
		lp.SetObjName(machine.DesiredState.ToString())
		lp.SetObjDir(glpk.MIN)
		lp.AddRows(len(machine.DesiredState))
		for index, desiredJoltage := range machine.DesiredState {
			lp.SetRowName(index+1, "c"+strconv.Itoa(index+1))
			lp.SetRowBnds(index+1, glpk.FX, float64(desiredJoltage), float64(desiredJoltage))
		}
		lp.AddCols(len(machine.Actions))
		ind := make([]int32, len(machine.Actions)+1)
		for index := range ind {
			ind[index] = int32(index)
		}
		coefficients := make([][]float64, len(machine.DesiredState)+1)
		for index := range coefficients {
			coefficients[index] = make([]float64, len(machine.Actions)+1)
		}
		for index, action := range machine.Actions {
			lp.SetColName(index+1, "x"+strconv.Itoa(index))
			lp.SetColBnds(index+1, glpk.LO, 0, 0)
			lp.SetColKind(index+1, glpk.IV)
			lp.SetObjCoef(index+1, 1.0)
			for _, affectedIndex := range action {
				coefficients[affectedIndex+1][index+1] = 1
			}
		}
		for index, row := range coefficients {
			if index == 0 {
				continue
			}
			lp.SetMatRow(index, ind, row)
		}

		iocp := glpk.NewIocp()
		iocp.SetMsgLev(glpk.MSG_ERR) // only errors & warnings
		iocp.SetPresolve(true)       // Enable presolving for efficiency
		if err := lp.Intopt(iocp); err != nil {
			log.Fatalf("ILP error: %v", err)
		}

		// Output results
		fmt.Printf("Found the Answer for %s -> %g (", lp.ObjName(), lp.MipObjVal())
		for i := 1; i <= len(machine.Actions); i++ {
			fmt.Printf("%g", lp.MipColVal(i))
			if i < len(machine.Actions) {
				fmt.Print(",")
			}
		}
		print(")\n")
		total += int(lp.MipObjVal())
	}
	fmt.Println("Total: ", total)
}
