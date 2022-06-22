package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	oldDir := ""
	fs, err := ioutil.ReadDir(oldDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	var newDir = ""
	for _, f := range fs {
		sp := strings.Split(f.Name(), ".")
		if strings.ToLower(sp[len(sp)-1]) != "md" {
			continue
		}
		var contents []byte
		var filePath = path.Join(oldDir, f.Name())
		var fd, err = os.Open(filePath)
		if err != nil {
			fmt.Println("read file err: ", filePath, " err: ", err)
			continue
		}

		scanner := bufio.NewScanner(fd)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			line := scanner.Bytes()
			//fmt.Println(string(line))
			if len(line) == 0 {
				contents = append(contents, '\n')
				continue
			}
			if line[0] != '#' {
				contents = append(contents, line...)
				contents = append(contents, '\n')
				continue
			}
			var i = 1
			for ; i < len(line); i++ {
				if line[i] != line[i-1] {
					if line[i] == ' ' {
						contents = append(contents, line...)
						goto end
					} else {
						contents = append(contents, line[:i]...)
						contents = append(contents, ' ')
						contents = append(contents, line[i:]...)
						goto end
					}
				}
			}

			contents = append(contents, line...)
		end:
			contents = append(contents, '\n')
		}

		var newFileName = path.Join(newDir, f.Name())

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

		fd.Close()
		fd, err = os.OpenFile(newFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			fmt.Printf("new file name: %s , err: %v\n", newFileName, err)
			continue
		}

		fd.Write(contents)
		fd.Sync()
		fd.Close()
	}
}
