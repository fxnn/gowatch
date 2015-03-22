package parser

import (
	"bufio"
	"log"
	"os"
)

type FileLineSource struct {
	filename string
	Err      error
}

func NewFileLineSource(filename string) (p *FileLineSource) {
	return &FileLineSource{filename, nil}
}

func (p *FileLineSource) Lines() <-chan string {
	out := make(chan string)

	if file := p.openOrFail(p.filename); file != nil {
		go p.doRead(*file, out)
	}

	return out
}

func (p *FileLineSource) openOrFail(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		p.Err = err
		log.Fatal(err)
		return nil
	}

	return file
}

func (p *FileLineSource) doRead(file os.File, out chan string) {
	defer file.Close()

	scanner := bufio.NewScanner(&file)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		p.Err = err
	}
}
