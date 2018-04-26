package scripts

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"syscall"
)

// RunScript represents a collection of JSON configured executable scripts.
type RunScript struct {
	Scripts []script `json:"scripts"`
}

// RunScriptResponse represents the results of each executed script
type RunScriptResponse struct {
	Results []result `json:"results"`
}

type result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	StatusCode int    `json:"status_code"`
}

type script struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func (r *RunScript) Unmarshal(script string) error {
	err := json.Unmarshal([]byte(script), &r)

	if err != nil {
		log.Fatalln("Error Unmarshalling: ", err)
	}

	return err
}

func (r *RunScript) Execute() (*RunScriptResponse, error) {
	results := make([]result, 0)
	for _, x := range r.Scripts {
		r, err := execScript(x)
		if err != nil {
			log.Println("ERROR :" + err.Error())
		}
		results = append(results, r)
	}
	return &RunScriptResponse{results}, nil
}

func execScript(s script) (result, error) {
	cmd := exec.Command(s.Command, s.Args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	r := result{
		stdout.String(),
		stderr.String(),
		-1,
	}
	if err == nil {
		r.StatusCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}
	return r, err
}
