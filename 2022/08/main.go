// First Go program
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Main function
func main() {
	fmt.Println("Starting")

	lines, err := get_file_lines("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	tree_grid, y_dim, x_dim, err := create_tree_grid(lines)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	vis_trees, _ := count_visible_trees(tree_grid, y_dim, x_dim)

	fmt.Printf("P1: There are %d visible trees!\n", vis_trees)
}

func count_visible_trees(tree_grid [][]int, y_dim uint, x_dim uint) (int, [][]bool) {
	// Create memorization matrixes
	highest_left := create_matrix_of_size(y_dim, x_dim, -1)
	highest_right := create_matrix_of_size(y_dim, x_dim, -1)
	highest_top := create_matrix_of_size(y_dim, x_dim, -1)
	highest_bottom := create_matrix_of_size(y_dim, x_dim, -1)
	vis_matrix := make([][]bool, y_dim)
	vis_count := 0

	// Iterate over tree_grid
	for y := 0; y < int(y_dim); y++ {
		vis_matrix[y] = make([]bool, x_dim)
		for x := 0; x < int(x_dim); x++ {
			height := tree_grid[y][x]
			// Check from left
			v_left := false // Smaller is default
			top_height := calc_left_vis(&tree_grid, &highest_left, y_dim, x_dim, y, x)
			if x == 0 {
				v_left = true
			} else if height == top_height {
				// Check if previous tree was the same height
				prev_top_height := calc_left_vis(&tree_grid, &highest_left, y_dim, x_dim, y, x-1)
				v_left = prev_top_height < top_height
			}

			// Check from right
			v_right := false // Smaller is default
			top_height = calc_right_vis(&tree_grid, &highest_right, y_dim, x_dim, y, x)
			if x == int(x_dim)-1 {
				v_right = true
			} else if height == top_height {
				// Check if previous tree was the same height
				prev_top_height := calc_right_vis(&tree_grid, &highest_right, y_dim, x_dim, y, x+1)
				v_right = prev_top_height < top_height
			}

			// Check from top
			v_top := false // Smaller is default
			top_height = calc_top_vis(&tree_grid, &highest_top, y_dim, x_dim, y, x)
			if y == 0 {
				v_top = true
			} else if height == top_height {
				// Check if previous tree was the same height
				prev_top_height := calc_top_vis(&tree_grid, &highest_top, y_dim, x_dim, y-1, x)
				v_top = prev_top_height < top_height
			}

			// Check from bottom
			v_bottom := false // Smaller is default
			top_height = calc_bottom_vis(&tree_grid, &highest_bottom, y_dim, x_dim, y, x)
			if y == int(y_dim)-1 {
				v_bottom = true
			} else if height == top_height {
				// Check if previous tree was the same height
				prev_top_height := calc_bottom_vis(&tree_grid, &highest_bottom, y_dim, x_dim, y+1, x)
				v_bottom = prev_top_height < top_height
			}

			// Check from all sides
			vis_matrix[y][x] = v_left || v_right || v_top || v_bottom
			if vis_matrix[y][x] {
				vis_count += 1
			}
		}
	}

	return vis_count, vis_matrix
}

func calc_left_vis(tree_grid *[][]int, vis_arr *[][]int, y_dim uint, x_dim uint, y_pos int, x_pos int) int {
	if (*vis_arr)[y_pos][x_pos] == -1 {
		if x_pos == 0 {
			(*vis_arr)[y_pos][x_pos] = (*tree_grid)[y_pos][x_pos]
		} else {
			(*vis_arr)[y_pos][x_pos] =
				max(
					calc_left_vis(tree_grid, vis_arr, y_dim, x_dim, y_pos, x_pos-1),
					(*tree_grid)[y_pos][x_pos])
		}
	}

	return (*vis_arr)[y_pos][x_pos]
}

func calc_right_vis(tree_grid *[][]int, vis_arr *[][]int, y_dim uint, x_dim uint, y_pos int, x_pos int) int {
	if (*vis_arr)[y_pos][x_pos] == -1 {
		if x_pos == int(x_dim)-1 {
			(*vis_arr)[y_pos][x_pos] = (*tree_grid)[y_pos][x_pos]
		} else {
			(*vis_arr)[y_pos][x_pos] =
				max(
					calc_right_vis(tree_grid, vis_arr, y_dim, x_dim, y_pos, x_pos+1),
					(*tree_grid)[y_pos][x_pos])
		}
	}

	return (*vis_arr)[y_pos][x_pos]
}

func calc_top_vis(tree_grid *[][]int, vis_arr *[][]int, y_dim uint, x_dim uint, y_pos int, x_pos int) int {
	if (*vis_arr)[y_pos][x_pos] == -1 {
		if y_pos == 0 {
			(*vis_arr)[y_pos][x_pos] = (*tree_grid)[y_pos][x_pos]
		} else {
			(*vis_arr)[y_pos][x_pos] =
				max(
					calc_top_vis(tree_grid, vis_arr, y_dim, x_dim, y_pos-1, x_pos),
					(*tree_grid)[y_pos][x_pos])
		}
	}

	return (*vis_arr)[y_pos][x_pos]
}

func calc_bottom_vis(tree_grid *[][]int, vis_arr *[][]int, y_dim uint, x_dim uint, y_pos int, x_pos int) int {
	if (*vis_arr)[y_pos][x_pos] == -1 {
		if y_pos == int(y_dim)-1 {
			(*vis_arr)[y_pos][x_pos] = (*tree_grid)[y_pos][x_pos]
		} else {
			(*vis_arr)[y_pos][x_pos] =
				max(
					calc_bottom_vis(tree_grid, vis_arr, y_dim, x_dim, y_pos+1, x_pos),
					(*tree_grid)[y_pos][x_pos])
		}
	}

	return (*vis_arr)[y_pos][x_pos]
}

func create_tree_grid(lines []string) ([][]int, uint, uint, error) {
	var x_dim uint = 0                // Left to right
	var y_dim uint = uint(len(lines)) // Top to bottom
	var tree_grid [][]int

	for y, line := range lines {
		if x_dim == 0 {
			x_dim = uint(len(line))
		} else if int(x_dim) != len(line) {
			return [][]int{}, 0, 0, errors.New("can't parse input! lines are of different lenght")
		}
		if tree_grid == nil {
			tree_grid = create_matrix_of_size(uint(y_dim), uint(x_dim), -1)
		}
		for x, char := range line {
			num, err := strconv.Atoi(string(char))
			if err != nil {
				return [][]int{}, 0, 0, fmt.Errorf("can't parse input cahr at row:%d col:%d\nreason: %s", y, x, err)
			}
			tree_grid[y][x] = num
		}
	}

	return tree_grid, y_dim, x_dim, nil
}

func create_matrix_of_size(y_dim uint, x_dim uint, init_val int) [][]int {
	matrix := make([][]int, y_dim)

	for y := 0; y < int(y_dim); y++ {
		matrix[y] = make([]int, x_dim)
		for x := 0; x < int(x_dim); x++ {
			matrix[y][x] = init_val
		}
	}

	return matrix
}

func max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
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
