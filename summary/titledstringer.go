package summary
import (
    "strings"
    "bytes"
    "fmt"
)

type TitledStringer struct{
    Title        string
    Stringer     fmt.Stringer
}

func (s *TitledStringer) String() string {
    var buffer bytes.Buffer

    buffer.WriteString(s.Title)
    buffer.WriteString("\n")
    buffer.WriteString(strings.Repeat("=", len(s.Title)))
    buffer.WriteString("\n")

    buffer.WriteString(s.Stringer.String())
    buffer.WriteString("\n")

    return buffer.String()
}