package main

import (
	"errors"
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

	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
)

//Very generic check function to reduce boilerplate.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Creates folder if it does not already exist.
func mayCreateFolder(path string) {
	_, err := os.Stat(path)
	if err == nil {
		//folder exists
	} else if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(path, 0755)
	} else {
		//Shadow realm
		log.Fatal(err)
	}
}

//These values describe our template pixels.
const (
	_ = iota
	Transparent
	Bit
	Accent
	Fill
	Outline
)

//Colors to match for baf and outline
var Black = color.RGBA{0, 0, 0, 255}
var Red = color.RGBA{255, 0, 0, 255}
var Green = color.RGBA{0, 255, 0, 255}
var Blue = color.RGBA{0, 0, 255, 255}
var Transp = color.RGBA{0, 0, 0, 0}

//In my mind, I think it might be better if we initialize our flags in a single map.
var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
var templateString = flag.String("template", "", "Choose template to render")
var foldPref = flag.String("fold", "None", "Sets fold preference for template (Even/Odd/None)")
var outlinePref = flag.Bool("outline", true, "Sets outline preference")
var hexPref = flag.String("hex", "#FFFFFF", "Sets colors based on a specified hex value.  if both -color and -hex are used, -color will be used.  -Blend overrules both.")
var colorPref = flag.String("color", "", "Sets color based on a specified color from Wikipedia's list of colors. if both -color and -hex are used, -color will be used. -Blend overrules both.")
var backgroundPref = flag.String("background", "", "Sets color of background. use a Hex value or a color name from wikipedia's list of colors.")

