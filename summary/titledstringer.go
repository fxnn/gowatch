package summary
import (
    "strings"
    "bytes"
    "fmt"
)

type TitledStringer struct{
    title		string
    stringer	fmt.Stringer
}

func (s *TitledStringer) String() string {
    var buffer bytes.Buffer

    buffer.WriteString(s.title)
    buffer.WriteString("\n")
    buffer.WriteString(strings.Repeat("=", len(s.title)))
    buffer.WriteString("\n")

    buffer.WriteString(s.stringer.String())
    buffer.WriteString("\n\n")

    return buffer.String()
}