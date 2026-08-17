package evaluator

import (
	"github.com/damascenov/monkey_interpreter/lexer"
	"github.com/damascenov/monkey_interpreter/object"
	"github.com/damascenov/monkey_interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"4", 4},
		{"7", 7},
		{"24", 24},
		{"-4", -4},
		{"-7", -7},
		{"-24", -24},
		{"4 + 7", 11},
		{"4 + 4 + 4", 12},
		{"4 * 7", 28},
		{"4 * 7 * 4", 112},
		{"-4 + 24 + -20", 0},
		{"4 + 7 * 24", 172},
		{"4 * 7 + 24", 52},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(4 + 7 * 24 + 5 / 3) * 3 + -10", 509},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"4 < 7", true},
		{"4 > 7", false},
		{"4 > 4", false},
		{"4 < 4", false},
		{"4 == 4", true},
		{"4 == 7", false},
		{"4 != 4", false},
		{"4 != 7", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"(4 < 7) == true", true},
		{"(4 < 7) == false", false},
		{"(4 > 7) == true", false},
		{"(4 > 7) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 4 }", 4},
		{"if (false) { 4 }", nil},
		{"if (1) { 7 }", 7},
		{"if (4 < 7) { 24 }", 24},
		{"if (4 > 7) { 24 }", nil},
		{"if (4 > 7) { 24 } else { 555 }", 555},
		{"if (4 < 7) { 24 } else { 555 }", 24},
	}

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        integer, ok := tt.expected.(int)
        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        } else {
            testNullObject(t, evaluated)
        }
    }
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
        {"return 4;", 4},
        {"return 7; 4;", 7},
        {"return 4 * 7;", 28},
        {"7; return 24; 4", 24},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func testNullObject(t *testing.T, obj object.Object) bool {
    if obj != NULL {
        t.Errorf("object is not NULL, got=%T (%+v)", obj, obj)
        return false
    }
    return true
}

func TestBangOperator(t *testing.T) {
	test := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range test {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}
