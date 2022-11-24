package fsm

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"github.com/eduvpn/eduvpn-common/types"
)

type (
	FSMStateID      int8
	FSMStateIDSlice []FSMStateID
)

func (v FSMStateIDSlice) Len() int {
	return len(v)
}

func (v FSMStateIDSlice) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v FSMStateIDSlice) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type FSMTransition struct {
	To          FSMStateID
	Description string
}

type (
	FSMStates map[FSMStateID]FSMState
)

type FSMState struct {
	Transitions []FSMTransition

	// Which state to go back to on a back transition
	BackState FSMStateID
}

type FSM struct {
	States  FSMStates
	Current FSMStateID

	// Info to be passed from the parent state
	Name          string
	StateCallback func(FSMStateID, FSMStateID, interface{}) bool
	Directory     string
	Debug         bool
	GetName       func(FSMStateID) string
}

func (fsm *FSM) Init(
	current FSMStateID,
	states map[FSMStateID]FSMState,
	callback func(FSMStateID, FSMStateID, interface{}) bool,
	directory string,
	nameGen func(FSMStateID) string,
	debug bool,
) {
	fsm.States = states
	fsm.Current = current
	fsm.StateCallback = callback
	fsm.Directory = directory
	fsm.GetName = nameGen
	fsm.Debug = debug
}

func (fsm *FSM) InState(check FSMStateID) bool {
	return check == fsm.Current
}

func (fsm *FSM) HasTransition(check FSMStateID) bool {
	for _, transitionState := range fsm.States[fsm.Current].Transitions {
		if transitionState.To == check {
			return true
		}
	}

	return false
}

func (fsm *FSM) getGraphFilename(extension string) string {
	debugPath := path.Join(fsm.Directory, "graph")
	return fmt.Sprintf("%s%s", debugPath, extension)
}

func (fsm *FSM) writeGraph() {
	graph := fsm.GenerateGraph()
	graphFile := fsm.getGraphFilename(".graph")
	graphImgFile := fsm.getGraphFilename(".png")
	f, err := os.Create(graphFile)
	if err != nil {
		return
	}

	_, writeErr := f.WriteString(graph)
	f.Close()
	if writeErr != nil {
		cmd := exec.Command("mmdc", "-i", graphFile, "-o", graphImgFile, "--scale", "4")
		// Generating is best effort
		_ = cmd.Start()
	}
}

func (fsm *FSM) GoBack() {
	fsm.GoTransition(fsm.States[fsm.Current].BackState)
}

func (fsm *FSM) GoTransitionRequired(newState FSMStateID, data interface{}) error {
	oldState := fsm.Current
	if !fsm.GoTransitionWithData(newState, data) {
		return types.NewWrappedError("failed required transition", fmt.Errorf("required transition not handled, from: %s -> to: %s", fsm.GetName(oldState), fsm.GetName(newState)))
	}
	return nil
}

func (fsm *FSM) GoTransitionWithData(newState FSMStateID, data interface{}) bool {
	ok := fsm.HasTransition(newState)

	handled := false
	if ok {
		oldState := fsm.Current
		fsm.Current = newState
		if fsm.Debug {
			fsm.writeGraph()
		}

		handled = fsm.StateCallback(oldState, newState, data)
	}

	return handled
}

func (fsm *FSM) GoTransition(newState FSMStateID) bool {
	// No data means the callback is never required
	return fsm.GoTransitionWithData(newState, "")
}

func (fsm *FSM) generateMermaidGraph() string {
	graph := "graph TD\n"
	sortedFSM := make(FSMStateIDSlice, 0, len(fsm.States))
	for stateID := range fsm.States {
		sortedFSM = append(sortedFSM, stateID)
	}
	sort.Sort(sortedFSM)
	for _, state := range sortedFSM {
		transitions := fsm.States[state].Transitions
		for _, transition := range transitions {
			if state == fsm.Current {
				graph += "\nstyle " + fsm.GetName(state) + " fill:cyan\n"
			} else {
				graph += "\nstyle " + fsm.GetName(state) + " fill:white\n"
			}
			graph += fsm.GetName(
				state,
			) + "(" + fsm.GetName(
				state,
			) + ") " + "-->|" + transition.Description + "| " + fsm.GetName(
				transition.To,
			) + "\n"
		}
	}
	return graph
}

func (fsm *FSM) GenerateGraph() string {
	if fsm.GetName != nil {
		return fsm.generateMermaidGraph()
	}

	return ""
}
