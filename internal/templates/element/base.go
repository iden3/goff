package element

const Base = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import (
	"math/big"
	"math/bits"
	"crypto/rand"
	"encoding/binary"
	"io"
	"sync"
	"unsafe"
	{{if eq .NoCollidingNames false}}"strconv"{{- end}}
)

// {{.ElementName}} represents a field element stored on {{.NbWords}} words (uint64)
// {{.ElementName}} are assumed to be in Montgomery form in all methods
// field modulus q =
// 
// {{.Modulus}}
type {{.ElementName}} [{{.NbWords}}]uint64

// {{.ElementName}}Limbs number of 64 bits words needed to represent {{.ElementName}}
const {{.ElementName}}Limbs = {{.NbWords}}

// {{.ElementName}}Bits number bits needed to represent {{.ElementName}}
const {{.ElementName}}Bits = {{.NbBits}}

// GetUint64 returns z[0],... z[N-1]
func (z {{.ElementName}}) GetUint64() []uint64  {
    return z[0:]
}

// SetUint64 z = v, sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetUint64(v uint64) *{{.IfaceName}} {
{{else}}
func (z *{{.ElementName}}) SetUint64(v uint64) {{.IfaceName}} {
{{end}}
	z[0] = v
	{{- range $i := .NbWordsIndexesNoZero}}
		z[{{$i}}] = 0
	{{- end}}
	return z.ToMont()
}

// Set z = x
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Set(x *{{.IfaceName}}) *{{.IfaceName}} {
{{else}}
func (z *{{.ElementName}}) Set(x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
        var xar = x.GetUint64()
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = xar[{{$i}}]
	{{- end}}
	return z
}

// Set z = x
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetFromArray(xar []uint64) *{{.IfaceName}} {
{{else}}
func (z *{{.ElementName}}) SetFromArray(xar []uint64) {{.IfaceName}} {
{{end}}

	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = xar[{{$i}}]
	{{- end}}
	return z
}

// SetZero z = 0
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetZero() *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) SetZero() {{.IfaceName}} {
{{end}}
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = 0
	{{- end}}
	return z
}

// SetOne z = 1 (in Montgomery form)
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetOne() *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) SetOne() {{.IfaceName}} {
{{end}}
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = {{index $.One $i}}
	{{- end}}
	return z
}


