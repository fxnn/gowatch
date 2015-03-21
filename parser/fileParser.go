package parser

import (
	"bufio"
	"github.com/fxnn/gowatch/logentry"
	"log"
	"os"
)

type FileParser struct {
	filename               string
	Err                    error
	logTextToEntryFunction func(line string) logentry.LogEntry
}

func NewFileParser(filename string, logTextToEntryFunction func(line string) logentry.LogEntry) (p *FileParser) {
	p = new(FileParser)

	p.filename = filename
	p.logTextToEntryFunction = logTextToEntryFunction

	return p
}

func (p *FileParser) Parse() <-chan logentry.LogEntry {
	out := make(chan logentry.LogEntry)

	if file := p.openOrFail(p.filename); file != nil {
		go p.doParse(*file, out)
	}

	return out
}

func (p *FileParser) openOrFail(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		p.Err = err
		log.Fatal(err)
		return nil
	}

	return file
}

func (p *FileParser) doParse(file os.File, out chan logentry.LogEntry) {
	defer file.Close()

	scanner := bufio.NewScanner(&file)
	for scanner.Scan() {
		out <- p.logTextToEntryFunction(scanner.Text())
	}
	close(out)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		p.Err = err
	}
}
