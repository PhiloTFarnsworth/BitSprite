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
)

type Pixel int

//These values describe our template pixels.
const (
	Background Pixel = iota
	Bit
	Accent
	Fill
	Outline
)

//Colors to match for fab and outline
var Black = color.Color(color.RGBA{0, 0, 0, 255})
var Red = color.Color(color.RGBA{255, 0, 0, 255})
var Green = color.RGBA{0, 255, 0, 255}
var Blue = color.RGBA{0, 0, 255, 255}
var Transp = color.RGBA{0, 0, 0, 0}
var LGray = color.RGBA{85, 85, 85, 255}
var HGray = color.RGBA{170, 170, 170, 255}
var White = color.RGBA{255, 255, 255, 255}

//The big idea here is that we want some granular control for the user, but without too much extra.  Ideally, the user
//could do bitsprite.exe -template=somefile and have something interesting pop out, while another user could load up on
//flags and have their specific needs met.
var templateString = flag.String("template", "", "Choose template to render")
var foldPref = flag.String("fold", "", "Sets fold preference for template if desired. (e=Even, o=odd)")
var colorPref = flag.String("color", "", "Sets color of activated bit pixels.")
var accentPref = flag.String("accent", "", "Sets the color of the accent pixels.")
var fillPref = flag.String("fill", "", "Sets the color of the fill pixels.")
var backgroundPref = flag.String("background", "", "Sets color of background.")
var outlineColorPref = flag.String("outcolor", "#000000", "Sets the color of the outline pixels.")
var outlinePref = flag.Bool("outline", true, "Sets outline preference")
var upscalePref = flag.Int("upscale", 1, "Increases the scale of the template's copies")
var compositePref = flag.Int("sheetwidth", 16, "Sets width of output sprite sheet, must return a whole number for 256/compositeWidth")
var legacyColors = flag.Bool("legacy", false, "Colors are based on a composite linear gradient of the YCbCr at .5 lumia if true")
var outputNamePref = flag.String("outname", "", "Sets the output files to be named after this string instead of the template name")
var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")

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
	upScale := *upscalePref
	compositeWidth := *compositePref
	outputName := *outputNamePref
	//We're going to support blends for all of these variables, so first we'll pass our flags split as though it's blend
	//code.
	chosenColorStrings := make(map[Pixel][]string)
	chosenColorStrings[Bit] = strings.Split(*colorPref, ":")
	chosenColorStrings[Accent] = strings.Split(*accentPref, ":")
	chosenColorStrings[Fill] = strings.Split(*fillPref, ":")
	chosenColorStrings[Background] = strings.Split(*backgroundPref, ":")
	chosenColorStrings[Outline] = strings.Split(*outlineColorPref, ":")

	//There's a few ways we can handle bad compositePrefs, but defaulting to 8 is one solution.
	if 256%compositeWidth != 0 || compositeWidth > 256 {
		compositeWidth = 16
		fmt.Print("Bad sheetWidth passed, defaulting to sheetWidth=8\n")
	}

	chosenColors := make(map[Pixel][]color.Color)
	//We check for length of our flag values, if we split into multiple strings, we expect to blend.
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

	//sanitize upScale
	if upScale < 1 {
		upScale = 1
	}

	//Open the templateFile
	templateFile, err := os.Open("Templates/" + templateName + ".png")
	check(err)
	defer templateFile.Close()

	if outputName != "" {
		templateName = outputName
	}

	//Prepare the generation directories for the file here
	currentDir, err := filepath.Abs("")
	check(err)
	genDirString := "/GenerationDirectory"
	generationDirectory := filepath.Join(currentDir, genDirString)
	mayCreateFolder(generationDirectory)
	dirString := "GenerationDirectory/" + templateName
	PlacementDirectory := filepath.Join(currentDir, dirString)
	mayCreateFolder(PlacementDirectory)
	individualSpriteDir := filepath.Join(currentDir, dirString+"/Individuals")
	mayCreateFolder(individualSpriteDir)

	//Grab our template pixels and the template config
	templateStream, err := png.Decode(templateFile)
	check(err)
	templateFile.Seek(0, 0)
	templateConfig, err := png.DecodeConfig(templateFile)
	check(err)

	//Grab our values to build an individual new image.
	var canvasWidth int
	var foldAt int
	if strings.EqualFold(folding, "even") || strings.EqualFold(folding, "e") {
		canvasWidth = (templateConfig.Width * 2)
		foldAt = (canvasWidth / 2)
	} else if strings.EqualFold(folding, "odd") || strings.EqualFold(folding, "o") {
		canvasWidth = ((templateConfig.Width * 2) - 1)
		foldAt = ((canvasWidth / 2) + 1)
	} else {
		canvasWidth = templateConfig.Width
		foldAt = canvasWidth
	}
	canvasHeight := templateConfig.Height

	//create an array and assign our pixels to it.
	var pixelList []Pixel
	for y := 0; y < templateConfig.Height; y++ {
		for x := 0; x < templateConfig.Width; x++ {
			//This was a bizarre fix.  Generally, .PNGs almost always return pixels encoded as RGBA, but by some happenstance
			//I managed to create a .PNG which was read as NRGBA.  At any rate, this should work no matter what tomfoolery
			//happens when you create the template .png.
			aPixel := color.RGBAModel.Convert(templateStream.At(x, y))
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
				pixelList = append(pixelList, Background)
			}
		}
	}

	//composite is our sprite sheet, we'll draw it up simultaneously with our individual images.
	composite := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth * upScale * compositeWidth, canvasHeight * upScale * 256 / compositeWidth}})

	//While we could benefit from making fewer work groups of routines, I don't find the performance penalty
	//on smaller files as too painful when compared to the gains this makes on larger files.  This was 2x faster
	//than a sequential write for a 5x5 image, and was very fast on a 77x154 image.
	var wg sync.WaitGroup
	wg.Add(256)

	for i := 0; i < 256; i++ {
		go func(i int) {

			defer wg.Done()
			//We'll create the modified template based on our pixel list, where we modify our outlines based
			//on the status of nearby bits.
			var newImage []Pixel
			bitsRead := 0
			//I think the bitarray has been more for my benefit, I should be able to write this without the bitarray being set.
			for j := 0; j < len(pixelList); j++ {
				if pixelList[j] == Bit {
					//We take our increment, shift it by the bitsRead, finally checking whether it is even odd.  This way 0 = all inactive,
					//255 = all active.
					if (i>>(bitsRead%8))&1 == 0 {
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
			//just want to check if it is a colored pixel, and if so, then we check if there are any Background
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
							if templateConfig.Width > pixelCoord.X+k && pixelCoord.X+k >= 0 {
								xIndex := pixelCoord.X + k + (pixelCoord.Y * templateConfig.Width)
								if newImage[xIndex] == Background {
									newImage[xIndex] = Outline
								}
							}
							if templateConfig.Height > pixelCoord.Y+k && pixelCoord.Y+k >= 0 {
								yIndex := pixelCoord.X + ((pixelCoord.Y + k) * templateConfig.Width)
								if newImage[yIndex] == Background {
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

			canvas := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{canvasWidth * upScale, canvasHeight * upScale}})

			//Instead of writing our picture out as a list, we are using the x and y loops to more
			//easily fold our images.
			var pixelIndex int
			//let's grab the base color for our image
			finalColors := make(map[Pixel]color.Color)
			if !legacy {
				for key, val := range chosenColors {
					if len(val) > 1 {
						finalColors[key] = chosenColors[key][i]
					} else {
						finalColors[key] = chosenColors[key][0]
					}
				}
			} else {
				//legacy ycbcr gradients
				finalColors[Bit] = color.YCbCr{128, uint8((i + 128) % 256), uint8(i % 256)}
				finalColors[Accent] = color.YCbCr{64, uint8((i + 128) % 256), uint8(i % 256)}
				finalColors[Fill] = color.YCbCr{192, uint8((i + 128) % 256), uint8(i % 256)}
				finalColors[Background] = Transp
				finalColors[Outline] = Black
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
					//Messy.  Essentially we read the pixel index on our newImage, then we set the pixels on the actual image while
					//accomodating for scale.
					for j := 0; j < upScale; j++ {
						for k := 0; k < upScale; k++ {
							canvas.Set((x*upScale)+j, (y*upScale)+k, finalColors[newImage[pixelIndex]])
							composite.Set((x*upScale)+j+canvasWidth*upScale*(i%compositeWidth), (y*upScale)+k+canvasHeight*upScale*(i/compositeWidth), finalColors[newImage[pixelIndex]])
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
