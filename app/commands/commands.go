package commands

type Add struct {
	A int
	B int

	Result struct {
		C     int
		Error error
	}
}

type Subtract struct {
	A int
	B int

	Result struct {
		C     int
		Error error
	}
}

type Multiply struct {
	A int
	B int

	Result struct {
		C     int
		Error error
	}
}

type Divide struct {
	A int
	B int

	Result struct {
		C     int
		Error error
	}
}

type Bogus struct {
	A int
	B int

	Result struct {
		C     int
		Error error
	}
}
