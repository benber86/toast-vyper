package syntax

import (
	"strings"
	"testing"
)

func TestVyperDetection(t *testing.T) {
	for _, path := range []string{"Vault.vy", "interfaces/Token.vyi", "VAULT.VY", "TOKEN.VYI"} {
		lang := ForPath(path)
		if lang == nil || lang.Name != "vyper" || lang != ForName("Vyper") {
			t.Fatalf("Vyper not registered for %s", path)
		}
	}
}

func TestVyperHighlight(t *testing.T) {
	src := `#pragma version 0.4.3
from ethereum.ercs import IERC20
initializes: controls[ownable := ownable]
uses: ownable
exports: controls.paused
implements: IERC20

event Deposit:
    sender: indexed(address)
    amount: uint256

struct Position:
    balance: uint256

flag Roles:
    ADMIN
    USER

interface Token:
    def balanceOf(owner: address) -> uint256: view
    def transfer(to: address, amount: uint256) -> bool: nonpayable

asset: public(immutable(address))
balances: HashMap[address, uint256]
LIMIT: constant(uint256) = 1_000

@deploy
def __init__(token: address):
    asset = token

@external
@nonreentrant
def deposit(amount: uint256) -> uint256:
    """Deposit tokens.
    A multiline docstring.
    """
    data: Bytes[4] = b"test"
    digest: bytes32 = keccak256(data)
    rate: decimal = 1.5
    active: bool = True
    balance: uint256 = staticcall Token(asset).balanceOf(self)
    assert extcall Token(asset).transfer(msg.sender, amount), "failed"
    for i: uint256 in range(3):
        balance += i
    log Deposit(sender=msg.sender, amount=amount)
    return balance
`
	h, err := NewHighlighter("Vault.vy", nil)
	if err != nil || !h.HasQuery() {
		t.Fatalf("Vyper query did not compile: %v", err)
	}
	defer h.parser.Close()
	defer h.query.Close()
	h.Parse([]byte(src))
	if h.tree == nil {
		t.Fatal("Vyper parser returned no tree")
	}
	if h.tree.RootNode().HasError() {
		t.Fatalf("Vyper parse failed: %v", h.tree.RootNode())
	}
	defer h.tree.Close()

	for _, tc := range []struct{ line, token, style string }{
		{"#pragma", "#pragma version 0.4.3", "comment"},
		{"initializes:", "initializes", "keyword"},
		{"initializes:", ":=", "operator"},
		{"uses:", "uses", "keyword"},
		{"exports:", "exports", "keyword"},
		{"implements:", "implements", "keyword"},
		{"event Deposit:", "event", "keyword"},
		{"event Deposit:", "Deposit", "type"},
		{"struct Position:", "Position", "type"},
		{"flag Roles:", "flag", "keyword"},
		{"flag Roles:", "Roles", "type"},
		{"    ADMIN", "ADMIN", "constant"},
		{"interface Token:", "Token", "type"},
		{"    def balanceOf", "balanceOf", "function"},
		{"    def balanceOf", "view", "keyword"},
		{"    def transfer", "nonpayable", "keyword"},
		{"    sender:", "indexed", "keyword"},
		{"asset:", "public", "keyword"},
		{"asset:", "immutable", "keyword"},
		{"asset:", "address", "type"},
		{"balances:", "HashMap", "type"},
		{"LIMIT:", "LIMIT", "constant"},
		{"LIMIT:", "1_000", "number"},
		{"@deploy", "@deploy", "attribute"},
		{"@external", "@external", "attribute"},
		{"@nonreentrant", "@nonreentrant", "attribute"},
		{"def deposit", "deposit", "function"},
		{"    A multiline", "    A multiline docstring.", "string"},
		{"    data:", "b\"test\"", "string"},
		{"    digest:", "bytes32", "type"},
		{"    digest:", "keccak256", "function"},
		{"    rate:", "1.5", "number"},
		{"    active:", "True", "constant"},
		{"    balance: uint256 =", "staticcall", "keyword"},
		{"    balance: uint256 =", "balanceOf", "function"},
		{"    assert extcall", "extcall", "keyword"},
		{"    for i:", "uint256", "type"},
		{"        balance +=", "+=", "operator"},
		{"    log Deposit", "log", "keyword"},
	} {
		t.Run(tc.line+"/"+tc.token, func(t *testing.T) {
			for row, line := range strings.Split(src, "\n") {
				if !strings.HasPrefix(line, tc.line) {
					continue
				}
				start := strings.Index(line, tc.token)
				for _, span := range h.HighlightLine(row, line) {
					if span.Start <= start && span.End >= start+len(tc.token) {
						if span.Style != tc.style {
							t.Errorf("%q style = %s, want %s", tc.token, span.Style, tc.style)
						}
						return
					}
				}
				t.Fatalf("no highlight for %q", tc.token)
			}
			t.Fatalf("line %q missing from fixture", tc.line)
		})
	}
}

func TestVyperPredicatesAndEdit(t *testing.T) {
	h, err := NewHighlighter("Token.vyi", nil)
	if err != nil || !h.HasQuery() {
		t.Fatalf("Vyper query did not compile: %v", err)
	}
	defer h.parser.Close()
	defer h.query.Close()
	src := "value: uint256 = amount\n"
	h.Parse([]byte(src))
	for _, s := range h.HighlightLine(0, src) {
		if src[s.Start:s.End] == "amount" && (s.Style == "type" || s.Style == "constant" || s.Style == "keyword") {
			t.Fatalf("ordinary identifier matched a predicate: %+v", s)
		}
	}
	updated := "value: uint256 = 42\n"
	oldTree := h.tree
	h.Edit([]byte(updated), 17, 22, 19, 0, 17, 0, 22, 0, 19)
	oldTree.Close()
	defer h.tree.Close()
	if h.tree.RootNode().HasError() {
		t.Fatal("incremental edit produced a parse error")
	}
	for _, s := range h.HighlightLine(0, updated) {
		if s.Style == "number" && updated[s.Start:s.End] == "42" {
			return
		}
	}
	t.Fatal("number not highlighted after incremental edit")
}
