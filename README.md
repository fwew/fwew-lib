# fwew Library
[![Build Status](https://travis-ci.com/tirea/fwew.svg?branch=master)](https://travis-ci.com/tirea/fwew) 
[![License: GPL v2](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html)

The Best Na'vi Dictionary library

## Development
This option is mostly for Contributors and Developers. Or people who like to compile stuff themselves.
You will need the [GO Programming Language](https://golang.org/) and [Git](https://git-scm.com/) installed. 

### Setup
We are using go modules so no GOPATH setup is needed.

To compile and run tests:
```shell script
cd ~/wherever/you/want
git clone https://github.com/fwew/fwew-lib
cd fwew-lib
go test ./...
```

Now make changes to the code and have fun.
Please also add tests for the new code, so we get a high code coverage.

## Usage
We have already two programs, that are using this library:
- A discord bot: https://github.com/fwew/discord-bot
- A CLI program: https://github.com/fwew/fwew

### Translate
Translating is possible from Na'vi or any other supported language.
As result, you will get an array of the Word struct, that can be used to create the output you desire.
```go
require (
    fwew "github.com/fwew/fwew-lib"
)

// Translate from a native language
navi = fwew.TranslateToNavi("search", "en")
fmt.Println(word.ToOutputLine(0, true, false, false, false, false, false))

// Translate a Na'vi word into the native language
navi = fwew.TranslateFromNavi("mllte", "en")
fmt.Println(word.ToOutputLine(0, true, false, false, false, false, false))
```

### Numbers
Numbers also can be translated in both directions.
The Na'vi number system is base 8 (Oktal) and therefore the Integers are base 8.
Number have to be in the range, that Na'vi is possible of saying 0o0 to 0o77777.
```go
require (
    fwew "github.com/fwew/fwew-lib"
)

// Translate an octal number to the Na'vi word
navi, err := fwew.NumberToNavi(0o56)
if err != nil {
    panic(err)
}
fmt.Println(navi)

// Translate the Na'vi word into the octal integer
number, err := fwew.NaviToNumber("mrrvomrr")
if err != nil {
    panic(err)
}
fmt.Println(number)
```

### List

`List()` is a powerful search feature of `fwew` that allows you to list all the words that satisfy a set of given conditions.
Every word has to be in the string array given to `List()`. Simply explode the given string at the space.

The syntax is as follows (cond is short for condition, spec is short for specification):
```
what cond spec [and what cond spec...]
```

`what` can be any one of the following:
```
pos          part of speech of na'vi word
word         na'vi word
syllables    number of syllables in the na'vi word
words        selection of na'vi words
stress       number of stressed syllables in the na'vi word
```

`cond` depends on the `what`. Here are the conditions that apply to each `what`:
pos:
```
has    part of speech has the following character sequence anywhere
is     part of speech is exactly the following character sequence
like   part of speech is like (matches) the following wildcard pattern
```
word:
```
starts    word starts with the following character sequence
ends      word ends with the following character sequence
has       word has the following character sequence anywhere
like      word is like (matches) the following wildcard pattern
```
syllables and stress:
```
<     less than the following number
<=    less than or equal to the following number
=     exactly equal to the following number
>=    greater than or equal to the following number
>     greater than the following number
```
words:
```
first    the first consecutive words in the datafile (chronologically oldest words) 
last     the last consecutive words in the datafile (chronologically newest words)
```

`spec` depends on the `cond`. Here are the specifications that apply to each `cond`:

`has`, `is`, `starts`, and `ends` all expect a character sequence to come next.

`<`, `<=`, `=`, `>=`, `>`, `first`, and `last` all expect a number to come next.

`like` expects a character sequence, usually containing at least one wildcard asterisk, to come next. 

#### Examples of List()

List all modal verbs:
```go
fwew.List([]string{"pos", "has", "v", "and", "pos", "has", "m.",})
```
List all stative verbs:
```go
fwew.List([]string{"pos", "has", "svin.",})
```
List all nouns that start with tì:
```go
fwew.List([]string{"word", "starts", "tì", "and", "pos", "is", "n.",})
```
List all 3 syllable transitive verbs:
```go
fwew.List([]string{"syllables", "=", "3", "and", "pos", "has", "vtr.",})
```
List the newest 25 words in the language:
```go
fwew.List([]string{"words", "last", "25",})
```

### Random

`Random()` is a random entry generator that generates the given number (or random number!) of random entries. 
It also features optional clause in which the `what cond spec` syntax from `List()` is supported to narrow down what kinds of random entries you get.

#### Examples of Random

List 10 random entries
```go
fwew.Random("en", 10, nil)
```
List 5 random transitive verbs
```go
fwew.Random("en", 5, []string{"pos", "has", "vtr",})
```
List a random number of random words
```go
fwew.Random("en", 0, nil)
```
List a random number of nouns
```go
fwew.Random("en", 0, []string{"pos", "is", "n.",})
```

### Update dictionary
`fwew.Update()` will update the dictionary file to the newest version, downloaded from https://tirea.learnnavi.org/dictionarydata/dictionary.txt.  
It will NOT update this library. To update the library you need to adjust the `go mod` of your project.
