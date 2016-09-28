package utils

import (
    "errors"
    // "fmt"
    "os"
    "path"
    "strings"
)

func IsDir(dir string) bool {
    f, e := os.Stat(dir)
    if e != nil {
        return false
    }
    return f.IsDir()
}

func statDir(dirPath, revPath string, includeDir, isDirOnly bool) ([]string, error) {
    dir, err := os.Open(dirPath)
    if err != nil {
        return nil, err
    }
    defer dir.Close()

    fis, err := dir.Readdir(0)
    if err != nil {
        return nil, err
    }

    // var statList []string
    statList := make([]string, 0)
    for _, fi := range fis {
        if strings.Contains(fi.Name(), ".DS_Store") {
            continue
        }
        relPath := path.Join(revPath, fi.Name())
        curPath := path.Join(dirPath, fi.Name())
        if fi.IsDir() {
            if includeDir {
                statList = append(statList, relPath+"/")
            }
            s, err := statDir(curPath, relPath, includeDir, isDirOnly)
            if err != nil {
                return nil, err
            }
            statList = append(statList, s...)
        } else if !isDirOnly {
            statList = append(statList, relPath)
        }
    }
    return statList, nil
}

func StatDir(rootPath string, includeDir ...bool) ([]string, error) {
    if !IsDir(rootPath) {
        return nil, errors.New("not a directory or does not exist: " + rootPath)
    }
    isIncludeDir := false
    if len(includeDir) >= 1 {
        isIncludeDir = includeDir[0]
    }
    return statDir(rootPath, "", isIncludeDir, false)
}

func GetAllSubDirs(rootPath string) ([]string, error) {
    if !IsDir(rootPath) {
        return nil, errors.New("not a directory or does not exist: " + rootPath)
    }
    return statDir(rootPath, "", true, true)
}
