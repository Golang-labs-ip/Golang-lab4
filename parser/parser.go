package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	. "../command"
	. "../commands"
)

type tokenType int

const (
	num tokenType = iota
	str
	sym
)

type parseRule struct {
	constructor func(args []interface{}) Command
	rule        []tokenType
}

var parseTable = map[string]parseRule{
	"print":     {rule: []tokenType{str}, constructor: printConst},
	"sha1":      {rule: []tokenType{str}, constructor: sha1Const},
}

func printConst(args []interface{}) Command {
	arg := args[0].(string)
	return &PrintCommand{Arg: arg}
}

func sha1Const(args []interface{}) Command {
	arg := args[0].(string)
	return &Sha1Command{Arg: arg}
}

type Parser struct {
	input *bufio.Scanner
	line  int
}

func (p *Parser) parseLine(line string) Command {
	parts := strings.Fields(line)
	cmd := parts[0]

	rule, ok := parseTable[cmd]
	if !ok {
		return p.errorCommand(fmt.Sprintf("Unknown command: %s", cmd))
	}
	if len(rule.rule) == 1 && rule.rule[0] == str {
		// The command receives only one argument - the whole string
		start := len(cmd)
		rest := strings.Trim(line[start:], " ")
		return rule.constructor([]interface{}{rest})
	}
	parts = parts[1:]
	if len(parts) != len(rule.rule) {
		return p.errorCommand(fmt.Sprintf("Error in %s: expect %d args, got %d",
			cmd,
			len(rule.rule),
			len(parts)),
		)
	}
	return p.mathArgs(cmd, rule, parts)
}

func (p *Parser) mathArgs(cmd string, rule parseRule, args []string) Command {
	res := []interface{}{}
	for i := range rule.rule {
		switch rule.rule[i] {
		case num:
			n, err := strconv.Atoi(args[i])
			if err != nil {
				return p.errorCommand(fmt.Sprintf("Error in %s while parsing number: %s", cmd, err.Error()))
			}
			res = append(res, n)
		case str:
			res = append(res, args[i])
		case sym:
			sym := args[i]
			if len([]rune(sym)) != 1 {
				return p.errorCommand(fmt.Sprintf("Error in %s: got string instead of single character", cmd))
			}
			res = append(res, args[i])
		}
	}
	return rule.constructor(res)
}

func (p *Parser) Parse() []Command {
	res := []Command{}
	for p.input.Scan() {
		p.line++
		line := p.input.Text()
		// Skip empty lines
		if len(strings.Trim(line, " \t\n")) == 0 {
			continue
		}
		cmd := p.parseLine(line)
		res = append(res, cmd)
	}
	return res
}

func (p *Parser) parseNext() Command {
	p.line++
	p.input.Scan()
	line := p.input.Text()
	return p.parseLine(line)
}

func NewParser(reader io.Reader) Parser {
	return Parser{input: bufio.NewScanner(reader), line: 0}
}

func (p *Parser) errorCommand(msg string) Command {
	msg += fmt.Sprintf(" on line %d\n", p.line)
	return &ErrorCommand{Msg: msg}
}
