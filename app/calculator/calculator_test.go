package calculator

import (
	"context"
	"testing"

	"github.com/mdwhatcott/calcy-apps/app/commands"
	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture
	add     FakeCalculator
	sub     FakeCalculator
	mul     FakeCalculator
	div     FakeCalculator
	handler *Handler
}

func (this *HandlerFixture) Setup() {
	this.add = FakeCalculator{Offset: 100}
	this.sub = FakeCalculator{Offset: 200}
	this.mul = FakeCalculator{Offset: 300}
	this.div = FakeCalculator{Offset: 400}
	this.handler = NewHandler(this.add, this.sub, this.mul, this.div)
}
func (this *HandlerFixture) TestUnrecognizedCommand() {
	this.So(func() { this.handler.Handle(context.Background(), "unrecognized") }, should.Panic)
}
func (this *HandlerFixture) TestAdd() {
	command := &commands.Add{A: 1, B: 2}
	this.handler.Handle(context.Background(), command)
	this.So(command.Result.C, should.Equal, 103)
	this.So(command.Result.Error, should.BeNil)
}
func (this *HandlerFixture) TestSubtract() {
	command := &commands.Subtract{A: 1, B: 2}
	this.handler.Handle(context.Background(), command)
	this.So(command.Result.C, should.Equal, 203)
	this.So(command.Result.Error, should.BeNil)
}
func (this *HandlerFixture) TestMultiply() {
	command := &commands.Multiply{A: 1, B: 2}
	this.handler.Handle(context.Background(), command)
	this.So(command.Result.C, should.Equal, 303)
	this.So(command.Result.Error, should.BeNil)
}
func (this *HandlerFixture) TestDivide() {
	command := &commands.Divide{A: 1, B: 2}
	this.handler.Handle(context.Background(), command)
	this.So(command.Result.C, should.Equal, 403)
	this.So(command.Result.Error, should.BeNil)
}

type FakeCalculator struct{ Offset int }

func (this FakeCalculator) Calculate(a, b int) int {
	return this.Offset + a + b
}
