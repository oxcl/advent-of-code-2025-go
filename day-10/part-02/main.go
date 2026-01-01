package main

import (
	"bufio"
	"errors"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const FILENAME = "../input.txt"

type State []int

type StateHash uint64

type Action []int

type Record struct {
	Steps       int
	multipliers []int
	State       *State
}

type Machine struct {
	DesiredState    State
	Actions         []Action
	Queue           *CostQueue[Record]
	ProcessedStates map[StateHash]*Record
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
			ProcessedStates: make(map[StateHash]*Record),
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

func (machine Machine) CalculateHeuristic(steps int, newState *State) float64 {
	maxDistanceToTarget := 0
	minDistanceToTarget := 0
	averageDistanceOfJoltages := 0.0
	totalDistance := 0
	for index, desiredJoltage := range machine.DesiredState {
		distance := desiredJoltage - (*newState)[index]
		averageDistanceOfJoltages += float64(distance) / float64(len(machine.DesiredState))
		totalDistance += distance
		if distance != 0 && distance > maxDistanceToTarget {
			maxDistanceToTarget = distance
		}
		if distance != 0 && distance < minDistanceToTarget {
			minDistanceToTarget = distance
		}
	}
	// prefer states that their joltages on average have a closer distance to targets
	// return float64(steps) + averageDistanceOfJoltages
	// return float64(steps) + float64(minDistanceToTarget)
	return float64(steps) + float64(totalDistance)
	// return float64(steps)
	// return float64(steps) + float64(maxDistanceToTarget)

}

func (machine Machine) Solve() (*Record, int, error) {
	tries := 0
	for machine.Queue.Size() > 0 && tries < 2_000_000 {
		record, err := machine.Queue.Pop()
		if err != nil {
			return nil, 0, err
		}
		if record.State.Hash() == machine.DesiredState.Hash() {
			return record, tries, nil
		}
		isActionInvalid := false
		for index, newJoltage := range *record.State {
			if newJoltage > machine.DesiredState[index] {
				isActionInvalid = true
				break
			}
		}
		if isActionInvalid {
			continue
		}
		if cachedRecord, exists := machine.ProcessedStates[record.State.Hash()]; exists {
			if cachedRecord.Steps <= record.Steps {
				continue
			} else {
				machine.ProcessedStates[record.State.Hash()] = record
			}
		}
		for index, action := range machine.Actions {
			minDistanceToTarget := math.MaxInt
			for _, index := range action {
				distance := machine.DesiredState[index] - (*record.State)[index]
				if distance < minDistanceToTarget {
					minDistanceToTarget = distance
				}
			}
			for multiplier := minDistanceToTarget; multiplier >= 1; multiplier-- {
				newState, err := record.State.ApplyAction(action, multiplier)
				if err != nil {
					return nil, 0, err
				}
				cost := machine.CalculateHeuristic(record.Steps+multiplier, &newState)
				newMultipliers := append([]int{}, record.multipliers...)
				newMultipliers[index] += multiplier
				machine.Queue.Add(&Record{
					Steps:       record.Steps + multiplier,
					multipliers: newMultipliers,
					State:       &newState,
				}, cost)
			}
		}
		machine.ProcessedStates[record.State.Hash()] = record
		tries++
	}
	return nil, 0, errors.New("failed to find the solution.")
}

func main() {
	machines, err := parseInput()
	if err != nil {
		println(err)
		return
	}
	total := 0
	for index, machine := range machines {
		initialState := make(State, len(machine.DesiredState))
		cost := machine.CalculateHeuristic(0, &initialState)
		machine.Queue.Add(&Record{
			State:       &initialState,
			Steps:       0,
			multipliers: make([]int, len(machine.Actions)),
		}, cost)
		answerRecord, tries, err := machine.Solve()
		if err != nil {
			println("failed to find the solution for the machine with desired state of index " + strconv.Itoa(index) + " error: " + err.Error())
			continue
		}
		println("found  the solution for the machine with desired state of", strconv.Itoa(index), " -> ", answerRecord.Steps, " in ", tries, "tries")
		total += answerRecord.Steps
	}
	println("Total: ", total)
}
