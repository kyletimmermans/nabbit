package nabbit

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/markkurossi/tabulate"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/net/html"
)

func bookmarkNumber(inputFile *os.File) int {

	var count int

	inputFile.Seek(0, io.SeekStart)

	tokenizer := html.NewTokenizer(inputFile)

	for {

		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		// Only get the html tag once
		if tokenType == html.StartTagToken {
			if token.Data == "a" {
				count++
			}
		}
	}

	return count
}

func listBookmarks(tokenizer *html.Tokenizer, token html.Token, tab *tabulate.Tabulate,
	writer *csv.Writer, writeOut string) {
	var url string

	for _, attr := range token.Attr {
		if attr.Key == "href" {
			url = attr.Val
		}
	}

	// Get Bookmark Name
	tokenType := tokenizer.Next()
	if tokenType == html.TextToken {
		row := tab.Row()
		name := tokenizer.Token().Data
		if len(name) > 50 {
			row.Column(name[:50] + "-")
		} else {
			row.Column(name)
		}

		// For the terminal only
		previewURL := strings.Replace(strings.Replace(url, "https://", "", 1), "http://", "", 1)

		if len(previewURL) > 50 {
			// No https in preview for more terminal space
			previewURL = previewURL[:50] + "-"
		}

		row.Column(previewURL)

		if len(writeOut) > 0 {
			// Full URL w/ --write-out
			CSVWrite(writer, []string{name, url})
		}
	}
}

func listFolders(inputFile *os.File, tokenizer *html.Tokenizer, token html.Token,
	tab *tabulate.Tabulate, writer *csv.Writer, writeOut string) {
	var count int = 0
	var name string

	tokenType := tokenizer.Next()
	if tokenType == html.TextToken {
		row := tab.Row()
		name = tokenizer.Token().Data
		if len(name) > 50 {
			row.Column(name[:50] + "-")
		} else {
			row.Column(name)
		}

		// Open up input file again, go to the line we just captured
		// and go up until the end of the folder, adding all the items to count
		copiedFile, _ := filepath.Abs(inputFile.Name())
		copyFileOpen, err := os.Open(copiedFile)
		if err != nil {
			log.Fatal("Error opening file: ", err)
			return
		}
		defer copyFileOpen.Close()

		// Track when folder closes
		trackerIndex := 0
		startPoint := false

		tokenizerTwo := html.NewTokenizer(copyFileOpen)

		for {
			tokenTypeTwo := tokenizerTwo.Next()
			if tokenTypeTwo == html.ErrorToken {
				break
			}

			tokenTwo := tokenizerTwo.Token()

			if tokenTypeTwo == html.TextToken && tokenTwo.Data == name {
				startPoint = true
			}

			if startPoint {
				if tokenTypeTwo == html.StartTagToken && tokenTwo.Data == "dl" {
					trackerIndex++
				} else if tokenTypeTwo == html.EndTagToken && tokenTwo.Data == "dl" {
					// Moving from right to left down tracker, figure out if
					// the original open tag for the folder is closed
					trackerIndex--
					if trackerIndex == 0 {
						row.Column(strconv.Itoa(count))
						startPoint = false
						break
					}
				} else if tokenTypeTwo == html.StartTagToken && (tokenTwo.Data == "a" || tokenTwo.Data == "h3") {
					count++
				}
			}
		}
	}

	if len(writeOut) > 0 {
		CSVWrite(writer, []string{name, strconv.Itoa(count)})
	}
}

