package sqlformatter

import (
	"net/url"
)

func parseQuery() int {
	u := url.Parse("https://example.org/?a=1&a=2&b=&=3&&&&")
	u.Query()
}
