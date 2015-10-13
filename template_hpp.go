package main

var hpp_template = `/*
 * CPP header for {{.Name}} codogram module
 * Auto-generated file! DO NOT MODIFY!
 */
#ifndef {{.FileName}}_GEN_CW_HPP
#define {{.FileName}}_GEN_CW_HPP

#include <QByteArray>

#include "{{.FileName}}.gen.h"

namespace Codograms {

#ifndef CLASS_CODOGRAM
#define CLASS_CODOGRAM
class Codogram {
public:
  virtual bool marshal() = 0;
  virtual bool unmarshal() = 0;
  virtual bool checkBuf() = 0;
  virtual void clearMessage() = 0;
  virtual const QByteArray& getBuf() = 0;
  virtual ~Codogram() {};
};
#endif

{{range .Codograms}}
class {{.Name}} : public Codogram {
public:
  {{.Name}}();
  bool marshal();
  bool unmarshal();
  bool checkBuf();
  const QByteArray& getBuf();
  void clearMessage();
  int msize();

  ::{{.Name}} m;
  QByteArray buf;
  static const size_t bufsize;
};
{{end}}

}
#endif // {{.FileName}}_GEN_CW_HPP
`