//Analogous colors are "...in harmony with one another." according to a color chart I just googled, so I figure we promote world peace in this fashion.
var fabPref = flag.String("fab", "Analogous", "Describes the color relationship between Fill, Accent and Bit pixels (Analogous, Triadic, SplitComplementary)")
var blendPref = flag.String("blend", "x:x", "Shifts between two colors as the program iterates through copies of the template. use 'Color:Color' or 'Hex:Hex' (i.e. 'Red:Blue' or '#FF0000:#0000FF'), based on Wikipedia's list of colors.")
var compositePref = flag.Int("sheetWidth", 8, "Sets width of output sprite sheet, must return a whole number for 256/compositeWidth")
var legacyColors = flag.Bool("legacy", false, "Colors are based on a composite linear gradient of the YCbCr at .5 lumia if true")

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
	fmt.Print("BitSprite: Making 256 versions of 1 thing since 2020\n")
	//Grab the flag values
	templateName := *templateString
	folding := *foldPref
	outlines := *outlinePref
	legacy := *legacyColors
	mainHex := *hexPref
	mainColor := *colorPref
	background := *backgroundPref
	fab := *fabPref
	blend := strings.Split(*blendPref, ":")
	compositeWidth := *compositePref

	//There's a few ways we can handle bad compositePrefs, but I figure just defaulting to 8 is better
	// than quitting or throwing an error.
	if 256%compositeWidth != 0 {
		compositeWidth = 8
		fmt.Print("Bad sheetWidth passed, defaulting to sheetWidth=8\n")
	}

	//here we need to make sense of user inputs on hex/color, FAB and blend.  First we take the bitColor, then we
	//the FAB relatives (by default Analogous).  Though if we want to use the blend values... we need to parse that string first.
	var bitColor color.Color
	var fillColor color.Color
	var accentColor color.Color
	//Looks like being a hack is back on the menu boys!
	var startColor color.Color
	var endColor color.Color
	var bgColor color.Color
	var bitColorList []color.Color
	var ok bool

	//First blend check
	if strings.HasPrefix(blend[0], "#") {
		startColor = gamut.Hex(blend[0])
		endColor = gamut.Hex(blend[1])
	} else {
		startColor, ok = palette.Wikipedia.Color(blend[0])
		if !ok {
			startColor = nil
		}
		endColor, ok = palette.Wikipedia.Color(blend[1])
		if !ok {
			endColor = nil
		}
	}

	//If we're blending, we'll grab the main color's blends, otherwise, we'll check -color and -hex.  If nothing, we set
	//everything to white.
	if startColor != nil && endColor != nil {
		bitColorList = gamut.Blends(startColor, endColor, 256)
		bitColor = nil
	} else {
		if mainColor == "" {
			bitColor = gamut.Hex(mainHex)
		} else {
			bitColor, ok = palette.Wikipedia.Color(mainColor)
			if !ok {
				bitColor, _ = palette.Wikipedia.Color("White")
			}
		}
		//check our fab to populate our single chromatic combos
		switch fab {
		case "Analogous":
			colorList := gamut.Analogous(bitColor)
			//arbitrary choice here.
			fillColor = colorList[0]
			accentColor = colorList[1]
		case "Triadic":
			colorList := gamut.Triadic(bitColor)
			fillColor = colorList[0]
			accentColor = colorList[1]
		case "SplitComplementary":
			colorList := gamut.SplitComplementary(bitColor)
			fillColor = colorList[0]
			accentColor = colorList[1]
		}
	}
	//finally parse the -background flag.
	if strings.HasPrefix(background, "#") {
		bgColor = gamut.Hex(background)
	} else {
		bgColor, ok = palette.Wikipedia.Color(background)
		if !ok {
			bgColor = Transp
		}
	}

	//Open the templateFile
	templateFile, err := os.Open("Templates/" + templateName + ".png")
	check(err)
	defer templateFile.Close()

	//Prepare the generation directories for the file here
	currentDir, err := filepath.Abs("")
	check(err)
	genDirString := "/GenerationDirectory"
	generationDirectory := filepath.Join(currentDir, genDirString)
	mayCreateFolder(generationDirectory)
	DirString := "GenerationDirectory/" + templateName
	PlacementDirectory := filepath.Join(currentDir, DirString)
	mayCreateFolder(PlacementDirectory)
	IndividualSpriteString := DirString + "/Individuals"
	individualSpriteDir := filepath.Join(currentDir, IndividualSpriteString)
	mayCreateFolder(individualSpriteDir)

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
			//We compare the template's pixels to our defined colors, then append them to pixelList
			switch aPixel {
			case Red:
				pixelList = append(pixelList, Outline)
			case Green:
				pixelList = append(pixelList, Accent)
			case Blue:
				pixelList = append(pixelList, Fill)
			case Black:
				pixelList = append(pixelList, Bit)
			default:
				pixelList = append(pixelList, Transparent)
			}
		}
	}
	//While we used Graphics Magic initially for creating a sprite sheet, we can probably incorporate this into the initial
	//encoding, getting rid of the external dependency.  Composite width is the number of images on each row.

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
					if !bitArray[bitsRead%8] {
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
			if outlines {
				for j := 0; j < len(newImage); j++ {
					if newImage[j] == Bit || newImage[j] == Fill || newImage[j] == Accent {
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
			newFile := individualSpriteDir + "/" + strconv.Itoa(i) + ".png"

			outfile, err := os.Create(newFile)
			check(err)

			canvas := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth, canvasHeight}})

			//Instead of writing our picture out as a list, we are using the x and y loops to more
			//easily fold our images.
			var pixelIndex int
			//let's grab the base color for our image
			var fColor color.Color
			var aColor color.Color
			var bColor color.Color
			if !legacy {
				if bitColor != nil {
					fColor = fillColor
					aColor = accentColor
					bColor = bitColor
				} else {
					bColor = bitColorList[i]
					switch fab {
					case "Analogous":
						colorList := gamut.Analogous(bColor)
						fColor = colorList[0]
						aColor = colorList[1]
					case "SplitComplementary":
						colorList := gamut.SplitComplementary(bColor)
						fColor = colorList[0]
						aColor = colorList[1]
					case "Triadic":
						colorList := gamut.Triadic(bColor)
						fColor = colorList[0]
						aColor = colorList[1]
					}
				}
			} else {
				//legacy ycbcr gradient
				bColor = color.YCbCr{128, uint8((i + 128) % 256), uint8(i % 256)}
				fColor = color.YCbCr{156, uint8((i + 128) % 256), uint8(i % 256)}
				aColor = color.YCbCr{32, uint8((i + 128) % 256), uint8(i % 256)}
				//Going to keep this as a note, While I prefer the above function, you could also
				//argue this covers more of the color plane.  On the otherhand, the extra color space
				//is almost entirely green and magenta, so it's not terribly interesting to double our
				//gradients.
				// 	if i < 128 {
				// 		bcolor = color.YCbCr{uint8(128), uint8((i*2+128)%256), uint8((i*2)%256)}
				// 	} else {
				// 		bcolor = color.YCbCrToRGB(uint8(128), uint8(255-i*2), uint8((i*2+128)%256))
				// 	}

			}
			for y := 0; y < canvasHeight; y++ {
				for x := 0; x < canvasWidth; x++ {
					//We want to start by converting our coordinate into an index position.  When we fold,
					//we put our index at the mirrored position.
					if x < foldAt {
						pixelIndex = x + (y * templateConfig.Width)
					} else {
						pixelIndex = (canvasWidth - x) + (y * templateConfig.Width) - 1
					}
					//looks great compared to the legacy interpretation!
					switch newImage[pixelIndex] {
					case Outline:
						canvas.Set(x, y, Black)
						composite.Set(x+canvasWidth*(i%8), y+canvasHeight*(i/8), Black)
					case Bit:
						canvas.Set(x, y, bColor)
						composite.Set(x+canvasWidth*(i%8), y+canvasHeight*(i/8), bColor)
					case Accent:
						canvas.Set(x, y, aColor)
						composite.Set(x+canvasWidth*(i%8), y+canvasHeight*(i/8), aColor)
					case Fill:
						canvas.Set(x, y, fColor)
						composite.Set(x+canvasWidth*(i%8), y+canvasHeight*(i/8), fColor)
					case Transparent:
						if background != "" {
							canvas.Set(x, y, bgColor)
							composite.Set(x+canvasWidth*(i%8), y+canvasHeight*(i/8), bgColor)
						}
					}

				}
			}
			//After building the sprite, we encode, then close the individual sprite file.
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
