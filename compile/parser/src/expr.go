package source

type Counter int

var c Counter

func count() {
	c += 1
	r := c
}
