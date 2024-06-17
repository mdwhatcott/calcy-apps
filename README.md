# calcy-apps

The 'main' part of a study in polymorphic deployment of a package across multiple UIs.

See github.com/mdwhatcott/calcy-lib for the first part.

--------

## Toward a full-fledged 'Bounded Context', Smarty style

--------

### Module A: `github.com/smarty/assertions`

Purpose: More effective assertions functions for testing.

Rationale:

- As test authors, we should declare what we actually expect instead of check for the presence of the opposite, which
  seems to be the norm in many Go projects out there in the wild.
- We believe generic failure messages can still be very helpful and effective.

Bad:

```go
if actual != expected {
    t.Error("the test failed blah blah blah")
}
```

Good:

```go
should.So(t, actual, should.Equal, expected)
```

The above assertion reads: "So, `actual` should equal `expected`." Nice!

Instructions:

Step 1: Author a package at `externals/should` that implements the following elements:

```go
package should

type testingT interface {
	Helper()
	Error(...any)
}

type assertion func(actual any, expected ...any) error

func So(t testingT, actual any, assert assertion, expected ...any) bool // TODO

func Equal(actual any, EXPECTED ...any) error // TODO
func BeTrue(actual any, _ ...any) error       // TODO
func BeFalse(actual any, _ ...any) error      // TODO
func BeNil(actual any, _ ...any) error        // TODO

type negated struct{}

var NOT negated

func (negated) Equal(actual any, expected ...any) error // TODO
func (negated) BeNil(actual any, _ ...any) error        // TODO
```

Step 2: Rewrite assertions in your tests to use your new package, verifying that a assertions fail and pass as expected,
and that failure messages are helpful. Spoiler: github.com/mdwhatcott/tiny-should

If this is your first time implementing a testing tool, congratulations! You've taken your first step into a larger
world...

Step 3: Replace all usage of your new package with `github.com/smarty/assertions/should`

Step 4: Poke around github.com/smarty/assertions to get a feel for how the assertion functions work, maybe even run the
tests and goof around a bit.

--------

### Module B: `github.com/smarty/gunit`

Purpose: x-Unit test fixtures via a reflection-based, table-driven test runner

Rationale:

- Separate instances of a struct-based test fixture can encapsulate the elements and state involved in each related test
  case, providing common setup/teardown behavior and facilitating concurrent and random execution of test cases.

Bad (package-level test state):

```go
package something_test

var state ...

func TestCase1(t *testing.T) {
    state = ...
    actual = SystemUnderTest(state)
    should.So(t, actual, should.Equal, ...)
}
func TestCase2(t *testing.T) {
    state = ...
    actual = SystemUnderTest(state)
    should.So(t, actual, should.Equal, ...)
}
```

Good (struct-level test state):

```go
package something_test

func TestSystemUnderTestFixture(t *testing.T) {
	gunit.Run(new(SystemUnderTestFixture), t)
}

type SystemUnderTestFixture struct {
	*gunit.Fixture
	state ...
}

func (this *SystemUnderTestFixture) Setup() {
    this.state = ...
}

func (this *SystemUnderTestFixture) TestCase1() {
    actual = SystemUnderTest(this.state)
    this.So(actual, should.Equal, ...)
}
func (this *SystemUnderTestFixture) TestCase2() {
    actual = SystemUnderTest(this.state)
    this.So(actual, should.Equal, ...)
}
```

Instructions:

Step 1: Author a package at `externals/gunit` that implements the following elements:

```go
package gunit

func Run(t *testing.T, fixture any) // TODO

type Fixture struct { *testing.T }

func (this *Fixture) So(actual any, assert assertion, expected ...any) // TODO

type assertion func(actual any, expected ...any) error
```

The `Run` func is the most difficult part. You must use the `reflect` package to scan the provided `fixture` for
a `Setup` method and any `Test...` methods. For each `Test...` method, instantiate a new instance (with reflection) of
the fixture type, call the `Setup` method, then call the `Test...` method.

Spoiler: https://www.smarty.com/blog/lets-build-xunit-in-go

Step 2: Rewrite the test cases in this project to use `gunit` fixtures, ensuring that each test case executes and that
failures are reported correctly.

Step 3: Replace usage of your new package with github.com/smarty/gunit.

Step 4: Poke around github.com/smarty/gunit to get a feel for how the package is laid out, maybe even run the tests and
goof around a bit.

--------

### Module C: `github.com/smarty/shuttle`

Purpose: Transform HTTP requests into intention-revealing, user instructions. After processing the given operation,
render the results of that operation back to the underlying HTTP response.

Rationale:

- Application logic should be kept very separate from transport protocols.

Bad:

