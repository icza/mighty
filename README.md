# mighty

[![GoDoc](https://godoc.org/github.com/icza/mighty?status.svg)](https://godoc.org/github.com/icza/mighty) [![Go Report Card](https://goreportcard.com/badge/github.com/icza/mighty)](https://goreportcard.com/report/github.com/icza/mighty)

Package `mighty` is a lightweight extension to Go's [testing](https://golang.org/pkg/testing/) package.

With this utility, you can make your Go test files super clean but still intuitive.
This package doesn't want to hide complex things from you, but wants to eliminate repetitive,
long code from your tests.

## Using `mighty`

You could create a value of `mighty.Myt` and use its methods like `m.Eq()`, `m.Neq()`, ... etc. as seen below:

	m := mighty.Myt{t} // t is of type *testing.T
	m.Eq(6, len("mighty")) // Expect len("mighty") to be 6
	r := bytes.NewBuffer([]byte{'a'})
	m.ExpEq(byte('a'))(r.ReadByte()) // Expect to read 'a' and no error

But the recommended, more intuitive and more compact way is to acquire and use [method values](https://golang.org/ref/spec#Method_values)
returned by functions of `mighty`:

	eq, expEq := mighty.Eq(t), mighty.ExpEq(t) // t is of type *testing.T
	eq(6, len("mighty")) // Expect len("mighty") to be 6
	r := bytes.NewBuffer([]byte{'a'})
	expEq(byte('a'))(r.ReadByte()) // Expect to read 'a' and no error

### Example #1: testing `math.Abs()`

Without `mighty` it could look like this:

	cases := []struct{ in, exp float64 }{{1, 1}, {-1, 1}}
	for _, c := range cases {
		if got := math.Abs(c.in); got != c.exp {
			t.Errorf("Expected: %v, got: %v", c.exp, got)
		}
	}

Using `mighty`:

	cases := []struct{ in, exp float64 }{{1, 1}, {-1, 1}}
	eq := mighty.Eq(t)
	for _, c := range cases {
		eq(c.exp, math.Abs(c.in))
	}

### Example #2: testing reading from `bytes.Buffer`

Without `mighty` it could look like this:

	r := bytes.NewBufferString("test-data")
	if b, err := r.ReadByte(); b != 't' || err != nil {
		t.Errorf("Expected: %v, got: %v, error: %v", 't', b, err)
	}
	if line, err := r.ReadString('-'); line != "est-" || err != nil {
		t.Errorf("Expected: %v, got: %v, error: %v", "est", line, err)
	}
	p := make([]byte, 4)
	if n, err := r.Read(p); n != 4 || string(p) != "data" || err != nil {
		t.Errorf("Expected: n=%v, p=%v; got: n=%v, p=%v; error: %v",
			4, "data", n, string(p), err)
	}

Using `mighty`:

	eq, expEq := mighty.Eq(t), mighty.ExpEq(t)
	r := bytes.NewBufferString("test-data")
	expEq(byte('t'))(r.ReadByte())
	expEq("est-")(r.ReadString('-'))
	p := make([]byte, 4)
	expEq(4)(r.Read(p))
	eq("data", string(p))

### More examples

Fore more usage examples, check out the following packages which use `mighty` extensively in their testing code:

- [github.com/icza/bitio](https://github.com/icza/bitio)
- [github.com/icza/session](https://github.com/icza/session)
