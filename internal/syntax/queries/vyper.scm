; Vyper-specific grammar; captures use Toast's existing theme keys.
(comment) @comment
(string) @string
(integer) @number
(float) @number
[(true) (false) (none)] @constant

["and" "as" "assert" "break" "continue" "def" "elif" "else"
 "for" "from" "if" "import" "in" "is" "not" "or" "pass"
 "raise" "return" "event" "struct" "interface" "enum" "flag"
 "log" "extcall" "staticcall"] @keyword

(function_definition name: (identifier) @function)
(interface_sig name: (identifier) @function mutability: (identifier) @keyword)
[(struct_definition name: (identifier) @type)
 (event_definition name: (identifier) @type)
 (interface_definition name: (identifier) @type)
 (enum_definition name: (identifier) @type)
 (flag_definition name: (identifier) @type)]

(decorator) @attribute
(type (identifier) @type)
(generic_type (identifier) @type)

((identifier) @type
 (#match? @type "^(bool|address|decimal|String|Bytes|HashMap|DynArray|u?int(8|16|24|32|40|48|56|64|72|80|88|96|104|112|120|128|136|144|152|160|168|176|184|192|200|208|216|224|232|240|248|256)|bytes(1|2|3|4|5|6|7|8|9|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|32))$"))

((call function: (identifier) @keyword)
 (#match? @keyword "^(public|immutable|constant|transient|indexed)$"))

((call function: (identifier) @function)
 (#not-match? @function "^(bool|address|decimal|String|Bytes|HashMap|DynArray|u?int(8|16|24|32|40|48|56|64|72|80|88|96|104|112|120|128|136|144|152|160|168|176|184|192|200|208|216|224|232|240|248|256)|bytes(1|2|3|4|5|6|7|8|9|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|32))$")
 (#not-match? @function "^(public|immutable|constant|transient|indexed)$"))
(call function: (attribute attribute: (identifier) @function))

((assignment left: (identifier) @keyword)
 (#match? @keyword "^(implements|uses|initializes|exports)$"))

((identifier) @constant
 (#match? @constant "^[A-Z][A-Z_0-9]*$"))

["+" "-" "*" "/" "//" "%" "**" "&" "|" "^" "~" "<<" ">>"
 "==" "!=" "<" ">" "<=" ">=" "=" ":=" "+=" "-=" "*=" "/=" "//="
 "%=" "**=" "&=" "|=" "^=" "<<=" ">>="] @operator
["(" ")" "[" "]" "{" "}" "," "." ";" ":" "->"] @punctuation
(ellipsis) @punctuation
