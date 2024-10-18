package http

import (
	"context"
	"errors"
	"testing"

	"github.com/mdw-smarty/calc-apps/app/commands"
	"github.com/mdw-smarty/calc-apps/http/inputs"
	"github.com/mdw-smarty/calc-apps/http/views"
	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestProcessorFixture(t *testing.T) {
	gunit.Run(new(ProcessorFixture), t)
}

type ProcessorFixture struct {
	*gunit.Fixture
	processor *Processor
	ctx       context.Context
	handler   *FakeApplicationHandler
}

func (this *ProcessorFixture) Setup() {
	this.ctx = context.WithValue(context.Background(), "test", this.Name())
	this.handler = NewFakeApplicationHandler()
	this.handler.ctx = func(ctx context.Context) {
		this.So(ctx.Value("test"), should.Equal, this.Name())
	}
	this.processor = NewProcessor(this.handler)
}

func (this *ProcessorFixture) TestUnrecognizedInput() {
	result := this.processor.Process(this.ctx, 42)
	this.So(result, should.Equal, internalServerError)
}

func (this *ProcessorFixture) TestAddition_ErrorFromHandler() {
	this.handler.err = errors.New("boink")
	result := this.processor.Process(this.ctx, &inputs.Addition{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Add{A: 1, B: 2}})
	this.So(result, should.Equal, additionFailure)
}
func (this *ProcessorFixture) TestSubtraction_ErrorFromHandler() {
	this.handler.err = errors.New("boink")
	result := this.processor.Process(this.ctx, &inputs.Subtraction{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Subtract{A: 1, B: 2}})
	this.So(result, should.Equal, subtractionFailure)
}
func (this *ProcessorFixture) TestMultiplication_ErrorFromHandler() {
	this.handler.err = errors.New("boink")
	result := this.processor.Process(this.ctx, &inputs.Multiplication{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Multiply{A: 1, B: 2}})
	this.So(result, should.Equal, multiplicationFailure)
}
func (this *ProcessorFixture) TestDivision_ErrorFromHandler() {
	this.handler.err = errors.New("boink")
	result := this.processor.Process(this.ctx, &inputs.Division{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Divide{A: 1, B: 2}})
	this.So(result, should.Equal, divisionFailure)
}

func (this *ProcessorFixture) TestAddition() {
	result := this.processor.Process(this.ctx, &inputs.Addition{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Add{A: 1, B: 2}})
	this.So(result, should.Equal, views.Addition{A: 1, B: 2, C: 42})
}
func (this *ProcessorFixture) TestSubtraction() {
	result := this.processor.Process(this.ctx, &inputs.Subtraction{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Subtract{A: 1, B: 2}})
	this.So(result, should.Equal, views.Subtraction{A: 1, B: 2, C: 42})
}
func (this *ProcessorFixture) TestMultiplication() {
	result := this.processor.Process(this.ctx, &inputs.Multiplication{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Multiply{A: 1, B: 2}})
	this.So(result, should.Equal, views.Multiplication{A: 1, B: 2, C: 42})
}
func (this *ProcessorFixture) TestDivision() {
	result := this.processor.Process(this.ctx, &inputs.Division{A: 1, B: 2})
	this.So(this.handler.handled, should.Equal, []any{commands.Divide{A: 1, B: 2}})
	this.So(result, should.Equal, views.Division{A: 1, B: 2, C: 42})
}
