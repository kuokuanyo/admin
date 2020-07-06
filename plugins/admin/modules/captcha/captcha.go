package captcha

type Captcha interface {
	Validate(token string) bool
}

var List = make(map[string]Captcha)

// �N�Ѽ�key�Bcaptcha�[�JList(make(map[string]Captcha))
func Add(key string, captcha Captcha) {
	if _, exist := List[key]; exist {
		panic("captcha exist")
	}
	List[key] = captcha
}

// �P�_List(make(map[string]Captcha))�̬O�_���Ѽ�key���Ȩæ^��Captcha(interface)
func Get(key string) (Captcha, bool) {
	captcha, ok := List[key]
	return captcha, ok
}
