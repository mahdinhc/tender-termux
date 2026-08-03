import "image"
import "helper/imgview"

img := image.load("gopher.png")

img.set_blue(img.get_blue()*0)
img.set_green(img.get_green()*0)

img.encode("png") |> imgview