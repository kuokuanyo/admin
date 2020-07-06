package captcha

type Captcha interface {
	Validate(token string) bool
}

var List = make(map[string]Captcha)

// 盢把计keycaptchaList(make(map[string]Captcha))
func Add(key string, captcha Captcha) {
	if _, exist := List[key]; exist {
		panic("captcha exist")
	}
	List[key] = captcha
}

// 耞List(make(map[string]Captcha))柑琌Τ把计key肚Captcha(interface)
func Get(key string) (Captcha, bool) {
	captcha, ok := List[key]
	return captcha, ok
}