```go
func (this *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    operand1, err := strconv.Atoi(request.Form.Get("a"))
    if err != nil {
        http.Error(response, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
        return
    }
    operand2, err := strconv.Atoi(request.Form.Get("b"))
    if err != nil {
        http.Error(response, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
        return
    }
    result := operand1 + operand2

    response.Header().Set("Content-Type", "text/plain; charset=utf-8")
    response.WriteHeader(http.StatusOK)
    _, _ = io.WriteString(response, strconv.Itoa(result))
}
```

Good:

```go
type Processor struct {
	handler contracts.Handler
}

func (this *Processor) Process(ctx context.Context, v any) any {
    switch input := v.(type) {
    case *inputs.Addition:
		return this.add(ctx, input)
    ...
    }
}
func (this *Processor) add(ctx context.Context, input *inputs.Addition) any {
	command := &commands.Add{A: input.A, B: input.B}
	this.handler.Handle(ctx, command)
	if command.Result.Error != nil {
		return additionFailure
	}
	return views.Addition{A: input.A, B: input.B, C: command.Result.C}
}

var additionFailure = shuttle.SerializeResult{
    StatusCode: http.StatusInternalServerError,
    Content: shuttle.InputError{
        Fields:  []string{"query:a", "query:b"},
        Name:    "calculation:addition-error",
        Message: fmt.Sprintf("The operands could not be added", verbPastParticiple),
    },
}
```

There's a lot to notice here:

1. There is no mention of `ServeHTTP(...)`, `*http.Request`, or `http.ResponseWriter` anywhere in the 'good' code.
2. There is no mention of http status codes, or response headers, etc...
3. There are more moving parts:
    - `*inputs.Addition` seems to hold the input operands (how did those get parsed from the `*http.Request`?)
    - `*commands.Add` seems to represent the user's intention that the application add the two operands from the input
    - `this.handler.Handle(...)` receives the command (and populates an error field).
    - `additionFailure` is an interesting data structure with all sorts of stuff...
    - `views.Addition` is another data structure which must, at some point, get serialized to the http response (as
      JSON).

We'll work in small steps to move toward this approach.

Step 1 (commands): define a new package called `app/commands` with a single file called `commands.go`. Define four
command structures called `Add`, `Subtract`, `Multiply`, and `Divide`. Here's what they should look like:

```go
type Add struct {
	A int
	B int

	Result struct {
		C     int
		Error error
	}
}
```

Step 2 (app handler): define a new package called `app/calculator`, containing `calculator.go`,
and `calculator_test.go` (two new Go files). Add the following to `calculator.go`:

```go
type Calculator interface{ Calculate(a, b int) int }

type Handler struct{ add, sub, mul, div Calculator }

func (this *Handler) Handle(ctx context.Context, commands ...any) // TODO
```

The purpose of the `Handler` (with its `Handle` method) is to receive incoming commands and fulfil the user's intention
by using the `Calculator` instances it received in its constructor to supply results back to each command.

Step 3 (app handler tests): Write a test suite in `calculator_test.go` as you implement the `Handler`
in `calculator.go`, which will prove that your handler implementation uses the correct `Calculator` for each
supplied `command` and assigns the calculated result on the command. Rather than using the actual implementations
of `Calculator` from your library module, define a `FakeCalculator` struct that implements the `Calculator` interface
for use in testing. We won't actually be assigning to the `Error` field on the commands' `Result` structure in this
example, but in the 'real world', when the application fails to carry out the command, it would set an error value.

At this point, we've moved the application logic quite far from the HTTP components. We are off to a good start!

Step 4 (http input models): In order to encourage separation of HTTP stuff from Application logic, this approach
provides one, and only one, opportunity to extract (or 'bind') data from the `*http.Request`: the "input model". So,
let's get started by creating a package at `http/inputs` with two files (to start with): `addition.go`
and `addition_test.go`. Start with the following in `addition.go`:

```go
type Addition struct {
	A int
	B int
}

func (this *Addition) Bind(request *http.Request) error
```

Implement the `Bind` method such that when invoked, it gets the query string parameters `a` and `b` from the
provided `request` and parses them as integers, setting the parsed values to the `A` and `B` fields of `this`. If
parsing fails, return something like this:

```go
InputError{
   Fields:  []string{fmt.Sprintf("query:%s", key)},
   Message: fmt.Sprintf("failed to parse '%s' parameter as integer: [%s]", key, raw),
}
```

