package v2ray

import (
	"github.com/QQGoblin/go-sdk/pkg/tmpl"
	"github.com/lithammer/dedent"
	"io/ioutil"
	"text/template"
)

func WriteConfig(c *Config, variables map[string]interface{}) error {
	var (
		tmplStr []byte
		content string
	)

	tmplStr, err := ioutil.ReadFile(c.Tmpl)
	if err != nil {
		return err
	}

	configTmpl := template.Must(template.New("config.json").Parse(dedent.Dedent(string(tmplStr))))

	content, err = tmpl.Render(configTmpl, variables)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Path, []byte(content), 0644)

}
