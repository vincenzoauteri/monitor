package main

import (
  "os"
  "fmt"
  "time"
  "regexp"
  "syscall"
)

var times map[string] int64

func parseDir (f *os.File) bool {
  result  := false
  dir, err := f.Readdir(100)
  if err != nil {
    return result
  }

  for i := range dir {
    if dir[i].Mode().IsDir() && dir[i].Name()[0] != '.'{
      os.Chdir(dir[i].Name())
      dirFile, err := os.Open(".")
      if (err == nil) {
        result = parseDir(dirFile)
        os.Chdir("..")
      } else {
        fmt.Println("Error")
      }
    } else {
      if matched, _ := regexp.MatchString(".go$",dir[i].Name()) ; matched { 
        path, _ := os.Getwd()
        key := path + "/" + dir[i].Name()
        value, none :=  times[key]
        currentModTime := dir[i].ModTime().Unix()
        if none {
          if currentModTime > value {
            fmt.Println("Modified: " + dir[i].Name())
            times[key] = currentModTime
            result = true
          }
        } else {
          times[key] = currentModTime
          result = true
        }
      }
    }
  }
  return result
}


const goDir = "/usr/local/bin/"
const appBins = "/Users/enzo/workspace/go/bin/"

func compileAndStart (targetDir string, appName string)  (*os.Process, error) {
  var procAttr os.ProcAttr
  procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
  args  := []string {"go", "install", targetDir + appName}
  p,err := os.StartProcess(goDir + "go", args, &procAttr)
  if err != nil {
    fmt.Println(err.Error())
    return p, err
  }
  p, err  = os.StartProcess(appBins + appName ,nil,&procAttr)
  return p, err
}


func main() {

  appName := "gotest"
  targetDir := "./"

  serverPid, err := compileAndStart(targetDir,appName)
  if err != nil {
    fmt.Println(err.Error())
  }
  times = make(map[string] int64)
  dir,err := os.Open(targetDir)
  if err != nil {
    fmt.Println(err)
    return
  }
  _ = parseDir(dir)
  for {
    dir,_ := os.Open(targetDir)
    mod := parseDir(dir)
    if mod {
      _ = serverPid.Signal(syscall.SIGKILL)
      serverPid, err= compileAndStart(targetDir,appName)
    }
    time.Sleep(1 * time.Second)
  }
}


