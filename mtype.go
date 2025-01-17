// If used, DeclaredDoctype and GuessedDoctype should be mutually exclusive.
//
// Reference material re. MDITA:
//  - https://github.com/jelovirt/dita-ot-markdown/wiki/Syntax-reference
//  - The format of local link targets is detected based on file extension.
//  - The following extensions are treated as DITA files:
//    `.dita` => dita ; `.xml` => dita ; `.md` => markdown ;
//    `.markdown` => markdown

package xmlutils

// An MType is specific to this app and/but is modeled after
// the prior concept of Mime-type. An MType has three fields.
//
// Its value is generally based on two to four inputs:
//   - The Mime-type guess returned by Go stdlib
//     func net/http.DetectContentType(data []byte) string
//     (which is based on https://mimesniff.spec.whatwg.org/ )
//     (The no-op default return value is "application/octet-stream")
//   - Our own shallow analysis of file contents
//   - The file extension (it is normally present)
//   - The DOCTYPE (iff XML, incl. HTML)
//
// Note that
//   - a plain text file MAY be presumed to be Markdown, altho it
//     is not clear (yet) which (TXT or MKDN) should take precedence.
//   - a Markdown file CAN and WILL be presumed to be LwDITA MDITA;
//     this may cause conflicts/problems for other dialects.
//   - mappings can appear bogus, for example HTTP stdlib "text/html"
//     might become MType "xml/html".
//
// String possibilities (but in LOWER CASE!) in each field:
//
//   - [0] XML, HTML, BIN, TXT, MKDN, (new!) DIRLIKE (i.e. non-contentful)
//   - We might (or not) keep XML and HTML distinct for a number of 
//     reasons, but partly because in the Go stdlib, they are processed 
//     quite differently, and we take advantage of it to keep HTML pro-
//     cessing free of nasty surprises and unhelpful strictness
//   - We might (or might not) keep MKDN distinct from TXT
//   - [1] CNT (Content), MAP (ToC), IMG, SCH(ema) [and maybe others TBD?]
//   - [2] Depends on [0]:
//     XML: per-DTD [and/or Pub-ID/filext];
//     HTML: per-DTD [and/or Pub-ID/filext];
//     BIN: format/filext;
//     SCH: format/filext [DTD,MOD,XSD,wotevs];
//     TXT: TBD
//     MKDN: flavor of Markdown (?) (note also AsciiDoc, RST, ...)
//     DIRLIKE: dir, symlink, pipe, socket, ...?
//
// Possible FIXME: Let [2] (3rd) be version info (html5, lwdiat, dita13)
// and then keep root tag info separately.
//
// Possible FIXME: Append version info, probably after a semicolon.
// .
type MType string
