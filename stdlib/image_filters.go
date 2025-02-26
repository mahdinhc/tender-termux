package stdlib

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/2dprototype/tender"
)

// makeImageFilters wraps an image.Image and returns an ImmutableMap of filter functions.
func makeImageFilters(img image.Image) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"blur": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					r, _ := tender.ToInt(args[0])
					newImg := applyBlurParallelOptimized(img, r)
					return makeImage(newImg), nil
				},
			},
			"bnw": &tender.UserFunction{
				Name: "bnw",
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					t, _ := tender.ToUint8(args[0])
					return makeImage(applyBnWParallel(img, t)), nil
				},
			},
			// "glitch" filter: optionally accepts a maxShift (default 20)
			"glitch": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					m, _ := tender.ToInt(args[0])
					return makeImage(applyGlitchParallel(img, m)), nil
				},
			},
			"invert": &tender.UserFunction{
				Name: "invert",
				Value: func(args ...tender.Object) (tender.Object, error) {
					newImg := applyInvertParallel(img)
					return makeImage(newImg), nil
				},
			},
			"grayscale": &tender.UserFunction{
				Name: "grayscale",
				Value: func(args ...tender.Object) (tender.Object, error) {
					newImg := applyGrayscaleParallel(img)
					return makeImage(newImg), nil
				},
			},
			"sepia": &tender.UserFunction{
				Name: "sepia",
				Value: func(args ...tender.Object) (tender.Object, error) {
					newImg := applySepiaParallel(img)
					return makeImage(newImg), nil
				},
			},
			"brightness": &tender.UserFunction{
				Name: "brightness",
				Value: func(args ...tender.Object) (tender.Object, error) {
					// Required: offset (can be negative)
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					offset, ok := tender.ToInt(args[0])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "offset",
							Expected: "int",
							Found:    args[0].TypeName(),
						}
					}
					newImg := applyBrightnessParallel(img, offset)
					return makeImage(newImg), nil
				},
			},
			"contrast": &tender.UserFunction{
				Name: "contrast",
				Value: func(args ...tender.Object) (tender.Object, error) {
					// Required: factor (float, e.g., 1.0 = no change)
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					factor, ok := tender.ToFloat64(args[0])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "factor",
							Expected: "float",
							Found:    args[0].TypeName(),
						}
					}
					newImg := applyContrastParallel(img, factor)
					return makeImage(newImg), nil
				},
			},
			"saturation": &tender.UserFunction{
				Name: "saturation",
				Value: func(args ...tender.Object) (tender.Object, error) {
					// Required: factor (float; 0 = grayscale, 1 = original)
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					factor, ok := tender.ToFloat64(args[0])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     "factor",
							Expected: "float",
							Found:    args[0].TypeName(),
						}
					}
					newImg := applySaturationParallel(img, factor)
					return makeImage(newImg), nil
				},
			},
			"sharpen": &tender.UserFunction{
				Name: "sharpen",
				Value: func(args ...tender.Object) (tender.Object, error) {
					newImg := applyConvolutionParallel(img, [][]float64{
						{0, -1, 0},
						{-1, 5, -1},
						{0, -1, 0},
					}, 1, 0)
					return makeImage(newImg), nil
				},
			},
			"emboss": &tender.UserFunction{
				Name: "emboss",
				Value: func(args ...tender.Object) (tender.Object, error) {
					// Using an emboss kernel with an offset to recenter colors.
					newImg := applyConvolutionParallel(img, [][]float64{
						{-2, -1, 0},
						{-1, 1, 1},
						{0, 1, 2},
					}, 1, 128)
					return makeImage(newImg), nil
				},
			},
			"edge": &tender.UserFunction{
				Name: "edge",
				Value: func(args ...tender.Object) (tender.Object, error) {
					// Edge detection kernel; offset added for visibility.
					newImg := applyConvolutionParallel(img, [][]float64{
						{1, 1, 1},
						{1, -8, 1},
						{1, 1, 1},
					}, 1, 128)
					return makeImage(newImg), nil
				},
			},
		},
	}
}

// -------------------
// Filter Implementations
// -------------------

