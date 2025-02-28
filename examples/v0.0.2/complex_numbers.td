import "cmplx"

real := 1
imagine := 5

c := complex(real, imagine)

b := cmplx.new(3, 1.5)

a := b + c + 100


println(c)
println(b)
println()
println(a)

println("real part:", a.real)
println("imagine part:", a.imag)
