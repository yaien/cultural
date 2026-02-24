package models

import (
	"fmt"
	"html/template"
)

func (c *pageData) Script() (template.HTML, error) {
	script := fmt.Sprintf("<script type=%q>\n(async () => {\n%s\n})();\n</script>", "text/javascript", c.Page.Script)
	return template.HTML(script), nil
}
