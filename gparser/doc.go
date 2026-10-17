//go:generate go run genpos.go

// Package gparser is a goyacc-based parser for OceanBase 4.x MySQL-mode SQL.
// It is a port of the official Bison grammar (sql_parser_mysql_mode.y) at the
// statement level, sharing the lexer and AST with the rest of the library.
package gparser
