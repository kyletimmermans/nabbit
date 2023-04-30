package nabbit

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var firefoxBookmarkFileIntro string = `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<meta http-equiv="Content-Security-Policy"
      content="default-src 'self'; script-src 'none'; img-src data: *; object-src 'none'"></meta>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks Menu</H1>

<DL><p>`

func CreateFolderFirefox(writer *bufio.Writer, name string, open int, close int) {
	if name != "" {
		var FolderString string
		openTabs := strings.Repeat("    ", open)
		endTabs := strings.Repeat("    ", close)

		FolderString = fmt.Sprintf("%s<DT><H3 ADD_DATE=\"%s\" LAST_MODIFIED=\"%s\">%s</H3>\n%s<DL><p>\n",
			openTabs, CurrentUnixTime, CurrentUnixTime, name, endTabs)

		FileWrite(writer, FolderString)
	}
}

func CreateBookmarkFirefox(writer *bufio.Writer, name string, link string, tabs int) {
	spaceString := strings.Repeat("    ", tabs)
	BookmarkString := fmt.Sprintf("%s<DT><A HREF=\"%s\" ADD_DATE=\"%s\" LAST_MODIFIED=\"%s\">%s</A>\n",
		spaceString, link, CurrentUnixTime, CurrentUnixTime, name)
	FileWrite(writer, BookmarkString)
}

func CreateBookmarkFileFirefox(inputFile *os.File, outputLocation string) {
	// Create bookmark file
	bookmarkFile := CreateFile(outputLocation)
	defer bookmarkFile.Close()

	// Create writer object
	writer := bufio.NewWriter(bookmarkFile)

	// Write in bookmark file intro syntax
	FileWrite(writer, firefoxBookmarkFileIntro+"\n")

	// For keeping track of when to add end tags for folders
	// [{num_of_tabs, closed?}]
	var closerTracker [][]int
	// Init with value for first run as no lines come before start
	tabCount := 0
	separatorTracker := false
	preventExtraTab := false

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
					FileWrite(writer, strings.Repeat("    ", i[0]+1)+"</DL><p>\n")
					closerTracker[idx][1] = 1 // Mark as closed
				}
			}
		}

		tabCount = strings.Count(line, "\t")
		line = strings.ReplaceAll(line, "\t", "")

		// If you have a separator right after another separator, no tabs
		if strings.Contains(line, "++++") {
			if separatorTracker == true {
				FileWrite(writer, "<HR>"+strings.Repeat(" ", 8))
			} else {
				FileWrite(writer, strings.Repeat("    ", tabCount+1)+"<HR>"+strings.Repeat(" ", 4))
			}

			separatorTracker = true
		} else {
			separatorTracker = false
		}

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

			if preventExtraTab == false {
				CreateBookmarkFirefox(writer, name, link, tabCount+1)
			} else {
				CreateBookmarkFirefox(writer, name, link, 0)
			}
		} else if !strings.Contains(line, "++++") {
			// Otherwise its a bookmark
			closerTracker = append(closerTracker, []int{tabCount, 0})
			if preventExtraTab == false {
				CreateFolderFirefox(writer, line, tabCount+1, tabCount+1)
			} else {
				CreateFolderFirefox(writer, line, 0, tabCount+1)
			}
		}

		if strings.Contains(line, "++++") {
			preventExtraTab = true
		} else {
			preventExtraTab = false
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
			FileWrite(writer, strings.Repeat("    ", x[0]+1)+"</DL><p>\n")
		}
	}

	// After all lines finished from input, put final end tags and newline
	FileWrite(writer, "</DL>\n")
}
