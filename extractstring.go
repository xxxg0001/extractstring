package extractstring

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sort"
)

type KeyWord struct {
	SyntaxType string
	Open       string
	Close      string
}

func proccessComment(close string, context []byte, nIndex, nSize int) int {
	i := nIndex
	for ; i < nSize; i++ {
		n := len(close)
		chs := context[i : i+n]
		if bytes.Equal(chs, []byte(close)) {
			return i + n - 1
		}
	}
	return i
}

func proccessString(close string, isEscape bool, context []byte, nIndex, nSize int) (int, int, int) {
	i := nIndex
	for ; i < nSize; i++ {
		if isEscape == true && context[i] == 0x5c {
			i++
			continue
		}
		n := len(close)
		chs := context[i : i+n]
		if bytes.Equal(chs, []byte(close)) {
			return i + n - 1, nIndex, i
		}
	}
	return i, -1, -1
}

type Filter func([]byte) bool

type File struct {
	keywords []KeyWord
}

func (p *File) GetStrings(context []byte, filter Filter) ([][]byte, []int, []int) {
	var entryStart []int
	var entryEnd []int
	var entryBytes [][]byte

	nSize := len(context)
	for i := 0; i < nSize; i++ {

		if context[i] == 0x5c {
			i++
			continue
		}
		keyword, offset := p.getKeyWordOpen(context, i, nSize)
		i += offset
		if keyword == nil {
			continue
		}
		switch keyword.SyntaxType {
		case "Comment":
			i = proccessComment(keyword.Close, context, i, nSize)
		case "String":
			start, end := -1, -1
			i, start, end = proccessString(keyword.Close, true, context, i, nSize)
			if start != -1 && end != -1 {
				result := context[start:end]
				if filter == nil || filter(result) {
					entryStart = append(entryStart, start)
					entryEnd = append(entryEnd, end)
					entryBytes = append(entryBytes, result)
				}
			}
		case "NoEscapeString":
			start, end := -1, -1
			i, start, end = proccessString(keyword.Close, false, context, i, nSize)
			if start != -1 && end != -1 {
				result := context[start:end]
				if filter == nil || filter(result) {
					entryStart = append(entryStart, start)
					entryEnd = append(entryEnd, end)
					entryBytes = append(entryBytes, result)
				}
			}
		}
	}
	return entryBytes, entryStart, entryEnd
}

func (p *File) getKeyWordOpen(context []byte, nIndex, nSize int) (*KeyWord, int) {
	for i := 0; i < len(p.keywords); i++ {
		n := len(p.keywords[i].Open)
		if nIndex+n > nSize {
			continue
		}
		chs := context[nIndex : nIndex+n]
		if bytes.Equal(chs, []byte(p.keywords[i].Open)) {
			return &p.keywords[i], n
		}
	}
	return nil, 0
}

type KeyWords []KeyWord

func (p KeyWords) Len() int           { return len(p) }
func (p KeyWords) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p KeyWords) Less(i, j int) bool { return len(p[i].Open) > len(p[j].Open) }

func (p *File) LoadSyntax(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	p.keywords = make([]KeyWord, 0)
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&p.keywords); err == io.EOF {

	} else if err != nil {
		return err
	}
	sort.Sort(KeyWords(p.keywords))
	return nil
}

func New(syntaxfile string) (*File, error) {
	p := new(File)
	err := p.LoadSyntax(syntaxfile)
	if err != nil {
		return nil, err
	}
	return p, nil
}
