package element

const Exp = `
// Exp z = x^exponent mod q
// (not optimized)
// exponent (non-montgomery form) is ordered from least significant word to most significant word
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Exp(x {{.ElementName}}, exponent ...uint64) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Exp(x {{.IfaceName}}, exponent ...uint64) {{.IfaceName}} {
{{end}}
	r := 0
	msb := 0
	for i := len(exponent) - 1; i>=0; i-- {
		if exponent[i] == 0 {
			r++
		} else {
			msb = (i * 64) + bits.Len64(exponent[i])
			break
		}
	} 
	exponent = exponent[:len(exponent)-r]
	if len(exponent) == 0 {
		return z.SetOne()
	}

        {{- if eq .IfaceName .ElementName}} 
	z.Set(&x)
        {{else}}
	z.Set(x)
        {{end}}

	l := msb - 2
	for i := l; i >= 0; i-- {
		z.Square(z)
		if exponent[i / 64]&(1<<uint(i%64)) != 0 {
                        {{- if eq .IfaceName .ElementName}} 
			z.MulAssign(&x)
                        {{else}}
			z.MulAssign(x)
                        {{end}}
 
		}
	}
	return z
}

`
