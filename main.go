package main

import (
	"fmt"
	"github.com/str1ngs/util/file"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

var (
	files = map[string]string{}
	dir   = "."
)

func init() {
	log.SetPrefix("testit: ")
	log.SetFlags(log.Lshortfile)
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	tick := time.Tick(time.Second)
	for _ = range tick {
		dirty, err := update_files()
		if err != nil {
			fmt.Println(err)
		}
		if dirty {
			doTests()
		}
	}
}

func doTests() {
	pname, _ := os.Getwd()
	pname = fmt.Sprintf("%s.test", path.Base(pname))
	exec.Command("killall", pname).Run()
	gotest := exec.Command("go", "test")
	gotest.Stderr = os.Stderr
	gotest.Stdout = os.Stdout
	if err := gotest.Start(); err != nil {
		log.Println(err)
	}
}

func update_files() (changed bool, err error) {
	markFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path[:1] == "." || info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		hash, err := file.Md5(path)
		if err != nil {
			return err
		}
		if _, exists := files[path]; !exists {
			changed = true
			fmt.Println("adding", hash, path)
			files[path] = hash
			return nil
		}
		if files[path] != hash {
			changed = true
			fmt.Println("changed", path)
			files[path] = hash
		}
		return nil
	}
	if file.Exists(".testit") {
		doTests()
		os.Remove(".testit")
	}
	return changed, filepath.Walk(".", markFn)
}
