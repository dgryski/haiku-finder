package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func appendIfUnique(w string, l []int, n int) []int {

	// ugly, but len(l) == 0 almost always, and 1 or 2 very rarely
	for _, v := range l {
		if v == n {
			return l
		}
	}

	return append(l, n)
}

func loadCmudict(path string) (map[string][]int, error) {

	m := make(map[string][]int, 140000)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(f)

	for scan.Scan() {
		s := scan.Text()
		if s[0] == ';' {
			// skip comments
			continue
		}

		// find first word
		idx := strings.Index(s, " ")
		w := s[0:idx]

		if w[idx-1] == ')' {
			w = w[:idx-3]
		}

		c := 0
		// count syllables == digits in remaining string
		for _, r := range s[idx:] {
			if r >= '0' && r <= '9' {
				c++
			}
		}

		syl := m[w]
		syl = appendIfUnique(w, syl, c)
		m[w] = syl
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func formatPoem(poem []string, syllables []int) (string, error) {

	l := 0
	c := syllables[l]

	h := &bytes.Buffer{}

	processed := 0

	for _, p := range poem {
		processed++
		word := p

		word = strings.ToUpper(word)

		word = strings.Map(func(r rune) rune {
			if 'A' <= r && r <= 'Z' {
				return r
			}
			return -1
		}, word)

		a := cmudict[word]
		if len(a) == 0 {

			prefixes := []struct {
				prefix    string
				syllables int
			}{
				{"ANTI", 2},
				{"DE", 1},
				{"DIS", 1},
				{"EX", 1},
				{"MEGA", 2},
				{"MINI", 2},
				{"MIS", 1},
				{"MULTI", 2},
				{"NON", 1},
				{"POST", 1},
				{"PRE", 1},
				{"PRO", 1},
				{"PROTO", 2},
				{"QUASI", 2},
				{"RE", 1},
				{"SEMI", 2},
				{"UN", 1},
				{"VICE", 1},
			}

			for _, p := range prefixes {
				if strings.HasPrefix(word, p.prefix) {
					w := strings.TrimPrefix(word, p.prefix)
					a = cmudict[w]
					if len(a) > 0 {
						a[0] += p.syllables
						break
					}
				}
			}

			if len(a) == 0 {
				return "", errors.New("unknown word: " + word)
			}
		}

		if len(a) > 1 {
			return "", errors.New("don't yet handle words with multiple syllable counts")
		}

		c -= a[0]

		if c < 0 {
			break
		}

		fmt.Fprint(h, p)

		if c > 0 {
			fmt.Fprint(h, " ")
		} else {
			fmt.Fprint(h, "\n")
			l++
			if l >= len(syllables) {
				break
			}
			c = syllables[l]
		}
	}

	if processed != len(poem) || c != 0 {
		return "", errors.New("not a haiku")
	}

	return h.String(), nil
}

var cmudict map[string][]int

func main() {
	var err error
	cmudict, err = loadCmudict("cmudict.0.7a")
	if err != nil {
		log.Println("failed to load cmudict0.7.a:", err)
		log.Println("Please download a copy with:")
		log.Println("curl -o cmudict.0.7a 'http://sourceforge.net/p/cmusphinx/code/11879/tree/trunk/cmudict/cmudict.0.7a?format=raw'")
		return
	}

	scan := bufio.NewScanner(os.Stdin)
	scan.Split(bufio.ScanWords)

	var poem []string

	var sentenceFinished bool

	for scan.Scan() {
		t := scan.Text()

		poem = append(poem, t)
		lastRune := t[len(t)-1]

		if lastRune == '.' || lastRune == '?' || lastRune == '!' {
			sentenceFinished = true
		}

		if sentenceFinished {

			switch t {
			case "Mr.":
				fallthrough
			case "Dr.":
				fallthrough
			case "Ms.":
				fallthrough
			case "Mrs.":
				fallthrough
			case "Sr.":
				sentenceFinished = false
			}

			if sentenceFinished {

				syllables := []int{5, 7, 5}

				p, err := formatPoem(poem, syllables)
				if err == nil {
					fmt.Println(p)
				}
				sentenceFinished = false
				poem = poem[:0]
			}
		}
	}
}
