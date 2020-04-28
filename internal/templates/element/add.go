package element

const Add = `
// Add z = x + y mod q
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Add( x, y *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Add( x, y {{.IfaceName}}) {{.IfaceName}} {
{{end}}

	var carry uint64
        var xar, yar = x.GetUint64(), y.GetUint64()
	{{$k := sub $.NbWords 1}}
	z[0], carry = bits.Add64(xar[0], yar[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
		{{- else}}
			z[{{$i}}], carry = bits.Add64(xar[{{$i}}], yar[{{$i}}], carry)
		{{- end}}
	{{- end}}
	{{- if .NoCarry}}
		z[{{$k}}], _ = bits.Add64(xar[{{$k}}], yar[{{$k}}], carry)
	{{- else }}
		z[{{$k}}], carry = bits.Add64(xar[{{$k}}], yar[{{$k}}], carry)
		// if we overflowed the last addition, z >= q
		// if z >= q, z = z - q
		if carry != 0 {
			// we overflowed, so z >= q
			z[0], carry = bits.Sub64(z[0], {{index $.Q 0}}, 0)
			{{- range $i := .NbWordsIndexesNoZero}}
				z[{{$i}}], carry = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, carry)
			{{- end}}
			return z
		}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}

// AddAssign z = z + x mod q
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) AddAssign(x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) AddAssign(x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	var carry uint64
        var xar = x.GetUint64()
	{{$k := sub $.NbWords 1}}
	z[0], carry = bits.Add64(z[0], xar[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
		{{- else}}
			z[{{$i}}], carry = bits.Add64(z[{{$i}}], xar[{{$i}}], carry)
		{{- end}}
	{{- end}}
	{{- if .NoCarry}}
		z[{{$k}}], _ = bits.Add64(z[{{$k}}], xar[{{$k}}], carry)
	{{- else }}
		z[{{$k}}], carry = bits.Add64(z[{{$k}}], xar[{{$k}}], carry)
		// if we overflowed the last addition, z >= q
		// if z >= q, z = z - q
		if carry != 0 {
			// we overflowed, so z >= q
			z[0], carry = bits.Sub64(z[0], {{index $.Q 0}}, 0)
			{{- range $i := .NbWordsIndexesNoZero}}
				z[{{$i}}], carry = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, carry)
			{{- end}}
			return z
		}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}

// Double z = x + x mod q, aka Lsh 1
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Double( x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Double( x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	var carry uint64
        var xar = x.GetUint64()
	{{$k := sub $.NbWords 1}}
	z[0], carry = bits.Add64(xar[0], xar[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
		{{- else}}
			z[{{$i}}], carry = bits.Add64(xar[{{$i}}], xar[{{$i}}], carry)
		{{- end}}
	{{- end}}
	{{- if .NoCarry}}
		z[{{$k}}], _ = bits.Add64(xar[{{$k}}], xar[{{$k}}], carry)
	{{- else }}
		z[{{$k}}], carry = bits.Add64(xar[{{$k}}], xar[{{$k}}], carry)
		// if we overflowed the last addition, z >= q
		// if z >= q, z = z - q
		if carry != 0 {
			// we overflowed, so z >= q
			z[0], carry = bits.Sub64(z[0], {{index $.Q 0}}, 0)
			{{- range $i := .NbWordsIndexesNoZero}}
				z[{{$i}}], carry = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, carry)
			{{- end}}
			return z
		}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}
`
