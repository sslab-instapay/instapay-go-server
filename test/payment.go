package main

import (
  "fmt"
  "reflect"
)

func main() {
  var a []string
  var b []string
  a = []string{"AA"}
  b = []string{"BB"}

  fmt.Printf("%t\n", reflect.DeepEqual(a, b))

  var c []string
  var d []string
  c = []string{"11", "22", "33"}

  d = c[0:2]
  fmt.Println(d)
  fmt.Println(c)
}
