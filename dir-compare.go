package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// struct to hold info for a file
type File struct {
	Name         string
	OriginalSize int64
	NewSize      int64
	Difference   int64
}

// slice of files with custom sorter
type BySize []File

func (files BySize) Len() int {
	return len(files)
}

func (files BySize) Swap(i int, j int) {
	files[i], files[j] = files[j], files[i]
}

func (files BySize) Less(i int, j int) bool {
	return files[i].OriginalSize < files[j].OriginalSize
}

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("not enough arguments suplied")
	}

	// check if src and dest exist
	info, err := os.Stat(args[1])
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("the directory %s does not exist", args[1])
		}
	}
	if !info.IsDir() {
		log.Fatal("first argument is not a directory")
	}

	info, err = os.Stat(args[2])
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("the directory %s does not exist", args[1])
		}
	}
	if !info.IsDir() {
		log.Fatal("second argument is not a directory")
	}

	// read files in src
	filenames, err := os.ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}

	fileDict := make(map[string]File)
	for _, file := range filenames {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		if info.Mode().IsRegular() {
			fileDict[info.Name()] = File{
				Name:         info.Name(),
				OriginalSize: info.Size(),
				NewSize:      0,
				Difference:   0.0,
			}
		}
	}

	// get sizes for files in dest if they exist
	filenames, err = os.ReadDir(args[2])
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range filenames {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		if info.Mode().IsRegular() {
			if _, exists := fileDict[info.Name()]; exists {
				file := fileDict[info.Name()]
				file.NewSize = info.Size()
				file.Difference = file.OriginalSize - file.NewSize
				fileDict[info.Name()] = file
			}
		}
	}

	// sort files by original size
	files := []File{}
	for file := range fileDict {
		files = append(files, fileDict[file])
	}
	sort.Sort(BySize(files))

	// build and display table
	filesTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#2DE224"))).
		Headers("File", "Orig. Size", "New Size", "Saved Space").
		StyleFunc(func(row int, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})

	var totalOrigSize int64
	var totalNewSize int64
	for _, file := range files {
		if file.NewSize != 0.0 {
			filesTable.Row(file.Name, convertBytes(file.OriginalSize), convertBytes(file.NewSize), convertBytes(file.Difference))
			totalOrigSize += file.OriginalSize
			totalNewSize += file.NewSize
		}
	}
	totalDifference := totalOrigSize - totalNewSize
	totalDifferencePercent := -(float64(totalNewSize) - float64(totalOrigSize)) / float64(totalOrigSize) * 100
	filesTable.Row("Total", convertBytes(totalOrigSize), convertBytes(totalNewSize), convertBytes(totalDifference)+fmt.Sprintf(" (%.1f %%)", totalDifferencePercent))

	fmt.Println(filesTable.Render())
}

// converts the given number of bytes to the closest representation as KiB, MiB, GiB, ...
func convertBytes(numBytes int64) string {
	representation := float64(numBytes)
	counter := 0
	for math.Abs(representation) > 1024.0 {
		representation /= 1024.0
		counter += 1
	}
	switch counter {
	case 0:
		return fmt.Sprintf("%.1f B", representation)
	case 1:
		return fmt.Sprintf("%.1f KiB", representation)
	case 2:
		return fmt.Sprintf("%.1f MiB", representation)
	case 3:
		return fmt.Sprintf("%.1f GiB", representation)
	case 4:
		return fmt.Sprintf("%.1f TiB", representation)
	case 5:
		return fmt.Sprintf("%.1f PiB", representation)
	case 6:
		return fmt.Sprintf("%.1f EiB", representation)
	}

	return "N/A"
}
