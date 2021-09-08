package tests

import (
    "testing"
    "fmt"
)

func TestFib(t *testing.T) {
    var (
        in int      = 7
        expected int  = 13
    )
    actual := Fib(in)
    if actual != expected {
        fmt.Printf("Fib(%d) = %d; expected %d", in, actual, expected)
    }
}
