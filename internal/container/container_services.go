package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/formating"
)

func (c *Container) serviceNoticeFormatter() *formating.NoticeFormatter {
	return once(func() *formating.NoticeFormatter {
		return formating.NewNoticeFormatter(
			c.serviceRenderer(),
			c.serviceCodePrinter(),
		)
	})
}
