// tender@v0.0.8
import "console"

header := ["Age", "Name", "Birth"]
body := [
	["29", "Chloe", "2000"],
	["21", "Kiwi", "2000".green],
	["22", "Crow", "2000"]
]


table := console.table(header, body)

println(table)