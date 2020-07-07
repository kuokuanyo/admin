package components

import (
	"bytes"
	"fmt"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"html/template"
	"strings"
)

// �����N�ŦXtemList(map[string]string)���ȥ[�Jtext(string)�A���۱N�ѼƤΥ\��K�[���s���ҪO�øѪR���ҪO���D��
// �N�Ѽ�compo�g�Jbuffer(bytes.Buffer)���̫��XHTML
func ComposeHtml(temList map[string]string, compo interface{}, templateName ...string) template.HTML {
	var text = ""
	// �Nmap[string]string�ŦXkey���ȥ[�Jtext
	for _, v := range templateName {
		text += temList["components/"+v]
	}

	// new�N���w���ѼƤ��t���s��HTML�ҪO
	// Funcs�K�[�s���\���ҪO
	// Parse�N�Ѽ�text�ѪR���ҪO���D��
	tmpl, err := template.New("comp").Funcs(template2.DefaultFuncMap).Parse(text)
	if err != nil {
		panic("ComposeHtml Error:" + err.Error())
	}
	buffer := new(bytes.Buffer)

	defineName := strings.Replace(templateName[0], "table/", "", -1)
	defineName = strings.Replace(defineName, "form/", "", -1)

	// �N�ĤT�ӰѼ�(compo)�g�Jbuffer��
	err = tmpl.ExecuteTemplate(buffer, defineName, compo)
	if err != nil {
		fmt.Println("ComposeHtml Error:", err)
	}

	// �Nbuffer��X��HTML
	return template.HTML(buffer.String())
}
