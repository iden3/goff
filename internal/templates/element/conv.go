package element

// note: not thourougly tested on moduli != .NoCarry
const FromMont = `
// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) FromMont() {{.IfaceName}} {
{{end}}
	fromMont{{.ElementName}}(z)
	return z 
}
`

const Conv = `
// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) ToMont() *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) ToMont() {{.IfaceName}} {
{{end}}

	var rSquare = {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	mulAssign{{.ElementName}}(z, &rSquare)
	return z
}

// ToRegular returns z in regular form (doesn't mutate z)
{{- if eq .IfaceName .ElementName}} 
func (z {{.ElementName}}) ToRegular() {{.ElementName}} {
	return *z.FromMont()
{{else}}
func (z {{.ElementName}}) ToRegular() {{.IfaceName}} {
	return z.FromMont()
{{end}}
}

// String returns the string form of an {{.ElementName}} in Montgomery form
func (z *{{.ElementName}}) String() string {
	var _z big.Int
	return z.ToBigIntRegular(&_z).String()
}

// ToByte returns the byte form of an {{.ElementName}} in Regular form
func (z {{.ElementName}}) ToByte() []byte {
{{- if eq .IfaceName .ElementName}} 
	t := z.ToRegular()
{{else}}
	t := z.ToRegular().(*{{.ElementName}})
{{end}}
	var _z []byte
	_z1 := make([]byte,8)
	{{- range $i := .NbWordsIndexesFull}}
		binary.LittleEndian.PutUint64(_z1, t[{{$i}}])
                _z = append(_z,_z1...)
	{{- end}}
	return _z
}

// FromByte returns the byte form of an {{.ElementName}} in Regular form (mutates z)
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) FromByte(x []byte) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) FromByte(x []byte) {{.IfaceName}} {
{{end}}
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = binary.LittleEndian.Uint64(x[{{$i}}*8:({{$i}}+1)*8])
	{{- end}}
        return z.ToMont()
}


// ToBigInt returns z as a big.Int in Montgomery form 
func (z *{{.ElementName}}) ToBigInt(res *big.Int) *big.Int {
      if bits.UintSize == 64 {
	bits := (*[{{.NbWords}}]big.Word)(unsafe.Pointer(z))
	return res.SetBits(bits[:])
      } else {
        var bits[{{.NbWords}}*2]big.Word
	{{- range $i := .NbWordsIndexesFull}}
            bits[{{$i}}*2] = big.Word(z[{{$i}}])
            bits[{{$i}}*2+1] = big.Word(z[{{$i}}] >> 32)
        {{- end }}
	return res.SetBits(bits[:])
      }
}

// ToBigIntRegular returns z as a big.Int in regular form 
func (z {{.ElementName}}) ToBigIntRegular(res *big.Int) *big.Int {
      if bits.UintSize == 64 {
	z.FromMont()
	bits := (*[{{.NbWords}}]big.Word)(unsafe.Pointer(&z))
	return res.SetBits(bits[:])
      } else {
        var bits[{{.NbWords}}*2]big.Word
	{{- range $i := .NbWordsIndexesFull}}
            bits[{{$i}}*2] = big.Word(z[{{$i}}])
            bits[{{$i}}*2+1] = big.Word(z[{{$i}}] >> 32)
        {{- end }}
	return res.SetBits(bits[:])
      }
}

// SetBigInt sets z to v (regular form) and returns z in Montgomery form
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetBigInt(v *big.Int) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) SetBigInt(v *big.Int) {{.IfaceName}} {
{{end}}
	z.SetZero()

	zero := big.NewInt(0)
	q := {{.ElementName}}Modulus()

	// fast path
	c := v.Cmp(q)
	if c == 0 {
		return z
	} else if c != 1 && v.Cmp(zero) != -1 {
		// v should
		vBits := v.Bits()
		for i := 0; i < len(vBits); i++ {
			z[i] = uint64(vBits[i])
		}
		return z.ToMont()
	}
	
	// copy input
	vv := new(big.Int).Set(v)
	vv.Mod(v, q)
	
	// v should
	vBits := vv.Bits()
        if bits.UintSize == 64 {
	   for i := 0; i < len(vBits); i++ {
		z[i] = uint64(vBits[i])
	   }
        } else {
	   for i := 0; i < len(vBits); i++ { 
              if i%2 == 0 {
                   z[i/2] = uint64(vBits[i])
              } else {
                   z[i/2] |= uint64(vBits[i]) << 32
              }
           }
        }
	return z.ToMont()
}

// SetString creates a big.Int with s (in base 10) and calls SetBigInt on z
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) SetString( s string) *{{.ElementName}} {
{{else}}
func (z *{{.ElementName}}) SetString( s string) {{.IfaceName}} {
{{end}}
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("{{.ElementName}}.SetString failed -> can't parse number in base10 into a big.Int")
	}
	return z.SetBigInt(x)
}

`
