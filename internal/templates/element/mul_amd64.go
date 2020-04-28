package element

const MontgomeryMultiplicationAMD64 = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

//go:noescape
func mulAssign{{.ElementName}}(res,y *{{.ElementName}})

//go:noescape
func fromMont{{.ElementName}}(res *{{.ElementName}}) 

//go:noescape
func reduce{{.ElementName}}(res *{{.ElementName}})  // for test purposes

// Mul z = x * y mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	if z == x {
		mulAssign{{.ElementName}}(z, y)
		return z
	} else if z == y {
		mulAssign{{.ElementName}}(z, x)
		return z
	} else {
		z.Set(x)
		mulAssign{{.ElementName}}(z, y)
		return z
	}
}
{{else}}
func (z *{{.ElementName}}) Mul(x, y {{.IfaceName}}) {{.IfaceName}} {
	if z == x {
		mulAssign{{.ElementName}}(z, y.(* {{.ElementName}}))
		return z
	} else if z == y {
		mulAssign{{.ElementName}}(z, x.(* {{.ElementName}}))
		return z
	} else {
		z.Set(x)
		mulAssign{{.ElementName}}(z, y.(* {{.ElementName}}))
		return z
	}
}
{{end}}


// MulAssign z = z * x mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
{{- if eq .IfaceName .ElementName}} 
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	mulAssign{{.ElementName}}(z, x)
	return z 
}
{{else}}
func (z *{{.ElementName}}) MulAssign(x {{.IfaceName}}) {{.IfaceName}} {
	mulAssign{{.ElementName}}(z, x.(* {{.ElementName}}))
	return z 
}
{{end}}
`
