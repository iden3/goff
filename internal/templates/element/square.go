package element

const SquareAMD64 = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

//go:noescape
func square{{.ElementName}}(res,y *{{.ElementName}})


// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	square{{.ElementName}}(z, x)
	return z
}
{{else}}
func (z *{{.ElementName}}) Square(x {{.IfaceName}}) {{.IfaceName}} {
	square{{.ElementName}}(z, x.(* {{.ElementName}}))
	return z
}
{{end}}

`

const SquareCIOSNoCarry = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Square(x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
        var xar = x.GetUint64()
	{{if .NoCarrySquare}}
		{{ template "square" dict "all" . "V1" "xar"}}
		{{ template "reduce" . }}
		return z 
	{{else if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "xar" "V2" "xar"}}
		{{ template "reduce" . }}
		return z 
	{{else }}
		return z.Mul(x, x)
	{{end}}
}

{{- if eq .ASM false }}
func square{{.ElementName}}(z,x *{{.ElementName}}) {
	{{if .NoCarrySquare}}
		{{ template "square" dict "all" . "V1" "x"}}
		{{ template "reduce" . }}
	{{else if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x"}}
		{{ template "reduce" . }}
	{{else }}
		z.Mul(x, x)
	{{end}}
}
{{ end }}

`
