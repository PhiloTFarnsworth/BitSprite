package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/muesli/gamut"
)

type Pixel int

//These values describe our template pixels.
const (
	Background Pixel = iota
	Bit
	Accent
	Fill
	Outline
	Delimiter
	PixelsDefined //Add any new tracked pixels above this.
)

var Black = color.RGBA{0, 0, 0, 255}
var Red = color.RGBA{255, 0, 0, 255}
var Green = color.RGBA{0, 255, 0, 255}
var Blue = color.RGBA{0, 0, 255, 255}
var Magenta = color.RGBA{255, 0, 255, 255}
var White = color.RGBA{255, 255, 255, 255}
var Transp = color.RGBA{0, 0, 0, 0}
var LGray = color.RGBA{85, 85, 85, 255}
var HGray = color.RGBA{170, 170, 170, 255}

//Flags.  Trying to be a bit more terse than the readme.  out- names are probably too abundant, and legacy is not ideal.
var templateString = flag.String("template", "", "Choose template to render, template must be in Templates folder.")
var foldPref = flag.String("fold", "", "Sets fold preference for template if desired, use even and odd, all other values default to no fold. (e, even=Even; o, odd=Odd)")
var vertFoldPref = flag.String("vertfold", "", "Sets fold preference accross bottom of image, use even and odd, all other values default to no fold. (e, even=Even; o, odd=Odd)")
var colorPref = flag.String("color", "", "Sets color of activated bit pixels, use Hex or Hex:Hex (#FFFFFF or #000000:#FFFFFF).")
var accentPref = flag.String("accent", "", "Sets the color of the accent pixels, use Hex or Hex:Hex (#FFFFFF or #000000:#FFFFFF).")
var fillPref = flag.String("fill", "", "Sets the color of the fill pixels,  use Hex or Hex:Hex (#FFFFFF or #000000:#FFFFFF).")
var backgroundPref = flag.String("background", "", "Sets color of background,  use Hex or Hex:Hex (#FFFFFF or #000000:#FFFFFF).")
var outlineColorPref = flag.String("outcolor", "#000000", "Sets the color of the outline pixels,  use Hex or Hex:Hex (#FFFFFF or #000000:#FFFFFF).")
var outlinePref = flag.Bool("outline", true, "Sets outline preference, use Golang Bool values.")
var upscalePref = flag.Int("upscale", 1, "Increases the scale of the template's copies, use a positive integer.")
var compositePref = flag.Int("sheetwidth", 16, "Sets width of output sprite sheet, use a factor of 256.")
var legacyColors = flag.Bool("legacy", false, "Colors are based on a composite linear gradient of the YCbCr at .5 lumia if true, use Golang Bool values.")
var outputNamePref = flag.String("outname", "", "Sets the output files to be placed in a generation directory named after the string provided.")
var individualsPref = flag.Bool("individuals", false, "Creates a directory of individual .png files for each image on the spritesheet")
var randSeedPref = flag.Bool("randseed", true, "Toggles random seed, used for debug/testing.")

//var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")

