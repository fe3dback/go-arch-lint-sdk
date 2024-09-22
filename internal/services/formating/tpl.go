package formating

import _ "embed"

//go:embed tpl_syntax.gohtml
var tplNoticeSyntax []byte

//go:embed tpl_orphans.gohtml
var tplNoticeOrphans []byte

//go:embed tpl_imports.gohtml
var tplNoticeImports []byte
