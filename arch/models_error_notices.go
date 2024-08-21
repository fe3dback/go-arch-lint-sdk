package arch

import (
	"bytes"
	"fmt"
)

type (
	ErrorWithNotices struct {
		OverallMessage           string
		Notices                  []Notice
		injectNoticesToErrorText bool
	}

	Notice struct {
		Message   string
		Reference Reference
	}
)

func NewErrorWithNotices(overallMessage string, notices []Notice, injectNoticesToErrorText bool) *ErrorWithNotices {
	return &ErrorWithNotices{
		OverallMessage:           overallMessage,
		Notices:                  notices,
		injectNoticesToErrorText: injectNoticesToErrorText,
	}
}

func (en ErrorWithNotices) Error() string {
	// for SDK usecase:
	if en.injectNoticesToErrorText {
		var buf bytes.Buffer

		buf.WriteString(en.OverallMessage)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("Found %d notices:\n", len(en.Notices)))

		for _, notice := range en.Notices {
			buf.WriteString(fmt.Sprintf("- %s\n", notice.Message))
		}

		return buf.String()
	}

	// for CLI usecase:
	// notices will be printed to stdout
	return fmt.Sprintf("%s (has %d notices)", en.OverallMessage, len(en.Notices))
}
