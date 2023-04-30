package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	utils "nabbit/utils"

	"github.com/spf13/pflag"
)

func main() {

	// Create command line arg flags

	// Official pflag package has not fixed issues#352 - Need this
	flag := pflag.NewFlagSet("flags", pflag.ContinueOnError)

	help := flag.BoolP("help", "h", false, "Display help/usage information")
	usage := flag.BoolP("usage", "u", false, "Display help/usage information")
	version := flag.BoolP("version", "v", false, "Display version")
	createFile := flag.StringP("create-file", "c", "", "Choose bookmark file type for creation from .txt file (chrome/firefox)")
	output := flag.StringP("output", "o", "", "Output filename of bookmark file (used w/ create-file)")
	nb := flag.BoolP("num-bookmarks", "b", false, "Get total number of bookmarks")
	nf := flag.BoolP("num-folders", "f", false, "Get total number of folders")
	ns := flag.BoolP("num-separators", "s", false, "Get total number of separators (Firefox only)")
	lb := flag.BoolP("list-bookmarks", "l", false, "List all of the bookmarks")
	lf := flag.BoolP("list-folders", "d", false, "List all of the folders")
	la := flag.BoolP("list-all", "a", false, "List all folders and the bookmarks (name & URL) within them")
	writeOut := flag.StringP("write-out", "w", "", "Output the results of any list functionality in CSV format")
	checkHealth := flag.BoolP("check-health", "k", false, "List all bookmark links that do not return a 200-OK response")

	// Override the default usage function
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	// Official pflag package has not fixed issues#352 - Need this
	err := flag.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		pflag.Usage()
		os.Exit(1)
	}

	if *help || *usage {
		helpString := "\nNabbit can be used to create Netscape bookmark files or it can\nbe used " +
			"to parse and gather information about them.\n\nSyntax: nabbit [file] " +
			"<options>\n\nExamples:\n\n-Create a bookmark file for Google Chrome from a" +
			".txt file:\n\n\tnabbit test.txt --create-file=chrome -o bookmarks.html\n\n" +
			"-Get the number of bookmarks and folders in a Netscape bookmark file:\n\n\tnabbit " +
			"bookmarks.html --num-bookmarks --num-folders\n\n-List all bookmarks and folders in a " +
			"Netscape bookmark file and write the output to a CSV file:\n\n\tnabbit bookmarks.html --list-all " +
			"--write-out=output.csv\n\n-Check the health of all bookmarks in a Netscape " +
			"bookmark file:\n\n\tnabbit bookmarks.html --check-health\n\n\nThe format of the " +
			"text file that can be used to create bookmark\nfiles can be found at:\nhttps://github" +
			".com/kyletimmermans/nabbit#text-bookmark-file-format\n"
		fmt.Println(helpString)
		pflag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Println("v1.0")
		os.Exit(0)
	}

	// If no args passed (only the program name)
	if len(os.Args) == 1 {
		err := errors.New("Must pass arguments. Use -h to see command usage")
		log.Fatal("Error: ", err)
	}

	// Find input filename in command line args
	inputFile := os.Args[1]
	if !strings.Contains(inputFile, ".txt") && !strings.Contains(inputFile, ".html") {
		err := errors.New("Input file must be a .txt or .html file")
		log.Fatal("Error: ", err)
	}

	// Need file for creation or inspection
	if len(inputFile) == 0 {
		err := errors.New("No file found in arguments for creation or inspection")
		log.Fatal("Error: ", err)
	}

	// Must have .txt for create-file flag
	if strings.Contains(inputFile, ".html") && len(*createFile) > 0 {
		fmt.Println(*createFile)
		err := errors.New("Input file must be .txt for the create-file flag")
		log.Fatal("Error: ", err)
	}

	// Must have .html for inspect flags
	if strings.Contains(inputFile, ".txt") && len(*createFile) == 0 {
		err := errors.New("Input file must be .html for the inspect flag(s)")
		log.Fatal("Error: ", err)
	}

	// Open file
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer file.Close()

	// Prevent redundant listing
	if *la && (*lb || *lf) {
		err := errors.New("Cannot use -la with -lb and/or -lf")
		log.Fatal("Error: ", err)
	}

	// Prevent create file mixing with info flags
	if len(*createFile) > 0 && (*nb || *nf || *ns || *lb || *lf || *la || *checkHealth) {
		err := errors.New("Cannot use the --create-file flag with any other flag")
		log.Fatal("Error: ", err)
	}

	// If create-file without output
	if (len(*createFile) > 0 && len(*output) == 0) || (len(*createFile) == 0 && len(*output) > 0) {
		err := errors.New("Must use --create-file with -o")
		log.Fatal("Error: ", err)
	}

	// If writeOut with no list function
	if len(*writeOut) > 0 && (!*lb && !*lf && !*la && !*checkHealth) {
		err := errors.New("Must use --write-out with a list function")
		log.Fatal("Error: ", err)
	}

	// Make sure writeOut is CSV file
	if len(*writeOut) > 0 && !strings.HasSuffix(*writeOut, ".csv") {
		err := errors.New("--write-out must be a .csv file")
		log.Fatal("Error: ", err)
	}

	// writeOut can only write the results of one list function
	if len(*writeOut) > 0 {
		check := 0

		if *lb {
			check++
		}
		if *lf {
			check++
		}
		if *la {
			check++
		}
		if *checkHealth {
			check++
		}

		if check > 1 {
			err := errors.New("--write-out can only write the results of one list function to a .csv file")
			log.Fatal("Error: ", err)
		}
	}

	// If proper create-file, move forward
	if len(*createFile) > 0 && len(*output) > 0 {

		if *createFile == "chrome" {
			utils.CreateBookmarkFileChrome(file, *output)
		} else if *createFile == "firefox" {
			utils.CreateBookmarkFileFirefox(file, *output)
		} else {
			err := errors.New("Must use \"chrome\" or \"firefox\" with --create-file")
			log.Fatal("Error: ", err)
		}

		return
	}

	// If proper inspect flags, move forward
	if *nb || *nf || *ns || *lb || *lf || *la || *checkHealth {
		var trueFlags []string

		// Function that checks flag value
		checkFlag := func(f *pflag.Flag) {
			if f.Value.String() == "true" && flag.Lookup(f.Name) != nil {
				trueFlags = append(trueFlags, flag.Lookup(f.Name).Name)
			}
		}
		// Run checkFlag on all created flags
		flag.VisitAll(checkFlag)

		if len(*writeOut) > 0 {
			// One to check, one to use
			trueFlags = append(trueFlags, "writeOut")
			trueFlags = append(trueFlags, *writeOut)
		}

		utils.ParseBookmarkFile(file, trueFlags)
		return
	}

	// If neither function runs, no option flags were passed to the program
	finalCheck := errors.New("Must apply option flags to the inputted file")
	log.Fatal("Error: ", finalCheck)

}
