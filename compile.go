package ctxrouter

import (
	"strings"
)

// Template is a compiled representation of path templates.
type Template struct {
	// OpCodes is a sequence of operations.
	OpCodes []int
	// Pool is a constant pool
	Pool []string
	// Verb is a VERB part in the template.
	Verb string
	// Fields is a list of field paths bound in this template.
	Fields []string
	// Original template (example: /v1/a_bit_of_everything)
	Template string
}

// Compiler compiles utilities representation of path templates into marshallable operations.
// They can be unmarshalled by runtime.NewPattern.
type Compiler interface {
	Compile() Template
}

type cop struct {
	// code is the opcode of the operation
	code OpCode

	// str is a string operand of the code.
	// num is ignored if str is not empty.
	str string

	// num is a numeric operand of the code.
	num int
}

func (w wildcard) compile() []cop {
	return []cop{
		{code: OpPush},
	}
}

func (w deepWildcard) compile() []cop {
	return []cop{
		{code: OpPushM},
	}
}

func (l literal) compile() []cop {
	return []cop{
		{
			code: OpLitPush,
			str:  string(l),
		},
	}
}

func (v variable) compile() []cop {
	var ops []cop
	for _, s := range v.segments {
		ops = append(ops, s.compile()...)
	}
	ops = append(ops, cop{
		code: OpConcatN,
		num:  len(v.segments),
	}, cop{
		code: OpCapture,
		str:  v.path,
	})

	return ops
}

func (t template) Compile() Template {
	var rawOps []cop
	for _, s := range t.segments {
		rawOps = append(rawOps, s.compile()...)
	}
	var (
		ops    []int
		pool   []string
		fields []string
	)
	hasEOF := strings.HasSuffix(t.template, "/")

	consts := make(map[string]int)
	for _, op := range rawOps {
		ops = append(ops, int(op.code))
		if op.str == "" && !hasEOF {
			ops = append(ops, op.num)
		} else {
			if _, ok := consts[op.str]; !ok {
				consts[op.str] = len(pool)
				pool = append(pool, op.str)
			}
			ops = append(ops, consts[op.str])
		}
		if op.code == OpCapture {
			fields = append(fields, op.str)
		}
	}
	return Template{
		OpCodes:  ops,
		Pool:     pool,
		Verb:     t.verb,
		Fields:   fields,
		Template: t.template,
	}
}
