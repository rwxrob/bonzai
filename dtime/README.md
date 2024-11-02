# Easy Date/Time Formats with Duration Spans

![WIP](https://img.shields.io/badge/status-wip-red)
[![GoDoc](https://godoc.org/github.com/rwxrob/dtime?status.svg)](https://godoc.org/github.com/rwxrob/dtime)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/dtime)](https://goreportcard.com/report/github.com/rwxrob/dtime)
![License](https://img.shields.io/github/license/rwxrob/dtime)

Returns one or two `*time.Time` pointers, one for the first second for given time, and a bounding second of the end of the duration (second after last second). This allows easy checks for give times within that duration.

```
+20m  - in 20 minutes
+1h   - in an hour
-2.5d - two and half days ago
.h    - this hour 
.d    - today (this day)
t     - tomorrow 
y     - yesterday
ly    - last year
nw    - next week
nm    - next month
lw    - last week
lmon  - last monday
ly+4w - last year for four weeks (24 days)
```

See the [test data](testdata/dtime.yaml) for hundreds of examples and the [PEG grammar](grammar.peg) for specifics.

## Motivation

When using a mobile device the only characters available on the default keyboard are alpha-numeric and the comma (,) and period (.). While it is only a minor convenience to shift to the character keyboard why not create a set of formats that worth with the least amount of trouble. Therefore, these formats use the shortest, best format possible to convey the most common references to dates and times. 

This also makes these time formats particularly useful to add to applications with a terse command-line interface.

## TODO

* Add the `dtime` command (with tab completion) to go with the package.

## See Also

### TJ Holowaychuk's `go-naturaldate` Package

TJ's [go-naturaldate](https://github.com/tj/go-naturaldate) package came out while I was developing this one. I noted his use of PEG and reworked the internals of my package to also use it. 

TJ's package is far better for conversational UIs. 

Mine started with emphasis on the least amount of typing possible and no spaces so that queries can easily be added as singular command-line arguments. 

Mine also comes with the `dtime` command for easy integration into shell scripts or while editing files with `vi/m` using "wand" syntax (`!!`,`!}`,`!G`}. 

I also focus mostly on time spans rather than specific dates.

### Andrew Snodgrass' PEG Golang Package

This [PEG package](https://github.com/pointlander/peg) is truly amazing. My days of writing ABNF are likely over.

## Design Decisions

* **Lowercase for all.** Since the primary motivation is efficiency both in inputting and parsing I decided to keep everything to lowercase which a few exceptions where uppercase has a different meaning.

* **Just one.** No need for more than one way to refer to things. Easier to memorize. 
