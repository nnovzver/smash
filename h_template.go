package main

var h_template = `/*
 * Header for {{.Name}} codogram module
 * Auto-generated file! DO NOT MODIFY!
 */
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

{{range .Codograms}}
{{.CMacros}}
typedef struct {{.Name}} {
{{range .Fields}}{{if and (ne .Type getBlobId) (ne .Type getTempId)}}  {{.CType}} {{.Name}}; // {{.Length}} bits
{{else if eq .Type getBlobId}}  uint8_t {{.Name}}[{{bytesInBits .Length}}];
{{end}}{{end}}` +
	`} {{.Name}};

int Marshal_{{.Name}}({{.Name}} *c, void *buff, size_t size);
int Unmarshal_{{.Name}}({{.Name}} *c, void *buff, size_t size);
{{if .CTest}}int is{{.Name}}(void *buff, size_t size);{{end}}
{{end}}

#ifdef __cplusplus
}
#endif
`
