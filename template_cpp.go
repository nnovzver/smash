package main

var cpp_template = `/*
 * CPP implementation for {{.Name}} codogram module
 * Auto-generated file! DO NOT MODIFY!
 */

#include "{{.FileName}}.gen.cw.hpp"

{{range .Codograms}}
Codograms::{{.Name}}::{{.Name}}()
{
  buf.fill(0, {{.Name}}__BUFSIZE);
  clearMessage();
}

const size_t Codograms::{{.Name}}::bufsize = {{.Name}}__BUFSIZE;

bool Codograms::{{.Name}}::marshal()
{
  if (Marshal_{{.Name}}(&m, buf.data(), buf.size()))
    return false;
  else
    return true;
}

bool Codograms::{{.Name}}::unmarshal()
{
  if (Unmarshal_{{.Name}}(&m, buf.data(), buf.size()))
    return false;
  else
    return true;
}

bool Codograms::{{.Name}}::checkBuf()
{
  if (Is_{{.Name}}(buf.data(), buf.size()))
    return true;
  else
    return false;  
}

void Codograms::{{.Name}}::clearMessage()
{
  memset(&m, 0, sizeof(m));
}

int Codograms::{{.Name}}::msize()
{
  return sizeof(m);
}
{{end}}
`
