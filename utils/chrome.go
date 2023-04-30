package nabbit

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// Starts every bookmark file
var chromeBookmarkFileIntro string = fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="%s" LAST_MODIFIED="%s" PERSONAL_TOOLBAR_FOLDER="true">Bookmarks Bar</H3>
    <DL><p>`, CurrentUnixTime, CurrentUnixTime)

func CreateFolderChrome(writer *bufio.Writer, name string, tabs int) {
	if name != "" {
		tabString := strings.Repeat("\t", tabs)
		FolderString := fmt.Sprintf("%s<DT><H3 ADD_DATE=\"%s\" LAST_MODIFIED=\"%s\">%s</H3>\n%s<DL><p>\n",
			tabString, CurrentUnixTime, CurrentUnixTime, name, tabString)
		FileWrite(writer, FolderString)
	}
}

func CreateBookmarkChrome(writer *bufio.Writer, name string, link string, tabs int) {
	tabString := strings.Repeat("\t", tabs)
	BookmarkString := fmt.Sprintf("%s<DT><A HREF=\"%s\" ADD_DATE=\"%s\">%s</A>\n", tabString, link, CurrentUnixTime, name)
	FileWrite(writer, BookmarkString)
}

func CreateBookmarkFileChrome(inputFile *os.File, outputLocation string) {
	// Create bookmark file
	bookmarkFile := CreateFile(outputLocation)
	defer bookmarkFile.Close()

	// Create writer object
	writer := bufio.NewWriter(bookmarkFile)

	// Write in bookmark file intro syntax
	FileWrite(writer, chromeBookmarkFileIntro+"\n")

	// For keeping track of when to add end tags for folders
	// [{num_of_tabs, closed?}]
	var closerTracker [][]int
	// Init with value for first run as no lines come before start
	tabCount := 0

	// Read input file
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Set before tabCount is initialized so tabCount will reflect previous line
		// If previous line differs from current line in tabs, we need to close
		if strings.Count(line, "\t") < tabCount {
			for idx, i := range closerTracker {
				// If done adding to folder (less tabs == closed)
				if strings.Count(line, "\t") <= i[0] && i[1] != 1 {
					FileWrite(writer, strings.Repeat("\t", i[0]+2)+"</DL><p>\n")
					closerTracker[idx][1] = 1 // Mark as closed
				}
			}
		}

		tabCount = strings.Count(line, "\t")
		line = strings.ReplaceAll(line, "\t", "")

		// Determine if bookmark or folder
		if strings.Contains(line, "----") {
			separator := strings.Index(line, "----")
			name := line[:separator]
			link := line[separator+4:]

			if len(link) == 0 {
				err := errors.New("Must have a non-empty link for all bookmarks")
				log.Fatal("Error: ", err)
			}
			if !strings.Contains(link, "http") {
				link = "https://" + line[separator+4:]
			}
			if link[len(link)-1:] != "/" {
				link = link + "/"
			}
			// Tabs + 2 bc everything starts under "Bookmarks Bar" base
			CreateBookmarkChrome(writer, name, link, tabCount+2)
		} else {
			closerTracker = append(closerTracker, []int{tabCount, 0})
			CreateFolderChrome(writer, line, tabCount+2)
		}
	}

	// Need biggest to smallest tab closers
	for i := 0; i < len(closerTracker)-1; i++ {
		maxidx := i
		for j := i + 1; j < len(closerTracker); j++ {
			if closerTracker[j][0] > closerTracker[maxidx][0] {
				maxidx = j
			}
		}
		temp := closerTracker[i]
		closerTracker[i] = closerTracker[maxidx]
		closerTracker[maxidx] = temp
	}

	// Close un-closed folders
	for _, x := range closerTracker {
		if x[1] == 0 {
			FileWrite(writer, strings.Repeat("\t", x[0]+2)+"</DL><p>\n")
		}
	}

	// After all lines finished from input, put final end tags and newline
	FileWrite(writer, "\t</DL><p>\n")
	FileWrite(writer, "</DL><p>\n")
}
