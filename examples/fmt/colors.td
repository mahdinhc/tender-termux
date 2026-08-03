// tender@v0.0.8
import "console"

cout := console.stdout()
cout.write("\x1b[0;32mhello world\x1b[0m\n")

cerr := console.stderr()
cerr.write("\x1b[0;31mtender script\x1b[0m\n")