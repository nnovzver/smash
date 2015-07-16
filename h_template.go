package main

var h_template = `/*
 * Header for {{.Name}} codogram module
 * Auto-generated file! DO NOT MODIFY!
 */
{{range .Codograms}}
typedef struct {{.Name}} {
  {{range .Fields}}
    {{.CType}} {{.Name}}; // {{.Length}} bits
  {{end}}
} {{.Name}};
{{end}}
`