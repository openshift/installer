package ovirt

import (
	"fmt"
	"sort"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func askTemplate(c *ovirtsdk4.Connection, p *ovirt.Platform) error {
	var templateName string
	var templateByNames = make(map[string]*ovirtsdk4.Template)
	var templateNames []string
	systemService := c.SystemService().TemplatesService()
	templateResp, err := systemService.List().Send()
	if err != nil {
		return err
	}
	templates, ok := templateResp.Templates()
	if !ok {
		return fmt.Errorf("there are no available templates")
	}

	for _, template := range templates.Slice() {
		templateByNames[template.MustName()] = template
		templateNames = append(templateNames, template.MustName())
	}
	if err := survey.AskOne(&survey.Select{
		Message: "Template",
		Help:    "The oVirt template for the deployed VMs. 'Blank' is the default template.",
		Options: templateNames,
	},
		&templateName,
		func(ans interface{}) error {
			choice := ans.(string)
			sort.Strings(templateNames)
			i := sort.SearchStrings(templateNames, choice)
			if i == len(templateNames) || templateNames[i] != choice {
				return fmt.Errorf("invalid template %s", choice)
			}
			template, ok := templateByNames[choice]
			if !ok {
				return fmt.Errorf("cannot find a template by name %s", choice)
			}
			p.TemplateName = template.MustName()
			return nil
		}); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	return nil
}
