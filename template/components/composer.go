package components

import (
	"bytes"
	"fmt"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"html/template"
	"strings"
)

// 盢才temList(map[string]string)text(string)钡帝盢把计の睰倒穝家狾秆猂家狾砰
// 盢把计compo糶buffer(bytes.Buffer)い程块HTML
func ComposeHtml(temList map[string]string, compo interface{}, templateName ...string) template.HTML {
	var text = ""
	// 盢map[string]string才keytext
	for _, v := range templateName {
		text += temList["components/"+v]
	}

	// new盢倒﹚把计だ皌倒穝HTML家狾
	// Funcs睰穝家狾
	// Parse盢把计text秆猂家狾砰
	tmpl, err := template.New("comp").Funcs(template2.DefaultFuncMap).Parse(text)
	if err != nil {
		panic("ComposeHtml Error:" + err.Error())
	}
	buffer := new(bytes.Buffer)

	defineName := strings.Replace(templateName[0], "table/", "", -1)
	defineName = strings.Replace(defineName, "form/", "", -1)

	// 盢材把计(compo)糶bufferい
	err = tmpl.ExecuteTemplate(buffer, defineName, compo)
	if err != nil {
		fmt.Println("ComposeHtml Error:", err)
	}

	// 盢buffer块ΘHTML
	return template.HTML(buffer.String())
}
