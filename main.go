package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sklyar/glox/internal/scanner"
	"io"
	"log/slog"
	"os"
)

type Lox struct {
	hadError bool
	logger   *slog.Logger
}

func (l *Lox) RunFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.run(b)

	return nil
}

func (l *Lox) RunPrompt() {
	r := bufio.NewReader(os.Stdin)

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			l.logger.Error("parse line", err.Error())
			continue
		}
		l.run(line)
		l.hadError = false
	}
}

func (l *Lox) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) run(source []byte) {
	s := scanner.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		var generalError *scanner.GeneralError
		if errors.As(err, &generalError) {
			l.Error(generalError.Line, generalError.Error())
			return
		}
		var unexpectedCharacterError *scanner.UnexpectedCharacterError
		if errors.As(err, &unexpectedCharacterError) {
			l.report(unexpectedCharacterError.Line, "", unexpectedCharacterError.Error())
			return
		}
		return
	}
	for _, token := range tokens {
		fmt.Println(token)
	}

}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	lox := &Lox{logger: logger}

	args := os.Args[1:]
	if len(args) == 0 {
		lox.RunPrompt()
		return
	}

	if len(args) == 1 {
		if err := lox.RunFile(args[0]); err != nil {
			logger.Error("run file", err.Error())
			return
		}
		if lox.hadError {
			os.Exit(65)
		}
		return
	}

	println("Usage: lox [script]")
}
