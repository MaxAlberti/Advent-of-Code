// First Go program
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Directory struct {
	Name     string
	Files    []File
	Children []Directory
	Parent   *Directory
}

type File struct {
	Name string
	Size int
}

type Command struct {
	Name string
	Arg  string
}

type OutputType int

const (
	DirectoryType = 0
	FileType      = 1
	CommandType   = 2
)

// Main function
func main() {
	fmt.Println("Starting")

	lines, err := get_file_lines("nav_test.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	_, err = parse_input(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func parse_input(lines []string) (Directory, error) {
	var root Directory
	root.Name = "/"
	root.Children = []Directory{}
	root.Files = []File{}
	root.Parent = nil

	var current_dir *Directory = nil
	for index, line := range lines {
		dir, file, cmd, itype, err := parse_line(line)
		if err != nil {
			return Directory{}, err
		}
		if index == 0 && itype != CommandType {
			return Directory{}, errors.New("first imput must be a command")
		}
		if itype == DirectoryType {
			current_dir, err = parse_dir(dir, current_dir)
			if err != nil {
				return Directory{}, err
			}
		} else if itype == FileType {
			_ = file.Name
		} else if itype == CommandType {
			current_dir, err = parse_cmd(cmd, &root, current_dir)
			if err != nil {
				return Directory{}, err
			}
		}
	}

	return root, nil
}

func parse_dir(dir Directory, cur *Directory) (*Directory, error) {
	for _, child := range cur.Children {
		if child.Name == dir.Name {
			// dir allready exists
			return cur, nil // error?
		}
	}

	/* Works, but can't go up in filesystem anymore bc pointer is changed
	tst := *cur
	dir.Parent = &tst
	tst.Children = append(tst.Children, dir)
	*cur = tst
	*/
	cur.Children = append(cur.Children, Directory{Parent: cur, Name: dir.Name})

	//dir.Parent = cur
	//cur.Children = append(cur.Children, dir)
	//cur.Children = append(cur.Children, Directory{Parent: cur, Name: dir.Name})

	return cur, nil
}

func parse_cmd(cmd Command, root *Directory, cur *Directory) (*Directory, error) {
	if cmd.Name == "cd" {
		if cmd.Arg == "/" {
			return root, nil
		} else if cmd.Arg == ".." {
			if cur.Parent != nil {
				return cur.Parent, nil
			} else {
				return cur, errors.New("cd error, dir has no parrent")
			}
		} else {
			for _, child := range cur.Children {
				if child.Name == cmd.Arg {
					return &child, nil
				}
			}

			/* ls should run bevore cd
			cur.Children = append(cur.Children, Directory{Parent: cur, Name: cmd.Arg})

			return &cur.Children[len(cur.Children)-1], nil
			*/
			return cur, errors.New("cd destination doesn't exist. (ls should run bevore cd, this should not happen)")
		}
	} else if cmd.Name == "ls" {
		// pass
		return cur, nil
	} else {
		return cur, errors.New("unknown command")
	}
}

func parse_line(line string) (Directory, File, Command, OutputType, error) {
	var re_cmd = regexp.MustCompile(`\$ (?P<name>cd|ls) ?(?P<arg>.+)?`)
	var re_dir = regexp.MustCompile(`dir (?P<name>.+)`)
	var re_file = regexp.MustCompile(`(?P<size>[0-9]+) (?P<name>.+)`)

	dir := Directory{}
	file := File{}
	cmd := Command{}
	var otype OutputType = -1

	match_cmd := re_cmd.FindStringSubmatch(line)
	match_dir := re_dir.FindStringSubmatch(line)
	match_file := re_file.FindStringSubmatch(line)

	if len(match_cmd) > 0 && len(match_dir) == 0 && len(match_file) == 0 {
		// Is cmd
		res := get_regex_goup_map(re_cmd, match_cmd)
		cmd.Name = res["name"]
		arg, ex := res["arg"]
		if ex {
			cmd.Arg = arg
		} else {
			cmd.Arg = ""
		}
		otype = CommandType
	} else if len(match_dir) > 0 && len(match_cmd) == 0 && len(match_file) == 0 {
		// Is dir
		res := get_regex_goup_map(re_dir, match_dir)
		dir.Name = res["name"]
		dir.Children = []Directory{}
		dir.Files = []File{}
		otype = DirectoryType
	} else if len(match_file) > 0 && len(match_dir) == 0 && len(match_cmd) == 0 {
		// Is file
		res := get_regex_goup_map(re_file, match_file)
		file.Name = res["name"]
		file.Size, _ = strconv.Atoi(res["size"])
		otype = FileType
	} else {
		return dir, file, cmd, otype, errors.New("can't parse line '" + line + "' multiple matches")
	}

	return dir, file, cmd, otype, nil
}

func get_regex_goup_map(re *regexp.Regexp, match []string) map[string]string {
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

func get_file_lines(filepath string) ([]string, error) {
	// open file
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	arr := []string{}
	for scanner.Scan() {
		// Read the line
		arr = append(arr, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	} else {
		return arr, err
	}
}
