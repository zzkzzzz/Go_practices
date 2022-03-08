package main

import (
    "fmt"
    "unicode/utf8"
)

// reference: https://go.dev/blog/strings
// Go source code is always UTF-8.
// A string holds arbitrary bytes.
// A string literal, absent byte-level escapes, always holds valid UTF-8 sequences.
// Those sequences represent Unicode code points, called runes.
// No guarantee is made in Go that characters in strings are normalized
func main() {

	// s is a string assigned a literal value representing the word “hello” in the Thai language. 
	// Go string literals are UTF-8 encoded text.
    const s = "สวัสดี"


	
	// Since strings are equivalent to []byte, this will produce the length of the raw bytes stored within.
    fmt.Println("Len:", len(s))


	fmt.Printf("plain string: ")
    fmt.Printf("%s", s)
    fmt.Printf("\n")

	// The %q (quoted) verb will escape any non-printable byte sequences in a string so the output is unambiguous.
    fmt.Printf("quoted string: ")
    fmt.Printf("%q", s)
    fmt.Printf("\n")

	// If we are unfamiliar or confused by strange values in the string, we can use the “plus” flag to the %q verb. 
	// This flag causes the output to escape not only non-printable sequences, but also any non-ASCII bytes, all while interpreting UTF-8. 
	fmt.Printf("quoted string: ")
    fmt.Printf("%+q", s)
    fmt.Printf("\n")

    fmt.Printf("hex bytes: ")
	// generates the hex values of all the bytes that constitute the code points in s.
    for i := 0; i < len(s); i++ {
        fmt.Printf("%x ", s[i])
    }
    fmt.Println()

	// To count how many runes are in a string, we can use the utf8 package.
    fmt.Println("Rune count:", utf8.RuneCountInString(s))

	// A range loop handles strings specially and decodes each rune along with its offset in the string.
    for idx, runeValue := range s {
        fmt.Printf("%#U starts at %d\n", runeValue, idx)
    }

	// achieve the same iteration by using the utf8.DecodeRuneInString function explicitly.
    fmt.Println("\nUsing DecodeRuneInString")
    for i, w := 0, 0; i < len(s); i += w {
        runeValue, width := utf8.DecodeRuneInString(s[i:])
        fmt.Printf("%#U starts at %d\n", runeValue, i)
        w = width

		// passing a rune value to a function.
        examineRune(runeValue)
    }
}

func examineRune(r rune) {
	// Values enclosed in single quotes are rune literals. 
    if r == 't' {
        fmt.Println("found tee")
    } else if r == 'ส' {
        fmt.Println("found so sua")
    }
}