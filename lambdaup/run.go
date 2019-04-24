package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func run(prog string, args ...string) ([]byte, error) {
	cmd := exec.Command(prog, args...)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return []byte{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []byte{}, err
	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	ch := make(chan bool)
	go func() {
		io.Copy(&outBuf, outPipe)
		ch <- true
	}()
	io.Copy(&errBuf, errPipe)
	_ = <-ch

	err = cmd.Wait()
	if err != nil {
		return outBuf.Bytes(), err
	}

	if len(errBuf.Bytes()) != 0 {
		return outBuf.Bytes(), fmt.Errorf("%s", errBuf.Bytes())
	}

	return outBuf.Bytes(), nil
}

// RoleList is the list of roles returned from the AWS iam list-roles command
type RoleList struct {
	Roles []Role
}

//A Role is an individual role returned in the AWS RoleList
type Role struct {
	RoleName string
	Arn      string
}

//RoleMap runs the AWS iam list-roles command, and parses the result into a RoleList which is simply a map of AWS role names and Arn's)
func RoleMap() (map[string]string, error) {
	res := make(map[string]string)
	data, err := run("aws", "iam", "list-roles")
	if err != nil {
		return res, err
	}

	var rlist RoleList
	err = json.Unmarshal(data, &rlist)
	if err != nil {
		return res, err
	}

	for _, v := range rlist.Roles {
		res[v.RoleName] = v.Arn
	}

	return res, nil
}
