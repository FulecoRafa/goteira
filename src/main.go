package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Ports []int64

// Set implements flag.Value.
func (p *Ports) Set(str string) error {
    num, err := strconv.ParseInt(str, 10, 64)
    if err != nil {
        return err
    }
    *p = append(*p, num)
    return nil
}

// String implements flag.Value.
func (p *Ports) String() string {
    var b strings.Builder
    for _, port := range *p {
        fmt.Fprint(&b, port)
        b.WriteString(", ")
    }
    return b.String()

}

var _ flag.Value = &Ports{}

var ports Ports

func init() {
	flag.Var(&ports, "p", "Ports to listen")
}

func main() {
    flag.Parse()

    fmt.Println(ports)
}