func listAll(inputFile *os.File, tokenizer *html.Tokenizer, token html.Token,
	tab *tabulate.Tabulate, writer *csv.Writer, bar *mpb.Bar, writeOut string, checkHealth bool) {

	tokenType := tokenizer.Next()
	if tokenType == html.TextToken {
		folderName := tokenizer.Token().Data

		// Open up input file again, go to the line we just captured
		// and go up until the end of the folder, adding all the items to count
		copiedFile, _ := filepath.Abs(inputFile.Name())
		copyFileOpen, err := os.Open(copiedFile)
		if err != nil {
			log.Fatal("Error opening file: ", err)
			return
		}
		defer copyFileOpen.Close()

		// Track when folder closes
		trackerIndex := 0
		startPoint := false

		tokenizerTwo := html.NewTokenizer(copyFileOpen)

		for {
			tokenTypeTwo := tokenizerTwo.Next()
			if tokenTypeTwo == html.ErrorToken {
				break
			}

			tokenTwo := tokenizerTwo.Token()

			if tokenTypeTwo == html.TextToken && tokenTwo.Data == folderName {
				startPoint = true
			}

			if startPoint {
				if tokenTypeTwo == html.StartTagToken && tokenTwo.Data == "dl" {
					trackerIndex++
				} else if tokenTypeTwo == html.EndTagToken && tokenTwo.Data == "dl" {
					// Moving from right to left down tracker, figure out if
					// the original open tag for the folder is closed
					trackerIndex--
					if trackerIndex == 0 {
						startPoint = false
						break
					}
				} else if tokenTypeTwo == html.StartTagToken && tokenTwo.Data == "a" {
					// 1 is our tab level for this folder, add only these
					if trackerIndex == 1 {
						var url string
						for _, attr := range tokenTwo.Attr {
							if attr.Key == "href" {
								url = attr.Val
							}
						}

						var responseCode string
						if checkHealth {
							responseCode = VisitURL(url)
							if responseCode == "200" {
								bar.Increment()
								continue
							}
						}

						previewURL := strings.Replace(strings.Replace(url,
							"https://", "", 1), "http://", "", 1)

						if checkHealth {
							if len(previewURL) > 40 {
								previewURL = previewURL[:40] + "-"
							}
						} else {
							if len(previewURL) > 45 {
								previewURL = previewURL[:45] + "-"
							}
						}

						var tabName string
						tokenTypeTwo = tokenizerTwo.Next()
						if tokenTypeTwo == html.TextToken {
							tabName = tokenizerTwo.Token().Data
						}

						row := tab.Row()

						if len(folderName) > 45 {
							row.Column(folderName[:45] + "-")
						} else {
							row.Column(folderName)
						}

						if len(tabName) > 50 {
							row.Column(tabName[:50] + "-")
						} else {
							row.Column(tabName)
						}

						row.Column(previewURL)

						if checkHealth {
							row.Column(responseCode)
							// Once non-200 added, then we can increment
							bar.Increment()
						}

						if len(writeOut) > 0 {
							if !checkHealth {
								CSVWrite(writer, []string{folderName, tabName, url})
							} else {
								CSVWrite(writer, []string{folderName, tabName, url, responseCode})
							}
						}
					}
				}
			}
		}
	}
}

