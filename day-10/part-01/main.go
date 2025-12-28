package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const FILENAME = "../input.txt"

type State string

type Action []int

type Record struct {
	Steps    int
	PrevStep *Record
	State    State
}

type Machine struct {
	DesiredState    State
	Actions         []Action
	Queue           *Queue[Record]
	ProcessedStates map[State]struct{}
}

func parseInput() ([]Machine, error) {
	file, err := os.Open(FILENAME)
	if err != nil {
		return nil, errors.New("there was an issue opening the file " + FILENAME)

	}
	lightStateRegexp, err := regexp.Compile(`^\[[.#]+\] `)
	if err != nil {
		panic("there is an issue with the regexp for light diagrams")
	}
	actionsRegexp, err := regexp.Compile(`\((\d+,?)+\)`)
	if err != nil {
		panic("there is an issue with the regexp for light diagrams")
	}

	scanner := bufio.NewScanner(file)
	var machines []Machine
	for scanner.Scan() {
		line := scanner.Text()

		lightStateString := lightStateRegexp.FindString(line)
		if len(lightStateString) == 0 {
			return nil, errors.New("this line contained no light diagram! line: " + line)
		}
		actionStrings := actionsRegexp.FindAllString(line, -1)
		if actionStrings == nil {
			return nil, errors.New("this line contained no actions! line: " + line)
		}

		desiredState := State(lightStateString[1 : len(lightStateString)-2])

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
		machines = append(machines, Machine{
			DesiredState:    desiredState,
			Actions:         actions,
			Queue:           NewQueue[Record](),
			ProcessedStates: make(map[State]struct{}),
		})
	}
	return machines, nil
}

func (state State) ApplyAction(action Action) (State, error) {
	newStateArr := strings.Split(string(state), "")
	for _, button := range action {
		switch newStateArr[button] {
		case ".":
			newStateArr[button] = "#"
		case "#":
			newStateArr[button] = "."
		default:
			return "", errors.New("the state contains a character other than . or # state:" + string(state))
		}
	}
	newState := strings.Join(newStateArr, "")
	return State(newState), nil
}

func (machine Machine) Solve() (*Record, error) {
	tries := 0
	for machine.Queue.Size() > 0 || tries < 1000 {
		record, err := machine.Queue.Dequeue()
		if err != nil {
			return nil, err
		}
		if record.State == machine.DesiredState {
			return &record, nil
		}
		if _, exists := machine.ProcessedStates[record.State]; exists {
			continue
		}
		for _, action := range machine.Actions {
			newState, err := record.State.ApplyAction(action)
			if err != nil {
				return nil, err
			}
			machine.Queue.Enqueue(Record{
				Steps:    record.Steps + 1,
				PrevStep: &record,
				State:    newState,
			})
		}
		machine.ProcessedStates[record.State] = struct{}{}
		tries++
	}
	return nil, errors.New("failed to find the solution.")
}

func main() {
	machines, err := parseInput()
	if err != nil {
		println(err)
		return
	}
	total := 0
	for _, machine := range machines {
		initialState := State(strings.Repeat(".", len(machine.DesiredState)))
		machine.Queue.Enqueue(Record{
			State:    initialState,
			Steps:    0,
			PrevStep: nil,
		})
		answerRecord, err := machine.Solve()
		if err != nil {
			println("failed to find the solution for the machine with desired state of " + string(machine.DesiredState) + " error: " + err.Error())
			return
		}
		println("found  the solution for the machine with desired state of", machine.DesiredState, " -> ", answerRecord.Steps)
		total += answerRecord.Steps
	}
	println("Total: ", total)
}
