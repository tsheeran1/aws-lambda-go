package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// LProject contains the json from the project.json file
type LProject struct {
	Name   string
	Bucket string
	Role   string
	path   string
}

//NewLProject takes a project.json file name and reads the json into LProject
func NewLProject(fname string) (LProject, error) {

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return LProject{}, err
	}

	var res LProject
	err = json.Unmarshal(data, &res)
	if err != nil {
		return LProject{}, err
	}
	res.path = path.Dir(fname)

	if strings.HasPrefix(res.Role, "arn") {
		return res, nil //Happy
	}

	rmp, err := RoleMap()
	if err != nil {
		return res, err
	}

	nRole, ok := rmp[res.Role]
	res.Role = nRole
	if !ok {
		return res, errors.New("Role Not found: " + res.Role)
	}

	return res, nil

}

//UploadLambda takes a file name and uploads it as a lambda function
func (lp LProject) UploadLambda(name string) error {
	fpath := path.Join(lp.path, name)

	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")

	fmt.Println("Building : " + fpath + ".go")
	_, err := run("go", "build", "-o", fpath, fpath+".go")
	if err != nil {
		return err
	}

	fmt.Println("Zipping to : " + fpath + ".zip")
	_, err = run("zip", "-j", fpath+".zip", fpath)
	if err != nil {
		return err
	}

	lamname := lp.Name + "_" + name

	upcmd := exec.Command("aws", "s3", "cp", fpath+".zip", "s3://"+lp.Bucket+"/"+lamname+".zip")

	upOut, err := upcmd.StdoutPipe()
	if err != nil {
		return err
	}

	fmt.Println("Starting Upload of : " + lamname)
	err = upcmd.Start()
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, upOut)
	err = upcmd.Wait()
	if err != nil {
		return err
	}

	fl, err := NewFunctionList()
	if err != nil {
		return err
	}

	if fl.HasFunction(lamname) {
		fmt.Println("Updating Function : " + lamname)
		resp, err := run("aws", "lambda", "update-function-code", "--function-name", lamname, "--s3-bucket", lp.Bucket, "--s3-key", lamname+".zip")
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	}

	fmt.Println("Creating Function")
	resp, err := run("aws", "lambda", "create-function", "--function-name", lamname, "--runtime", "go1.x", "--role", lp.Role, "--handler", name, "--code", "S3Bucket="+lp.Bucket+", S3Key="+lamname+".zip")
	if err != nil {
		return err
	}

	fmt.Println(string(resp))
	return nil
}

func main() {

	lname := flag.String("n", "", "Name of Lambda")
	confloc := flag.String("c", "project.json", "Location of Config file")
	flag.Parse()

	proj, err := NewLProject(*confloc)
	if err != nil {
		log.Fatal(err)
	}

	err = proj.UploadLambda(*lname)
	if err != nil {
		log.Fatal(err)
	}

}
