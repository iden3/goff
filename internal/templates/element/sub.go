package element

const Sub = `
// Sub  z = x - y mod q
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Sub( x, y *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) Sub( x, y {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	var b uint64
        var xar, yar = x.GetUint64(), y.GetUint64()
	z[0], b = bits.Sub64(xar[0], yar[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		z[{{$i}}], b = bits.Sub64(xar[{{$i}}], yar[{{$i}}], b)
	{{- end}}
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], {{index $.Q 0}}, 0)
		{{- range $i := .NbWordsIndexesNoZero}}
			{{- if eq $i $.NbWordsLastIndex}}
				z[{{$i}}], _ = bits.Add64(z[{{$i}}], {{index $.Q $i}}, c)
			{{- else}}
				z[{{$i}}], c = bits.Add64(z[{{$i}}], {{index $.Q $i}}, c)
			{{- end}}
		{{- end}}
	}
	return z
}

// SubAssign  z = z - x mod q
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SubAssign(x *{{.ElementName}}) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) SubAssign(x {{.IfaceName}}) {{.IfaceName}} {
{{end}}
	var b uint64
        var xar = x.GetUint64()
	z[0], b = bits.Sub64(z[0], xar[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		z[{{$i}}], b = bits.Sub64(z[{{$i}}], xar[{{$i}}], b)
	{{- end}}
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], {{index $.Q 0}}, 0)
		{{- range $i := .NbWordsIndexesNoZero}}
			{{- if eq $i $.NbWordsLastIndex}}
				z[{{$i}}], _ = bits.Add64(z[{{$i}}], {{index $.Q $i}}, c)
			{{- else}}
				z[{{$i}}], c = bits.Add64(z[{{$i}}], {{index $.Q $i}}, c)
			{{- end}}
		{{- end}}
	}
	return z
}
`
