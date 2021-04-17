package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"sync"
)

//"github.com/muesli/gamut"
//Very generic check function to reduce boilerplate.
func check(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

//A helper function to confirm that a folder does indeed, exist. Borrowed from
//https://stackoverflow.com/questions/34458625/run-multiple-exec-commands-in-the-same-shell-golang
func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//These values describe our template pixels.
const (
	_ = iota
	Transparent
	Bit
	Outline
	Fill
	Accent
)

//This feels somewhat wasteful, though at the moment we're only calling it on bool flags.  Borrowed from
//https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go/35809400
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

//In my mind, I think it might be better if we initialize our flags in a single map.
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var templateString = flag.String("template", "", "choose template to render")
var foldPref = flag.String("fold", "None", "sets fold preference for template (Even/Odd/None)")
var outlinePref = flag.Bool("outline", true, "sets outline preference")
var rgbPref = flag.String("rgb", "255:255:255", "renders all images with the RGB provided ('r:g:b', where r,g,b are values 0-255)")
var gradientPref = flag.String("gradient", "", "renders all images based on gradient chosen") //TODO
var compositePref = flag.Int("sheetWidth", 8, "sets width of output sprite sheet, must return a whole number for 256/compositeWidth")

//var color = gamut.Hex("#333")

func main() {
	//Profiling
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//Trace
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	//Now the program can begin
	fmt.Print("BitSprite: Making 256 versions of 1 thing since 2020")

	//Grab the flag values
	templateName := *templateString
	folding := *foldPref
	outlines := *outlinePref
	rgbString := strings.Split(*rgbPref, ":")

	//Parse the rgbstring
	var rDefined, gDefined, bDefined uint8
	tempR, err := strconv.Atoi(rgbString[0])
	check(err)
	tempG, err := strconv.Atoi(rgbString[1])
	check(err)
	tempB, err := strconv.Atoi(rgbString[2])
	check(err)

	rDefined = uint8(tempR)
	gDefined = uint8(tempG)
	bDefined = uint8(tempB)

	//Open the templateFile
	templateFile, err := os.Open("Templates/" + templateName + ".png")
	check(err)
	defer templateFile.Close()

	//Prepare the generation directories for the file here.
	DirString := "GenerationDirectory/" + templateName
	CurrentDir, err := filepath.Abs("")
	check(err)
	PlacementDirectory := filepath.Join(CurrentDir, DirString)
	_, err = folderExists(PlacementDirectory)
	if err == nil {
		os.Mkdir(DirString, 0755)
	}
	IndividualSpriteDir := DirString + "/Individuals"
	_, err = folderExists(IndividualSpriteDir)
	if err == nil {
		os.Mkdir(IndividualSpriteDir, 0755)
	}

	//Grab our template pixels and the template config
	templateStream, err := png.Decode(templateFile)
	check(err)
	//Be kind, Rewind
	templateFile.Seek(0, 0)
	templateConfig, err := png.DecodeConfig(templateFile)
	check(err)

	//Grab our values to build an individual new image.
	var canvasWidth int
	var foldAt int
	switch folding {
	case "Even":
		canvasWidth = templateConfig.Width * 2
		foldAt = canvasWidth / 2
	case "Odd":
		canvasWidth = (templateConfig.Width * 2) - 1
		foldAt = (canvasWidth / 2) + 1
	case "None":
		canvasWidth = templateConfig.Width
		foldAt = canvasWidth
	}
	canvasHeight := templateConfig.Height

	//With that information, we can read out the values of our templatefile.  After a few differnet ideas of how to optimize,
	//I think the best practice is just to create a single list that captures our template, with enumerated values to delineate
	//which pixels are outlines, bits, ect.
	var pixelList []int
	for y := 0; y < templateConfig.Height; y++ {
		for x := 0; x < templateConfig.Width; x++ {
			aPixel := templateStream.At(x, y)
			m, n, o, _ := aPixel.RGBA()
			r := uint8(m)
			g := uint8(n)
			b := uint8(o)
			//Since we're reading them as rows, we simply append to pixelList as we go along.
			if r == 255 && g == 0 && b == 0 {
				pixelList = append(pixelList, Outline)
			} else if r == 0 && g == 0 && b == 0 {
				pixelList = append(pixelList, Bit)
			} else if r == 0 && g == 0 && b == 255 {
				pixelList = append(pixelList, Fill)
			} else if r == 0 && g == 255 && b == 0 {
				pixelList = append(pixelList, Accent)
			} else {
				pixelList = append(pixelList, Transparent)
			}
		}
	}
	//While we used Graphics Magic initially for creating a sprite sheet, we can probably incorporate this into the initial
	//encoding, getting rid of the external dependency.  Composite width is the number of images on each row.
	compositeWidth := *compositePref
	composite := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth * compositeWidth, canvasHeight * 256 / compositeWidth}})

	//While we could benefit from making fewer work groups of routines, I don't find the performance penalty
	//on smaller files as too painful when compared to the gains this makes on larger files.  This was 2x faster
	//than a sequential write for a 5x5 image, and was very fast on a 77x154 image.
	var wg sync.WaitGroup
	wg.Add(256)

	for i := 0; i < 256; i++ {
		go func(i int) {

			defer wg.Done()
			//Create the bit array for each individual image.
			bitArray := make([]bool, 8)
			//Here, we are encoding our i value into binary, then using bitwise shifts to place it into
			//an array.  We'll then use this array to decide whether we want our 'Bit' pixels to represented.
			//Technically the array is backwards, but since we'll end up reading it left to right it's not a problem
			for j := 0; j < 8; j++ {
				k := uint(i) >> j
				if k&1 == 1 {
					bitArray[j] = true
				}
			}

			//We'll create the modified template based on our pixel list, where we modify our outlines based
			//on the status of nearby bits.
			var newImage []int

			//So we start by copying pixelList to newImage, applying the bitArray to our designated bit pixels.
			//Since I want to keep open the option for larger images, we'll keep track of bitsRead, so that we
			//just repeat through our bitArray if we include more than 8 'Bit' pixels.
			bitsRead := 0
			for j := 0; j < len(pixelList); j++ {
				if pixelList[j] == Bit {
					if bitArray[bitsRead%8] == false {
						newImage = append(newImage, Outline)
					} else {
						newImage = append(newImage, Bit)
					}
					bitsRead++
				} else {
					newImage = append(newImage, pixelList[j])
				}
			}

			//I've gone through a few implementations, but have landed on this for drawing our outlines.  We
			//just want to check if it is a colored pixel, and if so, then we check if there are any transparent
			//pixels adjacent.  This should also be an optional process, in case the user does not want to have an
			//outline
			if outlines == true {
				for j := 0; j < len(newImage); j++ {
					if newImage[j] == Bit || newImage[j] == Fill {
						//Here we want to check for bit pixels across the cardinal directions. I opted for
						//simplifying the loop to check these combinations, so we translate our index into a coordinate.
						pixelCoord := image.Point{(j % templateConfig.Width), int(j / templateConfig.Width)}
						for k := -1; k < 2; k = k + 2 {
							//another benefit of translation, easier to check for whether a pixel is
							//actually adjacent or whether the next pixel is on the next different row
							var xIndex int
							var yIndex int
							if templateConfig.Width > pixelCoord.X+k && pixelCoord.X+k >= 0 {
								xIndex = pixelCoord.X + k + (pixelCoord.Y * templateConfig.Width)
							} else {
								//Outside our bounds/not adjacent?  Mark it as -1 and move on.
								xIndex = -1
							}
							if templateConfig.Height > pixelCoord.Y+k && pixelCoord.Y+k >= 0 {
								yIndex = pixelCoord.X + ((pixelCoord.Y + k) * templateConfig.Width)
							} else {
								yIndex = -1
							}
							//Is it in bounds?  Check for a transparent pixel and replace it with an Outline
							if xIndex >= 0 {
								if newImage[xIndex] == Transparent {
									newImage[xIndex] = Outline
								}
							}
							if yIndex >= 0 {
								if newImage[yIndex] == Transparent {
									newImage[yIndex] = Outline
								}
							}
						}
					}
				}
			}

			// With the template adjusted, we create the output file for each image.
			newFile := IndividualSpriteDir + "/" + strconv.Itoa(i) + ".png"

			outfile, err := os.Create(newFile)
			check(err)

			canvas := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth, canvasHeight}})

			//Instead of writing our picture out as a list, we are using the x and y loops to more
			//easily fold our images.
			for y := 0; y < canvasHeight; y++ {
				for x := 0; x < canvasWidth; x++ {
					//We want to start by converting our coordinate into an index position.  When we fold,
					//we put our index at the mirrored position.
					var pixelIndex int
					if x < foldAt {
						pixelIndex = x + (y * templateConfig.Width)
					} else {
						pixelIndex = (canvasWidth - x) + (y * templateConfig.Width) - 1
					}
					// We then check what our newImage list says about this index.  Outline are black pixels,
					// while Bit and Fill have their own formulas.  Since a canvas begins Transparent,
					// we can skip those pixel indices
					if newImage[pixelIndex] != Transparent {
						var r, g, b uint8
						if newImage[pixelIndex] == Outline {
							r = 0
							g = 0
							b = 0
						} else if newImage[pixelIndex] == Bit {
							//TODO:  Add choices for how to color.
							//two Linear gradients over YCbCr at .5 lumia (Rainbow)
							//r, g, b = color.YCbCrToRGB(uint8(128), uint8((i+128)%256), uint8(i%256))
							//the opposing gradients
							//r, g, b = color.YCbCrToRGB(uint8(128), uint8(255-i), uint8((i+128)%256))
							if rDefined+gDefined+bDefined == 0 {
								//Let's try something really wacky.  All the colors!
								if i < 128 {
									r, g, b = color.YCbCrToRGB(uint8(128), uint8((i*2+128)%256), uint8((i*2)%256))
								} else {
									r, g, b = color.YCbCrToRGB(uint8(128), uint8(255-i*2), uint8((i*2+128)%256))
								}
							} else {
								r, g, b = rDefined, gDefined, bDefined
							}
						} else if newImage[pixelIndex] == Fill {
							//r, g, b = color.YCbCrToRGB(uint8(128), uint8((i+128)%256), uint8(i%256))
							//r, g, b = color.YCbCrToRGB(uint8(128), uint8(255-i), uint8((i+128)%256))
							if i < 128 {
								r, g, b = color.YCbCrToRGB(uint8(128), uint8((i*2+128)%256), uint8((i*2)%256))
							} else {
								r, g, b = color.YCbCrToRGB(uint8(128), uint8(255-i*2), uint8((i*2+128)%256))
							}
						} else if newImage[pixelIndex] == Accent {
							//r, g, b = color.YCbCrToRGB(uint8(96), uint8((i+128)%256), uint8(i%256))
							//r, g, b = color.YCbCrToRGB(uint8(128), uint8(255-i), uint8((i+128)%256))
							if i < 128 {
								r, g, b = color.YCbCrToRGB(uint8(32), uint8((i*2+128)%256), uint8((i*2)%256))
							} else {
								r, g, b = color.YCbCrToRGB(uint8(32), uint8(255-i*2), uint8((i*2+128)%256))
							}
						}
						//assign RGB to the pixels
						newColor := [4]uint8{r, g, b, 255}
						yOffset := y * canvasWidth * 4
						iOffset := ((i / compositeWidth) * canvasHeight * canvasWidth * compositeWidth) + ((i % compositeWidth) * canvasWidth)
						for z := 0; z < 4; z++ {
							//Find our spot for the x axis, multiply by the bytes contained in the pixel, offset by the value of the row we're
							//on and the pixel we want to write to.
							canvas.Pix[(x*4+z)+yOffset] = newColor[z]
							//This horrific line of code creates the composite sprite sheet.  Essentially the same as above, but also offseting by the
							//y- and x-offsets of the sprite sheet using our increment.
							composite.Pix[((iOffset+x)*4)+
								z+(yOffset*compositeWidth)] = newColor[z]
						}
					}
				}
			}
			//After building the sprite, you'd think we'd be done
			png.Encode(outfile, canvas)
			outfile.Close()
		}(i)
	}

	wg.Wait()

	compositeName := PlacementDirectory + "/" + templateName + "SpriteSheet.png"
	compositeFile, err := os.Create(compositeName)
	check(err)
	png.Encode(compositeFile, composite)
	compositeFile.Close()
}
