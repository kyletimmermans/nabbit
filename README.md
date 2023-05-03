![Version 1.0](https://img.shields.io/badge/version-v1.0-orange.svg)
![Go 1.20](https://img.shields.io/badge/Go-1.20-00acd7.svg)
![Last Updated](https://img.shields.io/github/last-commit/kyletimmermans/nabbit?color=success)
[![kyletimmermans Twitter](http://img.shields.io/twitter/url/http/shields.io.svg?style=social&label=Follow)](https://twitter.com/kyletimmermans)

# <div align="center">Nabbit</div>

<p align="center">A command-line tool for parsing and creating Netscape bookmark files (Chrome & Firefox)</p>

</br>

<p align="center">
  <img src="https://github.com/kyletimmermans/nabbit/raw/main/media/nabbit_512x512.png?raw=true" width="50%" height="50%" alt="Nabbit"/>
</p>


</br>


Table of Contents
=================

<!--ts-->
   * [Installation](#installation)
   * [Usage](#usage)
      * [Flags](#flags)
      * [Examples](#examples)
      * [Text Bookmark File Format](#text-bookmark-file-format)
   * [Changelog](#changelog)
<!--te-->

</br>

### Installation

###### Linux
```bash
curl -q -s -LJO "https://github.com/kyletimmermans/nabbit/releases/download/latest/nabbit-v1-linux-amd64.tar.gz" && tar -xzf nabbit-v1-linux-amd64.tar.gz && rm nabbit-v1-linux-amd64.tar.gz && chmod +x nabbit
```

###### Mac (Intel)
```bash
curl -q -s -LJO "https://github.com/kyletimmermans/nabbit/releases/download/latest/nabbit-v1-mac-amd64.tar.gz" && tar -xzf nabbit-v1-mac-amd64.tar.gz && rm nabbit-v1-mac-amd64.tar.gz && chmod +x nabbit
```

###### Mac (Apple Silicon - M1)
```bash
curl -q -s -LJO "https://github.com/kyletimmermans/nabbit/releases/download/latest/nabbit-v1-mac-arm64.tar.gz" && tar -xzf nabbit-v1-mac-arm64.tar.gz && rm nabbit-v1-mac-arm64.tar.gz && chmod +x nabbit
```

###### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri "https://github.com/kyletimmermans/nabbit/releases/download/latest/nabbit-v1-windows-amd64.zip" -OutFile "nabbit-v1-windows-amd64.zip"; Expand-Archive -LiteralPath "nabbit-v1-windows-amd64.zip" -DestinationPath "nabbit"; icacls ".\nabbit\nabbit.exe" /grant *S-1-1-0:(X)
```

</br>

### Usage

```
$ nabbit <file> [options]
```

### Flags
```
  -k, --check-health         List all bookmark links that do not return a 200-OK response
  -c, --create-file string   Choose bookmark file type for creation from .txt file (chrome/firefox)
  -h, --help                 Display help/usage information
  -a, --list-all             List all folders and the bookmarks (name & URL) within them 
  -l, --list-bookmarks       List all of the bookmarks
  -d, --list-folders         List all of the folders
  -b, --num-bookmarks        Get total number of bookmarks
  -f, --num-folders          Get total number of folders
  -s, --num-separators       Get total number of separators (Firefox only)
  -o, --output string        Output filename of bookmark file (used w/ create-file)
  -u, --usage                Display help/usage information
  -v, --version              Display version
  -w, --write-out string     Output the results of any list functionality in CSV format
```

### Examples

#### <ins>Create a Bookmark File from a Text File (Chrome)</ins>

###### Command
```
$ nabbit chrome_example.txt --browser-type=chrome -o bookmarks.html
```

###### chrome_example.txt (http/https:// or trailing forward slash not required)
```
Google----https://www.google.com/
YouTube.com----https://www.youtube.com/
Cool Stuff
	OpenAI----https://www.openai.com
	Machine Learning
		IBM fun stuff----https://www.ibm.com
	Wikipedia - Cool----https://wikipedia.com
Fun Stuff
	Amazon----https://www.amazon.com
CERN----http://info.cern.ch/
```

###### Resulting bookmark file will look like this in browser
<p align="left">
  <img src="https://github.com/kyletimmermans/nabbit/raw/main/media/bookmarks_chrome.png?raw=true" alt="Bookmarks"/>
</p>

</br>

#### <ins>Create a Bookmark File from a Text File (Firefox)</ins>

###### Command
```
$ nabbit firefox_example.txt --browser-type=firefox -o bookmarks.html
```

###### firefox_example.txt ('++++' represents a separator line)
```
Example.com----www.example.com
++++
Mozilla Firefox
	Google----https://www.google.com/
	YouTube----https://www.youtube.com/
	Cool Stuff
		OpenAI----https://www.openai.com/
	++++
	++++
	Mozilla----https://www.mozilla.org/
My Stuff
	GitHub----https://www.github.com/
	Kyle----https://www.kyles.world/
Microsoft----https://www.microsoft.com
```

###### Resulting bookmark file will look like this in browser
<p align="left">
  <img src="https://github.com/kyletimmermans/nabbit/raw/main/media/bookmarks_firefox.png?raw=true" alt="Bookmarks"/>
</p>

</br>

#### <ins>Get Number of Items in Bookmark File</ins>

###### Command
```
$ nabbit firefox_bookmarks.html --num-bookmarks --num-folders --num-separators
```

###### Output
```
Number of Bookmarks: 7

Number of Folders: 3

Number of Separators: 3
```

</br>

#### <ins>List Information from Bookmark File</ins>

###### Command
```
$ nabbit chrome_bookmarks.html --list-all --write-out=test.csv
```

###### Output (Note that long names/URLs will be spliced for terminal output, use --write-out to get the full lines in a .csv file)
```
Bookmark & Folder List:

┏━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━┓
┃ Folder Name      ┃ Tab Name         ┃ URL              ┃
┡━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━┩
│ Bookmarks Bar    │ Google           │ www.google.com/  │
│ Bookmarks Bar    │ YouTube.com      │ www.youtube.com/ │
│ Bookmarks Bar    │ CERN             │ info.cern.ch/    │
│ Cool Stuff       │ OpenAI           │ www.openai.com/  │
│ Cool Stuff       │ Wikipedia - Cool │ wikipedia.com/   │
│ Machine Learning │ IBM fun stuff    │ www.ibm.com/     │
│ Fun Stuff        │ Amazon           │ www.amazon.com/  │
└──────────────────┴──────────────────┴──────────────────┘
```

</br>


### Text Bookmark File Format
#### For creating Netscape bookmark files from text files

###### Bookmarks - The name and the actual URL must be separated by a '----' (quadruple hyphen). 
```
Google----https://www.google.com/
YouTube.com----https://www.youtube.com/
```

###### Folders - The bookmarks under a folder must be tabbed over to represent the folder hierarchy
```
My Folder
	Google----https://www.google.com/
	YouTube.com----https://www.youtube.com/
	Nested Folder
		Another Link----www.example.com
```

###### Separators - Denoted in a bookmark file by a '++++' (quadruple plus) (Firefox only)
```
Google----https://www.google.com/
++++
YouTube.com----https://www.youtube.com/
```

Notes on Text Bookmark File Format:
- [x] Bookmarks - https:// and a trailing forward slash are not necessary as they are automatically added if not detected in the input file. 
- [x] Bookmarks - https:// will not replace http:// if you added it manually.
- [x] General - Do not put bookmarks, folders, or separators on the same line of your file, they each get their own line
- [x] General - Empty lines will be removed
- [x] General - Chrome and Firefox will automatically retrieve and render bookmark favicons when you import the HTML bookmark file

</br>

### Changelog
| Version  | Notes |
| :---: | :---: |
| 1.0 | Initial-Relase |
