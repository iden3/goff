package element

const Reduce = `
{{ define "reduce" }}
// if z > q --> z -= q
// note: this is NOT constant time
if !({{- range $i := reverse .NbWordsIndexesNoZero}} z[{{$i}}] < {{index $.Q $i}} || ( z[{{$i}}] == {{index $.Q $i}} && (
{{- end}}z[0] < {{index $.Q 0}} {{- range $i :=  .NbWordsIndexesNoZero}} )) {{- end}} ){
	var b uint64
	z[0], b = bits.Sub64(z[0], {{index $.Q 0}}, 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{-  if eq $i $.NbWordsLastIndex}}
			z[{{$i}}], _ = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, b)
		{{-  else  }}
			z[{{$i}}], b = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, b)
		{{- end}}
	{{- end}}
}
{{-  end }}

`

const Reduce2 = `
{{ define "reduce2" }}
// if z > q --> z -= q
// note: this is NOT constant time
if !({{- range $i := reverse .NbWordsIndexesNoZero}} zar[{{$i}}] < {{index $.Q $i}} || ( zar[{{$i}}] == {{index $.Q $i}} && (
{{- end}}zar[0] < {{index $.Q 0}} {{- range $i :=  .NbWordsIndexesNoZero}} )) {{- end}} ){
	var b uint64
	zar[0], b = bits.Sub64(zar[0], {{index $.Q 0}}, 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{-  if eq $i $.NbWordsLastIndex}}
			zar[{{$i}}], _ = bits.Sub64(zar[{{$i}}], {{index $.Q $i}}, b)
		{{-  else  }}
			zar[{{$i}}], b = bits.Sub64(zar[{{$i}}], {{index $.Q $i}}, b)
		{{- end}}
	{{- end}}
}
        z.SetFromArray(zar)
{{-  end }}

`