(You'll need to define that `InputError` struct.)

Use `addition_test.go` to define a test suite of your own creation that proves the `Addition` struct's `Bind` method
works as described above.

Once you've got the `Addition` input model and tests, replicate that for `Subtraction`, `Multiplication`,
and `Division`. Yay, so much fun! (It will get tempting to wonder why we don't just use a consolidated input model data
structure since these operations all work with similar data. We're trying to give you a picture of several different
pathways through a system, where each pathway has its own input, output, and corresponding messages. Hang in there!)

Step 5 (http output views): Most of the output we produce from HTTP APIs is JSON. In Go, structs with JSON tags are very
easy to serialize, so that's what we return from our shuttle processors. Let's build the structures we'll use to return
JSON data to http responses. Create a package at `http/views` with data structures that look like this:

```go
type Addition struct {
	A int `json:"a"`
	B int `json:"b"`
	C int `json:"c"`
}
...
```

That's it!

Step 6 (processor): The processor is what receives the input model (after `Bind()` was called, and without returning any
error) and then translates that input model to a command (defined near the application logic) and sends it to the
application for processing. The value returned by the processor depends on what happens to the command. You can see a
pretty good example of what to build in the "Good" code snippet above (just before 'Step 1' of this section). Of course,
you won't forget to build a test suite to make sure that when the processor receives a populated input model it 1) turns
it into a command that 2) gets fed to the application handler (which will be a fake), which will 3) assign a result or
error, so that 4) the process can interpret the result and provide a return value.

Step 7 (shuttle in a bottle): Of course, we need some bit of library code to feed the `*http.Request` to the input
model's `Bind()` method, and then to feed the input model to the `Processor`, and then to take what the `Processor`
returns and write to the `http.ResponseWriter`. That's what we're going to build now, in a new package
at `externals/shuttle`. Here's a start:

```go
type (
	InputModel interface {
		Bind(request *http.Request) error
	}
	InputError struct {
		Fields  []string `json:"fields,omitempty"`
		Name    string   `json:"name,omitempty"`
		Message string   `json:"message,omitempty"`
	}
)

func (this InputError) Error() string // TODO

type (
	Processor interface {
		Process(ctx context.Context, v any) any
	}
	SerializeResult struct {
		StatusCode int
		Content    any
	}
)

func NewHandler(input func() InputModel, processor func() Processor) http.Handler
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	    // TODO:
	    // assign new input model
	    // bind request to input model (if err, return http.StatusBadRequest)
	    // pass input model to the processor, serialize result as JSON to http response
	})
}
```

Step 8 (routes and main):

Now for the fun part:

Main:

```go
func main() {
	appHandler := calculator.NewHandler(
		calcy.Addition{},
		calcy.Subtraction{},
		calcy.Multiplication{},
		calcy.Division{},
	)
	endpoint := "localhost:8080"
	log.Println("Listening on", endpoint)
	err := http.ListenAndServe(endpoint, HTTP.Router(appHandler))
	if err != nil {
		log.Fatalln(err)
	}
}
```

Routes:

```go
func Router(calculator contracts.Handler) http.Handler {
	h := http.NewServeMux()
	processor := func() shuttle.Processor { return NewProcessor(calculator) }
	h.Handle("/add", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewAddition() }, processor))
	h.Handle("/sub", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewSubtraction() }, processor))
	h.Handle("/mul", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewMultiplication() }, processor))
	h.Handle("/div", shuttle.NewHandler(func() shuttle.InputModel { return inputs.NewDivision() }, processor))
	return h
}
```

Try running the application and sending a few curl requests in!

Step 9 (drop-in smarty shuttle): Replace usage of your `shuttle` package with github.com/smarty/shuttle

Step 10: Plumb the layers and depths of github.com/smarty/shuttle. Find code that corresponds with each of your more simple `shuttle` code. 

Step 11: Pat yourself on the back. That was a lot of moving parts!

-------------

### Module D: `github.com/smarty/httprouter`

Purpose: Fast, flexible routing of incoming HTTP requests based on request method and path (which can include wildcard elements).

Rationale:

- A file path is a representation of a tree structure, so let's leverage that kind of data structure. The built-in http request router (http.ServeMux) has been pretty limited up until a very recent version of Go, and it's still a bit more loosey-goosey than we'd prefer.

Instructions:

Step 1: Implement a package at `ext/httprouter` with the following elements:

```go
func New(routes ...Route) (http.Handler, error)

type router struct {
	root *treeNode
}

func (this *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.root.Resolve(request.Method, request.URL.Path).ServeHTTP(response, request)
}

type treeNode struct {
	pathElement string
	static      []*treeNode // FUTURE: wildcard and variable nodes...
	handlers    *methodHandlers
}

func (this *treeNode) Add(route Route) error
func (this *treeNode) Resolve(method, path string) http.Handler
func notFoundHandler(response http.ResponseWriter, _ *http.Request)

type methodHandlers struct {
	get http.Handler // FUTURE: other methods
}

func (this *methodHandlers) Add(method string, handler http.Handler) bool
func (this *methodHandlers) Resolve(method string) http.Handler
func methodNotAllowedHandler(responseWriter http.ResponseWriter, _ *http.Request) 

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

func ParseRoute(method string, path string, handler http.Handler) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
```

