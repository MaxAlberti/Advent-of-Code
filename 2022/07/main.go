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

type Command struct {
	Name string
	Arg  string
}

type FileSysObj struct {
	Name     string
	Children []*FileSysObj
	Parent   *FileSysObj
	Size     int
	IsFile   bool
}

type OutputType int

const (
	FileSystemType = 1
	CommandType    = 2
)

type FileSystem map[string]FileSysObj

// Main function
func main() {
	fmt.Println("Starting")

	lines, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fs, err := parse_input(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	calc_p1(fs)
	calc_p2(fs)
}

func calc_p1(fs FileSystem) {
	dir_sizes := get_all_fs_dir_size(fs)
	res := 0
	for _, size := range dir_sizes {
		if size <= 100000 {
			res += size
		}
	}

	fmt.Printf("Part1: %d\n", res)
}

func calc_p2(fs FileSystem) {
	dir_sizes := get_all_fs_dir_size(fs)
	used_space := dir_sizes["/"]
	total_space := 70000000
	needed_space := 30000000
	free_space := total_space - used_space
	min_remove := needed_space - free_space

	// Find folder to delete
	res_name := ""
	res_size := total_space
	for key, size := range dir_sizes {
		if size >= min_remove && size < res_size {
			res_name = key
			res_size = size
		}
	}
	fmt.Printf("Part2: Remove dir '%s' wit a size of %d. (Min %d)\n", res_name, res_size, min_remove)
}

func parse_input(lines []string) (FileSystem, error) {
	// Create filesystem
	file_system := make(FileSystem)

	// Add root
	file_system["/"] = FileSysObj{Name: "/", Children: []*FileSysObj{}, Parent: nil, Size: 0, IsFile: false}

	var current_dir *FileSysObj
	rf := file_system["/"]
	current_dir = &rf

	for index, line := range lines {
		obj, cmd, itype, err := parse_line(line)
		if err != nil {
			return FileSystem{}, err
		}
		if index == 0 && itype != CommandType {
			return FileSystem{}, errors.New("first imput must be a command")
		}
		if itype == FileSystemType {
			current_dir, err = parse_fso(obj, current_dir, file_system)
			if err != nil {
				return FileSystem{}, err
			}
		} else if itype == CommandType {
			current_dir, err = parse_cmd(cmd, current_dir, file_system)
			if err != nil {
				return FileSystem{}, err
			}
		}
	}

	return file_system, nil
}

func parse_fso(obj FileSysObj, cur *FileSysObj, fs FileSystem) (*FileSysObj, error) {
	// Set the parent dir
	obj.Parent = cur
	// Get fso path
	fso_path := get_fso_path(&obj)
	// Add to file_system
	_, ex := fs[fso_path]
	if ex {
		// obj exists allready
		return cur, nil //error?
	} else {
		fs[fso_path] = obj
	}
	// Add obj to cur's children
	fso := fs[fso_path]
	cur.Children = append(cur.Children, &fso)

	// Update cur
	cur_path := get_fso_path(cur)
	fs[cur_path] = *cur

	return cur, nil
}

func parse_cmd(cmd Command, cur *FileSysObj, fs FileSystem) (*FileSysObj, error) {
	if cmd.Name == "cd" {
		if cmd.Arg == "/" {
			// Get root pointer
			rf := fs["/"]
			return &rf, nil
		} else if cmd.Arg == ".." {
			if cur.Parent != nil {
				return cur.Parent, nil
			} else {
				return cur, errors.New("cd error, dir has no parrent")
			}
		} else {
			for _, child := range cur.Children {
				if child.Name == cmd.Arg {
					return child, nil
				}
			}

			return cur, errors.New("cd destination doesn't exist. (ls should run bevore cd, this should not happen)")
		}
	} else if cmd.Name == "ls" {
		// pass
		return cur, nil
	} else {
		return cur, errors.New("unknown command")
	}
}

func parse_line(line string) (FileSysObj, Command, OutputType, error) {
	var re_cmd = regexp.MustCompile(`\$ (?P<name>cd|ls) ?(?P<arg>.+)?`)
	var re_dir = regexp.MustCompile(`dir (?P<name>.+)`)
	var re_file = regexp.MustCompile(`(?P<size>[0-9]+) (?P<name>.+)`)

	obj := FileSysObj{}
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
		obj.Name = res["name"]
		obj.Children = []*FileSysObj{}
		obj.Parent = nil
		obj.Size = 0
		obj.IsFile = false
		otype = FileSystemType
	} else if len(match_file) > 0 && len(match_dir) == 0 && len(match_cmd) == 0 {
		// Is file
		res := get_regex_goup_map(re_file, match_file)
		obj.Name = res["name"]
		obj.Children = []*FileSysObj{}
		obj.Parent = nil
		obj.Size, _ = strconv.Atoi(res["size"])
		obj.IsFile = true
		otype = FileSystemType
	} else {
		return obj, cmd, otype, errors.New("can't parse line '" + line + "' multiple matches")
	}

	return obj, cmd, otype, nil
}

func get_fso_path(obj *FileSysObj) string {
	if obj.Parent == nil && obj.Name == "/" {
		return "/"
	}

	path := obj.Name
	cur := obj
	for cur.Parent != nil {
		path = cur.Parent.Name + "/" + path
		cur = cur.Parent
	}

	return path[1:]
}

func get_fso_size(obj *FileSysObj) int {
	if obj.IsFile {
		return obj.Size
	} else {
		res := 0
		for _, c := range obj.Children {
			res += get_fso_size(c)
		}
		return res
	}
}

func get_all_fs_dir_size(fs FileSystem) map[string]int {
	res := make(map[string]int)

	for key, obj := range fs {
		if !obj.IsFile {
			res[key] = get_fso_size(&obj)
		}
	}

	return res
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