// Neg z = q - x 
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Neg( x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Neg( x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	if x.IsZero() {
		return z.SetZero()
	}
	var borrow uint64
        var xar = x.GetUint64()
	z[0], borrow = bits.Sub64({{index $.Q 0}}, xar[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
			z[{{$i}}], _ = bits.Sub64({{index $.Q $i}}, xar[{{$i}}], borrow)
		{{- else}}
			z[{{$i}}], borrow = bits.Sub64({{index $.Q $i}}, xar[{{$i}}], borrow)
		{{- end}}
	{{- end}}
	return z
}


// Div z = x*y^-1 mod q 
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Div( x, y *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Div( x, y {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	var yInv {{.ElementName}}
	yInv.Inverse( y)
	z.Mul( x, &yInv)
	return z
}

// Equal returns z == x
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Equal(x *{{.ElementName}}) bool {
{{else}}
func (z *{{.ElementName}}) Equal(x {{.IfaceName}}) bool {
{{end}}
        var xar = x.GetUint64()
	return {{- range $i :=  reverse .NbWordsIndexesNoZero}}(z[{{$i}}] == xar[{{$i}}]) &&{{end}}(z[0] == xar[0])
}

// IsZero returns z == 0
func (z *{{.ElementName}}) IsZero() bool {
	return ( {{- range $i :=  reverse .NbWordsIndexesNoZero}} z[{{$i}}] | {{end}}z[0]) == 0
}



// field modulus stored as big.Int 
var _{{.ElementName}}Modulus big.Int 
var once{{.ElementName}}Modulus sync.Once
func {{.ElementName}}Modulus() *big.Int {
	once{{.ElementName}}Modulus.Do(func() {
		_{{.ElementName}}Modulus.SetString("{{.Modulus}}", 10)
	})
	return &_{{.ElementName}}Modulus
}


{{/* We use big.Int for Inverse for these type of moduli */}}
{{if eq .NoCarry false}}

// Inverse z = x^-1 mod q 
// note: allocates a big.Int (math/big)
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Inverse( x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Inverse( x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	var _xNonMont big.Int
	x.ToBigIntRegular( &_xNonMont)
	_xNonMont.ModInverse(&_xNonMont, {{.ElementName}}Modulus())
	z.SetBigInt(&_xNonMont)
	return z
}

{{ else }}

// Inverse z = x^-1 mod q 
// Algorithm 16 in "Efficient Software-Implementation of Finite Fields with Applications to Cryptography"
// if x == 0, sets and returns z = x 
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Inverse(x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Inverse(x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	if x.IsZero() {
		return z.Set(x)
	}

	// initialize u = q
	var u = {{.ElementName}}{
		{{- range $i := .NbWordsIndexesFull}}
		{{index $.Q $i}},{{end}}
	}

	// initialize s = r^2
	var s = {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}

	// r = 0
	r := {{.ElementName}}{}

	v := x.GetUint64()

	var carry, borrow, t, t2 uint64
	var bigger, uIsOne, vIsOne bool

	for !uIsOne && !vIsOne {
		for v[0]&1 == 0 {
			{{ template "div2" dict "all" . "V" "v"}}
			if s[0]&1 == 1 {
				{{ template "add_q" dict "all" . "V1" "s" }}
			}
			{{ template "div2" dict "all" . "V" "s"}}
		} 
		for u[0]&1 == 0 {
			{{ template "div2" dict "all" . "V" "u"}}
			if r[0]&1 == 1 {
				{{ template "add_q" dict "all" . "V1" "r" }}
			}
			{{ template "div2" dict "all" . "V" "r"}}
		} 
		{{ template "bigger" dict "all" . "V1" "v" "V2" "u"}}
		if bigger  {
			{{ template "sub_noborrow" dict "all" . "V1" "v" "V2" "u" }}
			{{ template "bigger" dict "all" . "V1" "r" "V2" "s"}}
			if bigger {
				{{ template "add_q" dict "all" . "V1" "s" }}
			}
			{{ template "sub_noborrow" dict "all" . "V1" "s" "V2" "r" }}
		} else {
			{{ template "sub_noborrow" dict "all" . "V1" "u" "V2" "v" }}
			{{ template "bigger" dict "all" . "V1" "s" "V2" "r"}}
			if bigger {
				{{ template "add_q" dict "all" . "V1" "r" }}
			}
			{{ template "sub_noborrow" dict "all" . "V1" "r" "V2" "s" }}
		}
		uIsOne = (u[0] == 1) && ({{- range $i := reverse .NbWordsIndexesNoZero}}u[{{$i}}] {{if eq $i 1}}{{else}} | {{end}}{{end}} ) == 0
		vIsOne = (v[0] == 1) && ({{- range $i := reverse .NbWordsIndexesNoZero}}v[{{$i}}] {{if eq $i 1}}{{else}} | {{end}}{{end}} ) == 0
	}

	if uIsOne {
		z.Set(&r)
	} else {
		z.Set(&s)
	}

	return z
}

{{ end }}

// SetRandom sets z to a random element < q
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetRandom() *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) SetRandom() {{.IfaceName}} {
{{end}}
	bytes := make([]byte, {{mul 8 .NbWords}})
	io.ReadFull(rand.Reader, bytes)
	{{- range $i :=  .NbWordsIndexesFull}}
		{{- $k := add $i 1}}
		z[{{$i}}] = binary.BigEndian.Uint64(bytes[{{mul $i 8}}:{{mul $k 8}}]) 
	{{- end}}
	z[{{$.NbWordsLastIndex}}] %= {{index $.Q $.NbWordsLastIndex}}

	{{ template "reduce" . }}

	return z
}

{{ if .NoCollidingNames}}
{{ else}}

// One returns 1 (in montgommery form)
{{- if eq .IfaceName .ElementName}} 
func (z {{.ElementName}}) One()  *{{.ElementName}} {
{{else}}
func (z {{.ElementName}}) One()  {{.IfaceName}} {
{{end}}
	one := z
	one.SetOne()
	return &one
}

{{end}}


{{ define "bigger" }}
	// {{$.V1}} >= {{$.V2}}
	bigger = !({{- range $i := reverse $.all.NbWordsIndexesNoZero}} {{$.V1}}[{{$i}}] < {{$.V2}}[{{$i}}] || ( {{$.V1}}[{{$i}}] == {{$.V2}}[{{$i}}] && (
		{{- end}}{{$.V1}}[0] < {{$.V2}}[0] {{- range $i :=  $.all.NbWordsIndexesNoZero}} )) {{- end}} )
{{ end }}

{{ define "add_q" }}
	// {{$.V1}} = {{$.V1}} + q 
	{{$.V1}}[0], carry = bits.Add64({{$.V1}}[0], {{index $.all.Q 0}}, 0)
	{{- range $i := .all.NbWordsIndexesNoZero}}
		{{- if eq $i $.all.NbWordsLastIndex}}
			{{$.V1}}[{{$i}}], _ = bits.Add64({{$.V1}}[{{$i}}], {{index $.all.Q $i}}, carry)
		{{- else}}
			{{$.V1}}[{{$i}}], carry = bits.Add64({{$.V1}}[{{$i}}], {{index $.all.Q $i}}, carry)
		{{- end}}
	{{- end}}
{{ end }}

{{ define "sub_noborrow" }}
	// {{$.V1}} = {{$.V1}} - {{$.V2}}
	{{$.V1}}[0], borrow = bits.Sub64({{$.V1}}[0], {{$.V2}}[0], 0)
	{{- range $i := .all.NbWordsIndexesNoZero}}
		{{- if eq $i $.all.NbWordsLastIndex}}
			{{$.V1}}[{{$i}}], _ = bits.Sub64({{$.V1}}[{{$i}}], {{$.V2}}[{{$i}}], borrow)
		{{- else}}
			{{$.V1}}[{{$i}}], borrow = bits.Sub64({{$.V1}}[{{$i}}], {{$.V2}}[{{$i}}], borrow)
		{{- end}}
	{{- end}}
{{ end }}


{{ define "div2" }}
	// {{$.V}} = {{$.V}} >> 1
	{{- range $i :=  reverse .all.NbWordsIndexesNoZero}}
		{{- if eq $i $.all.NbWordsLastIndex}}
			t2 = {{$.V}}[{{$i}}] << 63
			{{$.V}}[{{$i}}] >>= 1
		{{- else}}
			t2 = {{$.V}}[{{$i}}] << 63
			{{$.V}}[{{$i}}] = ({{$.V}}[{{$i}}] >> 1) | t
		{{- end}}
		t = t2
	{{- end}}
	{{$.V}}[0] = ({{$.V}}[0] >> 1) | t
{{ end }}


`
