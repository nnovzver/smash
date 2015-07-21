package main

var c_template = `/*
 * Implementation for {{.Name}} codogram module.
 * Auto-generated file! DO NOT MODIFY!
 */

#include <string.h>

#include "{{.Name}}.gen.h"

#define MASK(start, end) ( \
  (((1 << (start - end + 1)) - 1) << end) \
)

{{range .Codograms}}
int Marshal_{{.Name}}({{.Name}} *c, void *buff, size_t size) {
  char *ch = buff;
  if (size < {{.CLength}} || buff == NULL) return -1;

  memset(buff, 0, size);
  {{range .Fields}}
    {{if eq .Type getConstId}}
  c->{{.Name}} = {{.Const}};
    {{end}}
  {{end}}

{{.CMarshal}}

  return 0;
}

int Unmarshal_{{.Name}}({{.Name}} *c, void *buff, size_t size) {
  char *ch = buff;
  if (size < {{.CLength}} || buff == NULL) return -1;

  memset(c, 0, size);

{{.CUnmarshal}}

  return 0;
}

{{if .CTest}}
int is{{.Name}}(void *buff, size_t size) {
  {{range .Fields}}
    {{if eq .Type getConstId}}
  {{.CType}} {{.Name}} = 0;
    {{end}}
  {{end}}
  char *ch = buff;
  if (size < {{.CLength}} || buff == NULL) return -1;

{{.CTest}}

  if (
  {{range $index, $field := .Fields}}
    {{if eq .Type getConstId}}
      {{if $index}}&& {{.Name}} == {{.Const}}
      {{else}}{{.Name}} == {{.Const}}
      {{end}}
    {{end}}
  {{end}}
  ) {
    return 1;
  } else {
    return 0;
  }
}
{{end}}
{{end}}
`