The idea here is that each slash-separated element of a path is represented as a level/node in a tree structure. We parse a registered route at startup to create the tree structure. Then at runtime we traverse the tree according to the incoming request path elements. If we found a matching terminal node, we serve the response from it, otherwise serve http 404 (not found) or, of the path matches, but the method doesn't match, serve http 415 (method not allowed). So, in summary, you'll need to process the path in slash-separate elements to construct and traverse a tree. (Don't forget to write tests.)

Step 2: Install your new router in `http/routes.go`

Step 3: (drop-in smarty httprouter): Replace usage of your `httprouter` package with github.com/smarty/httprouter

Step 4: explore the code of smarty/httprouter and learn how it supports wildcard and variable path elements.

### Module E: `github.com/smarty/httpstatus`

Purpose: Expose an HTTP route at `/status` that communicates the readiness of an application to serve requests.

Rationale: Allow operations team to eventually deploy new versions of software alongside older versions to facilitate zero-downtime rollouts. The HTTP handler that responds to the `/status` request will respond with one of four status responses:

1. Starting
2. Healthy
3. Failing
4. Stopping

At startup the handler will be in "Starting" mode. As a background goroutine succeeds in pinging some important resource (maybe a database), it will upgrade the mode to "Healthy". The status check will be repeated at regular intervals. If ever the check fails, the mode will transition to "Failing" until a check succeeds again. When the application is in process of shutting down, the mode will transition to "Stopping". Any mode but "Healthy" will result in HTTP 503 Service Unavailable.

Instructions:

Step 1: Implement a package at `/ext/httpstatus` with the following elements:

```go
type HealthCheck interface {
	Status(ctx context.Context) error
}

type Handler struct {
	state         uint32
	hardContext   context.Context
	softContext   context.Context
	shutdown      context.CancelFunc
	healthCheck   HealthCheck
	timeout       time.Duration
	frequency     time.Duration
	shutdownDelay time.Duration
}

func NewHandler(ctx context.Context, check HealthCheck, timeout, frequency, shutdownDelay time.Duration) *Handler

func (this *Handler) ServeHTTP(response http.ResponseWriter, _ *http.Request)
func (this *Handler) Listen()
func (this *Handler) Close() error

const (
	stateStarting = iota
	stateHealthy
	stateFailing
	stateStopping
)
```

Considerations:

- At first all `ServeHTTP` needs to do is write the text "Starting" to the http response as plain text, which is what you'll focus on for the first unit test.
- From there, things get interesting. The `Listen` method is long-lived and will be called from a different goroutine. It will run until the provided context is cancelled, which may happen as a result of `Close` being called. Once running the `Listen` method will call the provided `HealthCheck` and transition the `state` field accordingly. The `context.Context` provided will be used to create a derived/child context which will be passed to the `HealthCheck`. In the event that the health check returns no error value, transition to the 'healthy' state. In the event that the error represents a context cancellation (perhaps because of a timeout), transition to 'Stopping' and return from `Listen`. In the event that the error is not nil (and not a context cancellation), transition to 'Failing'. In all cases except for transitioning to 'Stopping', sleep for the provided `delay` before repeating the health check.
- Testing suggestion: use very small time duration values for the timeout, frequency, and delay fields.
- Testing suggestion: don't execute tests with the `-race` flag at first. (See 'Note' in next bullet point) 
- NOTE: Because the `Listen` and `ServeHTTP` methods refer to the `state` field from different goroutines there is a very real possibility of a data race, creating undefined behavior (most likely a program crash). Use atomic operations or a mutex to protect against such an unpleasant outcome. Once that solution is in place, calling `go test` with `-race` should pass. Oh, and if your test cases reference any state over multiple goroutines you'll need to install similar atomic/mutex treatment there too. 
- **Note to mentor:** While not very much behavior, this is tricky stuff. It may be more effective to pair program with the mentee and to refer often to github.com/smarty/httpstatus to really grasp the handling of the context.Context values, the atomic operations on the state field, and how the tests leverage the monitor interface to make the tests more deterministic given the concurrent nature of this component.

Step 2: Install a new HTTP route at `/status` that points to the `Handler` in your new `/ext/httpstatus` package.

Step 3: (drop-in smarty httpstatus): Replace usage of your `httpstatus` package with github.com/smarty/httpstatus (this will require getting to know various functional options).

Step 4: explore the code of smarty/httpstatus and learn how it precomputes the 4 status handlers, as well as how it communicates with a monitor interface.
