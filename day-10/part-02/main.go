package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const FILENAME = "../sample-input.txt"

type State []int

type StateHash uint64

type Action []int

type Record struct {
	Steps    int
	PrevStep *Record
	State    *State
}

type Machine struct {
	DesiredState    State
	Actions         []Action
	Queue           *Queue[Record]
	ProcessedStates map[StateHash]*State
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
			DesiredState:    desiredState,
			Actions:         actions,
			Queue:           NewQueue[Record](),
			ProcessedStates: make(map[StateHash]*State),
		})
	}
	return machines, nil
}

func (state State) ApplyAction(action Action, multiplier int) (State, error) {
	newState := append(State{}, state...)
	for _, button := range action {
		newState[button] += multiplier
	}
	return newState, nil
}

func (machine Machine) Solve() (*Record, error) {
	tries := 0
	for machine.Queue.Size() > 0 || tries < 1_000_000 {
		record, err := machine.Queue.Dequeue()
		if err != nil {
			return nil, err
		}
		if record.State.Hash() == machine.DesiredState.Hash() {
			return &record, nil
		}
		if _, exists := machine.ProcessedStates[record.State.Hash()]; exists {
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
