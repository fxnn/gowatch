package parser

type SimpleLineSource struct {
	lines []string
}

func NewSimpleLineSource() *SimpleLineSource {
	return &SimpleLineSource{make([]string, 0, 10)}
}

func (ls *SimpleLineSource) AddLine(line string) {
	ls.lines = append(ls.lines, line)
}

func (ls *SimpleLineSource) Lines() <-chan string {
	out := make(chan string, 10)
	go func() {
		for _, line := range ls.lines {
			out <- line
		}
		close(out)
	}()
	return out
}
