import "crypto"

var str = "Hello World"

md5 := crypto.md5(str)
sha1 := crypto.sha1(str)
sha256 := crypto.sha256(str)
sha512 := crypto.sha512(str)
sha3_256 := crypto.sha3_256(str)

println(md5)
println(sha1)
println(sha256)
println(sha512)
println(sha3_256)

