package nabbit

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var CurrentUnixTime string = strconv.FormatInt(time.Now().Unix(), 10)

func FileWrite(writer *bufio.Writer, content string) {
	_, err := fmt.Fprint(writer, content)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}
	writer.Flush()
}

func CSVWrite(writer *csv.Writer, items []string) {
	err := writer.Write(items)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
		return
	}
}

func CreateFile(fileName string) *os.File {
	createdFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	return createdFile
}

func VisitURL(url string) string {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 20
	transport.ResponseHeaderTimeout = 20 * time.Second
	transport.IdleConnTimeout = 20 * time.Second

	custom := &http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 20 { // 10 -> 20
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	response, err := custom.Get(url)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no such host"):
			return "No such host"
		case strings.Contains(err.Error(), "i/o timeout"):
			return "Unreachable"
		case strings.Contains(err.Error(), "connection refused"):
			return "Refused"
		case strings.Contains(err.Error(), "stopped after"):
			return "Too many redirects"
		case strings.Contains(err.Error(), "Client.Timeout exceeded"):
			return "Timeout"
		case strings.Contains(err.Error(), "context deadline exceeded"):
			return "Timeout"
		default:
			log.Fatal("Error making HTTP request: ", err)
		}
	}

	defer response.Body.Close()

	statusCode := response.StatusCode
	return strconv.Itoa(statusCode)
}
