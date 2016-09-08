# mighty

Lightweight extension to Go's [testing](https://golang.org/pkg/testing/) package.

With this utility, you can make your Go test files super clean but still intuitive.
This package doesn't want to hide complex things from you, but wants to eliminate repetitive,
long code from your tests.

Example: Let's write testing code for a `Sign(int) int` function,
which returns `1` if argument is positive, `-1` if it is negative, `0` otherwise.
It could look like this:

	func TestSign(t *testing.T) {
		cases := []struct{ in, exp int }{{12, 1}, {-12, -1}, {0, 0}}
		for _, c := range cases {
			if got := Sign(c.in); got != c.exp {
				t.Errorf("Expected: %v, got: %v", c.exp, got)
			}
		}
	}

Using `mighty`:

	func TestSign(t *testing.T) {
		cases := []struct{ in, exp int }{{12, 1}, {-12, -1}, {0, 0}}
		eq := mighty.Eq(t)
		for _, c := range cases {
			eq(c.exp, Sign(c.in))
		}
	}