func listItems(inputFile *os.File, itemType string, writeOut string) {
	var outputFile *os.File
	var writer *csv.Writer
	var p *mpb.Progress
	var bar *mpb.Bar

	// Progress Bar
	if itemType == "health" {
		p = mpb.New(mpb.WithWidth(64))
		bar = p.AddBar(int64(bookmarkNumber(inputFile)),
			mpb.PrependDecorators(
				decor.Name("Progress: "),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncSpace), "Done!",
				),
			),
		)
	}

	if len(writeOut) > 0 {
		outputFile = CreateFile(writeOut)
		writer = csv.NewWriter(outputFile)
		defer writer.Flush()
	}

	tab := tabulate.New(tabulate.Unicode)

	switch itemType {
	case "bookmark":
		tab.Header("Name").SetAlign(tabulate.ML)
		tab.Header("URL").SetAlign(tabulate.ML)
	case "folder":
		tab.Header("Name").SetAlign(tabulate.ML)
		tab.Header("No. of Items").SetAlign(tabulate.ML)
	case "all":
		tab.Header("Folder Name").SetAlign(tabulate.ML)
		tab.Header("Tab Name").SetAlign(tabulate.ML)
		tab.Header("URL").SetAlign(tabulate.ML)
	case "health":
		tab.Header("Folder Name").SetAlign(tabulate.ML)
		tab.Header("Tab Name").SetAlign(tabulate.ML)
		tab.Header("URL").SetAlign(tabulate.ML)
		tab.Header("Error").SetAlign(tabulate.ML)
	}

	// Go back to the start of the file as if it was read before
	// the file pointer is at the end and doesn't reset automatically
	inputFile.Seek(0, io.SeekStart)

	// Create tokenized html file (record keeper of html parts)
	tokenizer := html.NewTokenizer(inputFile)

	// Iterate through tokenized html file
	// tokenizer.Next() to keep moving forward each step
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		// Only get the html tag once
		if tokenType == html.StartTagToken {
			// anchor tag is a bookmark
			switch itemType {
			case "bookmark":
				if token.Data == "a" {
					listBookmarks(tokenizer, token, tab, writer, writeOut)
				}
			case "folder":
				if token.Data == "h3" {
					listFolders(inputFile, tokenizer, token, tab, writer, writeOut)
				}
			case "all":
				if token.Data == "h3" {
					listAll(inputFile, tokenizer, token, tab, writer, bar, writeOut, false)
				}
			case "health":
				if token.Data == "h3" {
					listAll(inputFile, tokenizer, token, tab, writer, bar, writeOut, true)
				}
			}
		}
	}

	// Wait for bar to complete and flush
	if itemType == "health" {
		p.Wait()
	}

	switch itemType {
	case "bookmark":
		fmt.Println("\nBookmark List:\n")
		tab.Print(os.Stdout)
	case "folder":
		fmt.Println("\nFolder List:\n")
		tab.Print(os.Stdout)
	case "all":
		fmt.Println("\nBookmark & Folder List:\n")
		tab.Print(os.Stdout)
	case "health":
		fmt.Println("\nBroken Bookmark Links:\n")
		tab.Print(os.Stdout)
	}
}

func getNumber(inputFile *os.File, itemType string) {

	var count int

	// Go back to the start of the file as if it was read before
	// the file pointer is at the end and doesn't reset automatically
	inputFile.Seek(0, io.SeekStart)

	// Create tokenized html file (record keeper of html parts)
	tokenizer := html.NewTokenizer(inputFile)

	// Iterate through tokenized html file
	// tokenizer.Next() to keep moving forward each step
	for {

		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		// Only get the html tag once
		if tokenType == html.StartTagToken {
			// anchor tag is a bookmark
			switch itemType {
			case "bookmark":
				if token.Data == "a" {
					count++
				}
			case "folder":
				if token.Data == "h3" {
					count++
				}
			case "separator":
				if token.Data == "hr" {
					count++
				}
			}
		}
	}

	switch itemType {
	case "bookmark":
		fmt.Printf("\nNumber of Bookmarks: %d\n", count)
	case "folder":
		fmt.Printf("\nNumber of Folders: %d\n", count)
	case "separator":
		fmt.Printf("\nNumber of Separators: %d\n", count)
	}
}

func ParseBookmarkFile(inputFile *os.File, args []string) {
	var writeOut string

	if len(args) > 2 {
		if args[len(args)-2] == "writeOut" {
			writeOut = args[len(args)-1]
		} else {
			writeOut = ""
		}
	} else {
		writeOut = ""
	}

	// Separate and do numbers first so they show up in stdout first
	for _, i := range args {
		switch i {
		case "b", "num-bookmarks":
			getNumber(inputFile, "bookmark")
		case "f", "num-folders":
			getNumber(inputFile, "folder")
		case "s", "num-separators":
			getNumber(inputFile, "separator")
		}
	}

	for _, i := range args {
		switch i {
		case "l", "list-bookmarks":
			listItems(inputFile, "bookmark", writeOut)
		case "d", "list-folders":
			listItems(inputFile, "folder", writeOut)
		case "a", "list-all":
			listItems(inputFile, "all", writeOut)
		case "k", "check-health":
			listItems(inputFile, "health", writeOut)
		}
	}
}
