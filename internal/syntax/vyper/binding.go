// Package vyper provides the Vyper Tree-sitter grammar for Toast.
package vyper

// #include "tree_sitter/parser.h"
// const TSLanguage *tree_sitter_vyper(void);
import "C"

import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	return sitter.NewLanguage(unsafe.Pointer(C.tree_sitter_vyper()))
}
