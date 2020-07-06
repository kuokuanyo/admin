package types

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type DisplayFnGenerator interface {
	Get(args ...interface{}) FieldFilterFn
	JS() template.HTML
	HTML() template.HTML
}

type BaseDisplayFnGenerator struct{}

func (base *BaseDisplayFnGenerator) JS() template.HTML   { return "" }
func (base *BaseDisplayFnGenerator) HTML() template.HTML { return "" }

var displayFnGens = make(map[string]DisplayFnGenerator)

func RegisterDisplayFnGenerator(key string, gen DisplayFnGenerator) {
	if _, ok := displayFnGens[key]; ok {
		panic("display function generator has been registered")
	}
	displayFnGens[key] = gen
}

type FieldDisplay struct {
	Display              FieldFilterFn
	DisplayProcessChains DisplayProcessFnChains
}

func (f FieldDisplay) ToDisplay(value FieldModel) interface{} {
	val := f.Display(value)

	if _, ok := val.(template.HTML); !ok {
		if _, ok2 := val.([]string); !ok2 {
			valStr := fmt.Sprintf("%v", val)
			for _, process := range f.DisplayProcessChains {
				valStr = process(valStr)
			}
			return valStr
		}
	}

	return val
}

// func(value string) stringFieldDisplay.DisplayProcessFnChains([]DisplayProcessFn)
// 硓筁把计limit耞func(value string)肚
func (f FieldDisplay) AddLimit(limit int) DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value string) string {
		if limit > len(value) {
			return value
		} else if limit < 0 {
			return ""
		} else {
			return value[:limit]
		}
	})
}

// func(value string) stringFieldDisplay.DisplayProcessFnChains([]DisplayProcessFn)
// func(value string)肚strings.TrimSpace(value)
func (f FieldDisplay) AddTrimSpace() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value string) string {
		return strings.TrimSpace(value)
	})
}

func (f FieldDisplay) AddSubstr(start int, end int) DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value string) string {
		if start > end || start > len(value) || end < 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(value) {
			end = len(value)
		}
		return value[start:end]
	})
}

// func(value string) stringFieldDisplay.DisplayProcessFnChains([]DisplayProcessFn)
func (f FieldDisplay) AddToTitle() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value string) string {
		return strings.Title(value)
	})
}

func (f FieldDisplay) AddToUpper() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value string) string {
		return strings.ToUpper(value)
	})
}

func (f FieldDisplay) AddToLower() DisplayProcessFnChains {
	return f.DisplayProcessChains.Add(func(value string) string {
		return strings.ToLower(value)
	})
}

type DisplayProcessFn func(string) string

type DisplayProcessFnChains []DisplayProcessFn

func (d DisplayProcessFnChains) Valid() bool {
	return len(d) > 0
}

// 盢把计f(func(string) string)globalDisplayProcessChains([]DisplayProcessFn)
func (d DisplayProcessFnChains) Add(f DisplayProcessFn) DisplayProcessFnChains {
	return append(d, f)
}

func (d DisplayProcessFnChains) Append(f DisplayProcessFnChains) DisplayProcessFnChains {
	return append(d, f...)
}

func (d DisplayProcessFnChains) Copy() DisplayProcessFnChains {
	if len(d) == 0 {
		return make(DisplayProcessFnChains, 0)
	} else {
		var newDisplayProcessFnChains = make(DisplayProcessFnChains, len(d))
		copy(newDisplayProcessFnChains, d)
		return newDisplayProcessFnChains
	}
}

func chooseDisplayProcessChains(internal DisplayProcessFnChains) DisplayProcessFnChains {
	if len(internal) > 0 {
		return internal
	}
	return globalDisplayProcessChains.Copy()
}

// globalDisplayProcessChains摸[]DisplayProcessFnDisplayProcessFn摸func(string) string
var globalDisplayProcessChains = make(DisplayProcessFnChains, 0)

// 盢把计f(func(string) string)globalDisplayProcessChains([]DisplayProcessFn)
func AddGlobalDisplayProcessFn(f DisplayProcessFn) {
	// type DisplayProcessFn func(string) string
	globalDisplayProcessChains = globalDisplayProcessChains.Add(f)
}

// func(value string) string把计globalDisplayProcessChains([]DisplayProcessFn)
// 硓筁把计limit耞func(value string)肚
func AddLimit(limit int) DisplayProcessFnChains {
	return addLimit(limit, globalDisplayProcessChains)
}

// func(value string) string把计globalDisplayProcessChains([]DisplayProcessFn)
// func(value string)肚strings.TrimSpace(value)
func AddTrimSpace() DisplayProcessFnChains {
	return addTrimSpace(globalDisplayProcessChains)
}

// func(value string) string把计globalDisplayProcessChains([]DisplayProcessFn)
// 硓筁把计startend耞func(value string)肚
func AddSubstr(start int, end int) DisplayProcessFnChains {
	return addSubstr(start, end, globalDisplayProcessChains)
}

// func(value string) stringglobalDisplayProcessChains([]DisplayProcessFn)
// func(value string)肚strings.Title(value)
func AddToTitle() DisplayProcessFnChains {
	return addToTitle(globalDisplayProcessChains)
}

// func(value string) stringglobalDisplayProcessChains([]DisplayProcessFn)
// func(value string)肚strings.ToUpper(value)
func AddToUpper() DisplayProcessFnChains {
	return addToUpper(globalDisplayProcessChains)
}

// func(value string) stringglobalDisplayProcessChains([]DisplayProcessFn)
// func(value string)肚strings.ToLower(value)
func AddToLower() DisplayProcessFnChains {
	return addToLower(globalDisplayProcessChains)
}

// func(value string) stringglobalDisplayProcessChains([]DisplayProcessFn)
// func(value string)肚html.EscapeString(value)
func AddXssFilter() DisplayProcessFnChains {
	return addXssFilter(globalDisplayProcessChains)
}

// func(value string) stringglobalDisplayProcessChains([]DisplayProcessFn)
// func(value string)肚replacer.Replace(value)
func AddXssJsFilter() DisplayProcessFnChains {
	return addXssJsFilter(globalDisplayProcessChains)
}

// func(value string) string把计chains([]DisplayProcessFn)
// 硓筁把计limit耞func(value string)肚
func addLimit(limit int, chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		if limit > len(value) {
			return value
		} else if limit < 0 {
			return ""
		} else {
			return value[:limit]
		}
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// func(value string)肚strings.TrimSpace(value)
func addTrimSpace(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		return strings.TrimSpace(value)
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// 硓筁把计startend耞func(value string)肚
func addSubstr(start int, end int, chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		if start > end || start > len(value) || end < 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(value) {
			end = len(value)
		}
		return value[start:end]
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// func(value string)肚strings.Title(value)
func addToTitle(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		return strings.Title(value)
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// func(value string)肚strings.ToUpper(value)
func addToUpper(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		return strings.ToUpper(value)
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// func(value string)肚strings.ToLower(value)
func addToLower(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		return strings.ToLower(value)
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// func(value string)肚html.EscapeString(value)
func addXssFilter(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		return html.EscapeString(value)
	})
	return chains
}

// func(value string) string把计chains([]DisplayProcessFn)
// func(value string)肚replacer.Replace(value)
func addXssJsFilter(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value string) string {
		replacer := strings.NewReplacer("<script>", "&lt;script&gt;", "</script>", "&lt;/script&gt;")
		return replacer.Replace(value)
	})
	return chains
}
