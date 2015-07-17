package main

var h_template = `/*
 * Header for {{.Name}} codogram module
 * Auto-generated file! DO NOT MODIFY!
 */
#include <stdint.h>

{{range .Codograms}}
typedef struct {{.Name}} {
  {{range .Fields}}
    {{.CType}} {{.Name}}; // {{.Length}} bits
  {{end}}
} {{.Name}};

int Marshal_{{.Name}}({{.Name}} *c, void *buff, size_t size);
int Unmarshal_{{.Name}}({{.Name}} *c, void *buff, size_t size);
int is{{.Name}}(void *buff, size_t size);
{{end}}
`
