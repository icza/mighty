# mighty

Lightweight extension to Go's [testing](https://golang.org/pkg/testing/) package.

With this utility, you can make your Go test files super clean but still intuitive.
This package doesn't want to hide complex things from you, but wants to eliminate repetitive,
long code from your tests.

## Using `mighty`

You could create a value of `mighty.Myt` and use its methods like `Myt.Eq()`, `Myt.Neq()`, ... etc. as seen below:

	m := mighty.Myt{t}
	m.Eq(1, Sign(5)) // Expect Sign(5) to return 1
	m.Eq(0, Sign(0)) // Expect Sign(0) to return 0

But the recommended way is to acquire [method values](https://golang.org/ref/spec#Method_values) returned by functions of `mighty`:

	eq := mighty.Eq(t)
	eq(1, Sign(5)) // Expect Sign(5) to return 1
	eq(0, Sign(0)) // Expect Sign(0) to return 0

### Example #1: testing the `Sign(int) int` function

Which returns `1` if argument is positive, `-1` if it is negative, `0` otherwise.
Without `mighty` it could look like this:

	func TestSign(t *testing.T) {
		cases := []struct{ in, exp int }{{5, 1}, {-5, -1}, {0, 0}}
		for _, c := range cases {
			if got := Sign(c.in); got != c.exp {
				t.Errorf("Expected: %v, got: %v", c.exp, got)
			}
		}
	}

Using `mighty`:

	func TestSign(t *testing.T) {
		cases := []struct{ in, exp int }{{5, 1}, {-5, -1}, {0, 0}}
		eq := mighty.Eq(t)
		for _, c := range cases {
			eq(c.exp, Sign(c.in))
		}
	}

### Example #2: testing reading from `bytes.Buffer`

Without `mighty` it could look like this:

	r := bytes.NewBufferString("test-data") // Acquire the Buffer
	if b, err := r.ReadByte(); b != 't' || err != nil {
		t.Errorf("Expected: %v, got: %v, error: %v", 't', b, err)
	}
	if line, err := r.ReadString('-'); line != "est-" || err != nil {
		t.Errorf("Expected: %v, got: %v, error: %v", "est", line, err)
	}
	p := make([]byte, 4)
	if n, err := r.Read(p); n != 4 || string(p) != "data" || err != nil {
		t.Errorf("Expected: n=%v, p=%v; got: n=%v, p=%v; error: %v", 4, "data", n, string(p), err)
	}

Using `mighty`:

	eq, expEq := mighty.Eq(t), mighty.ExpEq(t)
	r := bytes.NewBufferString("test-data") // Acquire the reader
	expEq(byte('t'))(r.ReadByte())
	expEq("est-")(r.ReadString('-'))
	p := make([]byte, 4)
	expEq(4)(r.Read(p))
	eq("data", string(p))
