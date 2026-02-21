//go:build embed_notice

package main

import _ "embed"

//go:embed NOTICE
var noticeContent string