// applyBlurParallelOptimized applies a blur filter using a sliding window approach in parallel.
func applyBlurParallelOptimized(img image.Image, radius int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	blurImg := image.NewRGBA(bounds)
	pixels := blurImg.Pix

	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := height / numWorkers

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := i * rowHeight
		endY := (i + 1) * rowHeight
		if i == numWorkers-1 {
			endY = height
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := 0; x < width; x++ {
					var rSum, gSum, bSum, count float64
					for ky := -radius; ky <= radius; ky++ {
						for kx := -radius; kx <= radius; kx++ {
							nx, ny := x+kx, y+ky
							if nx >= 0 && nx < width && ny >= 0 && ny < height {
								r, g, b, _ := img.At(nx+bounds.Min.X, ny+bounds.Min.Y).RGBA()
								rSum += float64(r >> 8)
								gSum += float64(g >> 8)
								bSum += float64(b >> 8)
								count++
							}
						}
					}
					offset := (y*width + x) * 4
					pixels[offset] = uint8(clamp(rSum/count, 0, 255))
					pixels[offset+1] = uint8(clamp(gSum/count, 0, 255))
					pixels[offset+2] = uint8(clamp(bSum/count, 0, 255))
					pixels[offset+3] = 255
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return blurImg
}

// applyBnWParallel converts the image to black and white using a luminance threshold.
func applyBnWParallel(img image.Image, threshold uint8) image.Image {
	bounds := img.Bounds()
	bnwImg := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					lum := 0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8)
					var val uint8
					if lum > float64(threshold) {
						val = 255
					}
					bnwImg.Set(x, y, color.RGBA{val, val, val, uint8(a >> 8)})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return bnwImg
}

// applyGlitchParallel applies a glitch effect by shifting RGB channels randomly.
func applyGlitchParallel(img image.Image, maxShift int) image.Image {
	bounds := img.Bounds()
	glitchImg := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					rx := x + rand.Intn(maxShift) - maxShift/2
					gx := x
					bx := x - rand.Intn(maxShift) + maxShift/2
					if rx < bounds.Min.X || rx >= bounds.Max.X {
						rx = x
					}
					if gx < bounds.Min.X || gx >= bounds.Max.X {
						gx = x
					}
					if bx < bounds.Min.X || bx >= bounds.Max.X {
						bx = x
					}
					r, _, _, _ := img.At(rx, y).RGBA()
					_, g, _, _ := img.At(gx, y).RGBA()
					_, _, b, a := img.At(bx, y).RGBA()
					glitchImg.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return glitchImg
}

// applyInvertParallel inverts the image colors.
func applyInvertParallel(img image.Image) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					orig := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					inv := color.RGBA{
						R: 255 - orig.R,
						G: 255 - orig.G,
						B: 255 - orig.B,
						A: orig.A,
					}
					out.Set(x, y, inv)
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// applyGrayscaleParallel converts the image to grayscale.
func applyGrayscaleParallel(img image.Image) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					orig := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					// Using standard luminance formula.
					lum := uint8(0.299*float64(orig.R) + 0.587*float64(orig.G) + 0.114*float64(orig.B))
					out.Set(x, y, color.RGBA{R: lum, G: lum, B: lum, A: orig.A})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// applySepiaParallel applies a sepia tone to the image.
func applySepiaParallel(img image.Image) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					orig := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					r := clamp(0.393*float64(orig.R)+0.769*float64(orig.G)+0.189*float64(orig.B), 0, 255)
					g := clamp(0.349*float64(orig.R)+0.686*float64(orig.G)+0.168*float64(orig.B), 0, 255)
					b := clamp(0.272*float64(orig.R)+0.534*float64(orig.G)+0.131*float64(orig.B), 0, 255)
					out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), orig.A})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// applyBrightnessParallel adjusts brightness by adding an offset to each color channel.
func applyBrightnessParallel(img image.Image, offset int) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					orig := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					r := clamp(float64(orig.R+uint8(offset)), 0, 255)
					g := clamp(float64(orig.G+uint8(offset)), 0, 255)
					b := clamp(float64(orig.B+uint8(offset)), 0, 255)
					out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), orig.A})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// applyContrastParallel adjusts contrast by a factor.
func applyContrastParallel(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					orig := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					// Adjust contrast: new = ((old - 128) * factor) + 128.
					r := clamp(((float64(orig.R)-128)*factor)+128, 0, 255)
					g := clamp(((float64(orig.G)-128)*factor)+128, 0, 255)
					b := clamp(((float64(orig.B)-128)*factor)+128, 0, 255)
					out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), orig.A})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// applySaturationParallel adjusts saturation by interpolating between grayscale and original color.
func applySaturationParallel(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := bounds.Dy() / numWorkers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*rowHeight
		endY := bounds.Min.Y + (i+1)*rowHeight
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					orig := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					// Calculate luminance.
					lum := 0.299*float64(orig.R) + 0.587*float64(orig.G) + 0.114*float64(orig.B)
					// Interpolate between gray and original.
					r := clamp(lum+factor*(float64(orig.R)-lum), 0, 255)
					g := clamp(lum+factor*(float64(orig.G)-lum), 0, 255)
					b := clamp(lum+factor*(float64(orig.B)-lum), 0, 255)
					out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), orig.A})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// applyConvolutionParallel applies a convolution filter using the provided kernel, divisor, and offset.
func applyConvolutionParallel(img image.Image, kernel [][]float64, divisor, offset float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	out := image.NewRGBA(bounds)
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	rowHeight := height / numWorkers
	kSize := len(kernel)
	kRadius := kSize / 2

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startY := i * rowHeight
		endY := (i + 1) * rowHeight
		if i == numWorkers-1 {
			endY = height
		}
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := 0; x < width; x++ {
					var rAcc, gAcc, bAcc float64
					for ky := 0; ky < kSize; ky++ {
						for kx := 0; kx < kSize; kx++ {
							ix := x + kx - kRadius
							iy := y + ky - kRadius
							// Clamp to bounds.
							if ix < 0 {
								ix = 0
							}
							if iy < 0 {
								iy = 0
							}
							if ix >= width {
								ix = width - 1
							}
							if iy >= height {
								iy = height - 1
							}
							r, g, b, _ := img.At(ix+bounds.Min.X, iy+bounds.Min.Y).RGBA()
							weight := kernel[ky][kx]
							rAcc += float64(r>>8) * weight
							gAcc += float64(g>>8) * weight
							bAcc += float64(b>>8) * weight
						}
					}
					rVal := clamp((rAcc/divisor)+offset, 0, 255)
					gVal := clamp((gAcc/divisor)+offset, 0, 255)
					bVal := clamp((bAcc/divisor)+offset, 0, 255)
					out.Set(x+bounds.Min.X, y+bounds.Min.Y, color.RGBA{uint8(rVal), uint8(gVal), uint8(bVal), 255})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return out
}

// clamp ensures value is within min and max.
func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}