func main() {
	// //Profiling
	flag.Parse()
	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	pprof.StartCPUProfile(f)
	// 	defer pprof.StopCPUProfile()
	// }

	// //Trace
	// f, err := os.Create("trace.out")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// err = trace.Start(f)
	// if err != nil {
	// 	panic(err)
	// }
	// defer trace.Stop()

	//rand seed
	randomSeed := *randSeedPref
	if randomSeed {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(1)
	}
	//Grab the flag values
	templateName := *templateString
	folding := *foldPref
	outlines := *outlinePref
	legacy := *legacyColors
	upScale := *upscalePref
	compositeWidth := *compositePref
	outputName := *outputNamePref
	individuals := *individualsPref
	vertFold := *vertFoldPref

	//We'll preemptively break down our colors strings as though they were blend values.  We'll use our
	//enumerated Pixel values to put them on a temporary map.  Sorta chunky, but it's readable enough.
	chosenColorStrings := make(map[Pixel][]string)
	chosenColorStrings[Bit] = strings.Split(*colorPref, ":")
	chosenColorStrings[Accent] = strings.Split(*accentPref, ":")
	chosenColorStrings[Fill] = strings.Split(*fillPref, ":")
	chosenColorStrings[Background] = strings.Split(*backgroundPref, ":")
	chosenColorStrings[Outline] = strings.Split(*outlineColorPref, ":")

	//Converts those hexes into colors.
	chosenColors := make(map[Pixel][]color.Color)
	for key, val := range chosenColorStrings {
		if len(val) == 1 {
			//We'll use these default values if nothing is defined
			if val[0] == "" {
				switch key {
				case Bit:
					chosenColors[key] = append(chosenColors[key], White)
				case Accent:
					chosenColors[key] = append(chosenColors[key], LGray)
				case Fill:
					chosenColors[key] = append(chosenColors[key], HGray)
				case Background:
					chosenColors[key] = append(chosenColors[key], Transp)
				case Outline:
					chosenColors[key] = append(chosenColors[key], Black)
				}
			} else {
				//Otherwise add one color to the chosen colors list
				chosenColors[key] = append(chosenColors[key], gamut.Hex(val[0]))
			}
		} else {
			//Add the blend to the list of chosen colors.  Should consider doing multiple blends.
			chosenColors[key] = gamut.Blends(gamut.Hex(val[0]), gamut.Hex(val[1]), 256)
		}
	}

	//There's a few ways we can handle bad sheetwidth flags, defaulting to 16 is one solution.
	if 256%compositeWidth != 0 || compositeWidth > 256 || compositeWidth < 1 {
		compositeWidth = 16
		fmt.Print("Bad sheetWidth passed, defaulting to sheetWidth=16\n")
	}

	//sanitize upScale
	if upScale < 1 {
		upScale = 1
	}
	//Open the templateFile
	currentDir, err := filepath.Abs("")
	check(err)
	templateFile, err := os.Open(filepath.Join(currentDir, "/Templates/"+templateName+".png"))
	check(err)
	defer templateFile.Close()

	var dirString string
	var PlacementDirectory string
	var individualSpriteDir string
	//Prepare the generation directories for the file here.
	if strings.EqualFold(outputName, "docs") {
		dirString = "docs/example"
		PlacementDirectory = filepath.Join(currentDir, dirString)
		mayCreateFolder(PlacementDirectory)
		individualSpriteDir = filepath.Join(currentDir, dirString+"/Individuals")
		if individuals {
			mayCreateFolder(individualSpriteDir)
		}
	} else {
		if outputName != "" {
			templateName = outputName
		}
		generationDirectory := filepath.Join(currentDir, "/GenerationDirectory")
		mayCreateFolder(generationDirectory)
		dirString = "GenerationDirectory/" + templateName
		PlacementDirectory = filepath.Join(currentDir, dirString)
		mayCreateFolder(PlacementDirectory)
		individualSpriteDir = filepath.Join(currentDir, dirString+"/Individuals")
		if individuals {
			mayCreateFolder(individualSpriteDir)
		}
	}

	//Grab our template pixels and the template config
	templateStream, err := png.Decode(templateFile)
	check(err)
	templateFile.Seek(0, 0)
	templateConfig, err := png.DecodeConfig(templateFile)
	check(err)

	//Use folding to determine the dimensions of the output images.
	var canvasWidth int
	var canvasHeight int
	var foldY int
	var foldX int
	if strings.EqualFold(folding, "even") || strings.EqualFold(folding, "e") {
		canvasWidth = (templateConfig.Width * 2)
		foldY = templateConfig.Width
	} else if strings.EqualFold(folding, "odd") || strings.EqualFold(folding, "o") {
		canvasWidth = ((templateConfig.Width * 2) - 1)
		foldY = (canvasWidth / 2) + 1
	} else {
		canvasWidth = templateConfig.Width
		foldY = canvasWidth
	}

	if strings.EqualFold(vertFold, "even") || strings.EqualFold(vertFold, "e") {
		canvasHeight = templateConfig.Height * 2
		foldX = templateConfig.Height
	} else if strings.EqualFold(vertFold, "odd") || strings.EqualFold(vertFold, "o") {
		canvasHeight = (templateConfig.Height * 2) - 1
		foldX = (canvasHeight / 2) + 1
	} else {
		canvasHeight = templateConfig.Height
		foldX = canvasHeight
	}

	//Translate the image into a simple array.
	var templateArray []Pixel
	var delimiters []int //indexes where we want to change our bit array
	for y := 0; y < templateConfig.Height; y++ {
		for x := 0; x < templateConfig.Width; x++ {
			//Convert pixel model to RGBA.
			aPixel := color.RGBAModel.Convert(templateStream.At(x, y))
			//We compare the template's pixels to our defined colors, then append them to templateArray
			switch aPixel {
			case Red:
				templateArray = append(templateArray, Outline)
			case Green:
				templateArray = append(templateArray, Accent)
			case Blue:
				templateArray = append(templateArray, Fill)
			case Black:
				templateArray = append(templateArray, Bit)
			case Magenta:
				templateArray = append(templateArray, Background)
				delimiters = append(delimiters, x+y*templateConfig.Width)
			default:
				templateArray = append(templateArray, Background)
			}
		}
	}

	//Generate number list for delimited segments of the input image.
	var randomArrays [][]int
	for i := 0; i < len(delimiters); i++ {
		randomArrays = append(randomArrays, rand.Perm(256))
	}
	//composite is our sprite sheet, and we'll draw it up simultaneously with our individual images.
	composite := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth * upScale * compositeWidth, canvasHeight * upScale * 256 / compositeWidth}})

	//This is admittedly lazy, but as it stands I don't have a great solution in mind for scaling wait groups based on the pixels we write.  There is definitely a
	//point where you gain some extra performance by using fewer wait groups that have responsibility for multiple images, but it's a little fuzzy and probably
	//not worth the testing time and added code complexity to find those points.
	var wg sync.WaitGroup
	wg.Add(256)

	for i := 0; i < 256; i++ {
		go func(i int) {
			defer wg.Done()
			//newImage will hold a modified template array, based on how we read our bit pixels and our outline settings.
			var newImage []Pixel
			bitsRead := 0
			resolutionNumber := i
			delimitersRead := 0
			for j := 0; j < len(templateArray); j++ {
				if returnIndex(delimiters, j) != -1 {
					delimitersRead = returnIndex(delimiters, j)
					bitsRead = 0
					resolutionNumber = randomArrays[delimitersRead][i]
				}
				if templateArray[j] == Bit {
					//We take our increment, shift it by the bitsRead, finally checking whether it is even or odd.  This way 0 = all inactive,
					//255 = all active.
					if (resolutionNumber>>(bitsRead%8))&1 == 0 {
						newImage = append(newImage, Outline)
					} else {
						newImage = append(newImage, Bit)
					}
					bitsRead++
				} else {
					newImage = append(newImage, templateArray[j])
				}

			}
			//checks neighbors of active, colored pixels.  If the neighboring pixel is a background, replace it with an outline
			//pixel.  Disabled by -outline=false
			if outlines {
				for j := 0; j < len(newImage); j++ {
					if newImage[j] == Bit || newImage[j] == Fill || newImage[j] == Accent {
						for k := -1; k < 2; k = k + 2 {
							//left, right; the remainder of the index gives us an x coordinate
							if templateConfig.Width > (j%templateConfig.Width)+k && (j%templateConfig.Width)+k >= 0 {
								//if x-coord is good we can just add k to our index
								if newImage[j+k] == Background {
									newImage[j+k] = Outline
								}
							}
							//up, down; use int()'s inherent round down ability to create our y coordinate
							if templateConfig.Height > int(j/templateConfig.Width)+k && int(j/templateConfig.Width)+k >= 0 {
								//looks gross, but it works.
								yIndex := (j % templateConfig.Width) + ((int(j/templateConfig.Width) + k) * templateConfig.Width)
								if newImage[yIndex] == Background {
									newImage[yIndex] = Outline
								}
							}
						}
					}
				}
			}
			//TODO: Reduce Option. Here we would run through the image again to reduce
			var outfile *os.File
			var canvas *image.RGBA
			if individuals {
				// With the template adjusted, we create the output file for each image.
				newFile := individualSpriteDir + "/" + strconv.Itoa(i) + ".png"
				outfile, err = os.Create(newFile)
				check(err)
				canvas = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth * upScale, canvasHeight * upScale}})
			}

			//let's grab the base color for our image
			var finalColors [PixelsDefined][]color.Color
			var placeholderIndex int
			if len(delimiters) > 0 {
				placeholderIndex = len(delimiters)
			} else {
				placeholderIndex = 1
			}
			for j := 0; j < placeholderIndex; j++ {
				if len(delimiters) == 0 {
					resolutionNumber = i
				} else {
					resolutionNumber = randomArrays[j][i]
				}
				if !legacy {
					for key, val := range chosenColors {
						if len(val) > 1 {
							finalColors[key] = append(finalColors[key], chosenColors[key][resolutionNumber])
						} else {
							finalColors[key] = append(finalColors[key], chosenColors[key][0])
						}
					}
				} else {
					//legacy ycbcr gradients
					finalColors[Bit] = append(finalColors[Bit], color.YCbCr{128, uint8((resolutionNumber + 128) % 256), uint8(resolutionNumber % 256)})
					finalColors[Accent] = append(finalColors[Accent], color.YCbCr{64, uint8((resolutionNumber + 128) % 256), uint8(resolutionNumber % 256)})
					finalColors[Fill] = append(finalColors[Fill], color.YCbCr{192, uint8((resolutionNumber + 128) % 256), uint8(resolutionNumber % 256)})
					finalColors[Background] = append(finalColors[Background], Transp)
					finalColors[Outline] = append(finalColors[Outline], Black)
				}
			}

			//Finally, with colors and a template secured, we can write to our individual canvas and collective composite.
			var pixelIndex int
			for y := 0; y < canvasHeight; y++ {
				for x := 0; x < canvasWidth; x++ {
					//We want to start by converting our coordinate into an index position.  When we fold,
					//we put our index at the mirrored position.
					if x < foldY {
						if y < foldX {
							pixelIndex = x + (y * templateConfig.Width)
						} else {
							pixelIndex = x + ((canvasHeight - y - 1) * templateConfig.Width)
						}
					} else {
						if y < foldX {
							pixelIndex = (canvasWidth - x) + (y * templateConfig.Width) - 1
						} else {
							pixelIndex = (canvasWidth - x) + ((canvasHeight - y - 1) * templateConfig.Width) - 1
						}
					}

					if y < foldX {
						if returnIndex(delimiters, pixelIndex) != -1 {
							delimitersRead = returnIndex(delimiters, pixelIndex)
						}
					} else {
						//when we flip, we need to consider that we're reading upside down, so adjust
						// our pixel index down
						modifiedPixelIndex := x + ((canvasHeight - y) * templateConfig.Width)
						if returnIndex(delimiters, modifiedPixelIndex) != -1 {
							delimitersRead = returnIndex(delimiters, modifiedPixelIndex) - 1
						} //Hacky hack for reading that last delimiter
						if delimitersRead == -1 {
							delimitersRead = 0
						}
					}

					//A little messy, but we account for upScale here.
					for j := 0; j < upScale; j++ {
						for k := 0; k < upScale; k++ {
							if individuals {
								canvas.Set((x*upScale)+j, (y*upScale)+k,
									finalColors[newImage[pixelIndex]][delimitersRead])
							}
							composite.Set((x*upScale)+j+canvasWidth*upScale*(i%compositeWidth),
								(y*upScale)+k+canvasHeight*upScale*(i/compositeWidth),
								finalColors[newImage[pixelIndex]][delimitersRead])
						}
					}
				}
			}
			//After building the sprite, we encode, then close the individual sprite file.
			if individuals {
				png.Encode(outfile, canvas)
				outfile.Close()
			}
		}(i)
	}
	wg.Wait()
	compositeName := PlacementDirectory + "/" + templateName + "SpriteSheet.png"
	compositeFile, err := os.Create(compositeName)
	check(err)
	png.Encode(compositeFile, composite)
	compositeFile.Close()
}

//Very generic check function to reduce boilerplate.  Since we are creating files, I figure we err on the side of caution and
//just fatal log any errors that come.
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

// return index of matched value, otherwise return -1
func returnIndex(list []int, find int) int {
	i := 0
	for i < len(list) {
		if find == list[i] {
			return i
		}
		i++
	}
	return -1
}
