package types

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

func (p *Prompt) Slugify() {
	a := time.Now().Unix()
	b := slug.Make(p.Title)
	if len(b) > 32 {
		b = b[:32]
	}
	p.Slug = fmt.Sprintf("%d_%s", a, b)
}
