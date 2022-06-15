package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputDir string
var outputDir string
var currDirOverrideFlag bool
var subfolderFlag bool
var verbosityFlag bool
var copyFolderStructure bool

var outDirPerm os.FileMode

func sortDir(cwd os.File) {
	contents, err := cwd.ReadDir(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range contents {
		if f.IsDir() && subfolderFlag {
			nwd, err := os.Open(cwd.Name() + "/" + f.Name())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer sortDir(*nwd)
		}
		if !f.IsDir() {
			filename := strings.Split(f.Name(), ".")
			var realFolderPath string
			if copyFolderStructure {
				realFolderPath = strings.Replace(cwd.Name(), inputDir, "", 1)
				realFolderPath = strings.TrimPrefix(realFolderPath, "/")+"/"
			}
			var ext string
			if len(filename) == 1 {
				fi, err := os.Lstat(cwd.Name()+"/"+f.Name())
				if err != nil {
					fmt.Println(err)
					return
				}
				if strings.Count(fi.Mode().String(), "x") >= 1 {
					ext = "EXECUTABLE"
				} else {ext = "NO EXTENSION"}		
				
			} else {ext = filename[len(filename)-1]}

			if verbosityFlag {
				fmt.Println("Moving", f.Name(), "to", outputDir+"/"+ext+"/"+realFolderPath+f.Name())
			}
			os.MkdirAll(outputDir+"/"+ext+"/"+realFolderPath, outDirPerm)
			err := os.Rename(cwd.Name()+"/"+f.Name(), outputDir+"/"+ext+"/"+realFolderPath+"/"+f.Name())
			if err != nil {
				fmt.Println(err)
				return
			}
			
		}
	}
}

func main() {
	flag.BoolVar(&currDirOverrideFlag, "d", false, "(Default) Sorts only files in the working directory. If called overrides any other behaviour flags.")
	flag.BoolVar(&verbosityFlag, "v", false, "Verbose output")
	flag.StringVar(&inputDir, "i", "./", "Directory containing files to be sorted. If -t is set, subfolders will also have their content sorted.")
	flag.BoolVar(&subfolderFlag, "s", false, "Sort also files in subfolders")
	flag.BoolVar(&copyFolderStructure, "t", false, "When set, folder structure gets copied into each extension folder")
	flag.StringVar(&outputDir, "o", "./", "Directory to output the sorted files to.")
	flag.Parse()

	fi, err := os.Lstat(outputDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	outDirPerm = fi.Mode()

	cwd, err := os.Open(inputDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	sortDir(*cwd)

}
