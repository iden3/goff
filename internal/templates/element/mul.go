package element

const MontgomeryMultiplication = `
// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Mul(x, y {{.IfaceName}}) {{.IfaceName}} {
{{end}}
        var xar, yar = x.GetUint64(), y.GetUint64()
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "xar" "V2" "yar"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "xar" "V2" "yar" "NoReturn" false}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

// MulAssign z = z * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) MulAssign(x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
        var xar = x.GetUint64()
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "xar"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "xar" "NoReturn" false}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

{{- if eq .ASM false }}
{{- if eq .IfaceName .ElementName}} 
func mulAssign{{.ElementName}}(z,x *{{.ElementName}}) {
{{else}}
func MulAssign{{.ElementName}}(z,x {{.IfaceName}}) {
{{end}}
        var xar = x.GetUint64()
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "xar"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "xar" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}

func fromMont{{.ElementName}}(z *{{.ElementName}}) {
	// the following lines implement z = z * 1
	// with a modified CIOS montgomery multiplication
	{{- range $j := .NbWordsIndexesFull}}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * {{index $.QInverse 0}}
		C := madd0(m, {{index $.Q 0}}, z[0])
		{{- range $i := $.NbWordsIndexesNoZero}}
			C, z[{{sub $i 1}}] = madd2(m, {{index $.Q $i}}, z[{{$i}}], C)
		{{- end}}
		z[{{sub $.NbWords 1}}] = C
	}
	{{- end}}

	{{ template "reduce" .}}
}

// for test purposes
func reduce{{.ElementName}}(z *{{.ElementName}})  {
	{{ template "reduce" . }}
}
{{- end}}


`
