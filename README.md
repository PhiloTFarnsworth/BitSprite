# BitSprite
## A Free 8-Bit Sprite Generator.  Create 256 variants from a single template .PNG

### What?
Bitsprite is a program that creates variants of an image based on prepared template image, as well
as creating a total sprite sheet of the resulatant images.
After downloading the package, or the .exe, the user can use the command prompt to run the program,
using flags to control the results.  The user has several options to designate colors, as well as
'unfold' the template to create symettrical images.

### Install
First, the user can either download the package, or they merely Bitsprite.Exe if they 
don't care to install the go language.  If just using the .EXE, make sure to create a 'Templates' 
folder, as BitSprite will look for a Templates folder in its directory to find templates.

### Using BitSprite
#### Creating a Template
To create a template, open your favorite image editor.  In this example, we'll just use the venerable
Microsoft Paint, but any editor should be fine, so long as you encode the template as a png.  Templates
have a special color coding to determine how the pixels are read and later written.

#### Pixel Legend
|Bit|Black (RGB 0,0,0)|
|Accent|Green (RGB 0,255,0)|
|Fill|Blue (RGB 0,0,255)|
|Outline|Red (RGB 255,0,0)|
|Background|White (RGB 255,255,255)|
Any deviation from these specific colors on a template will result in the offending color being treated
as background.  So if you output is entirely transparent, check that your pixels are correctly colored.

Bit pixels are how our generator creates the variety between images, as these pixels switch between
being active and displaying a set color, or inactive and becoming an outline pixel.  Accent and fill
pixels are static, allowing the user to set pixels that are always active.  Outline pixels are generated
any time an active (colored) pixel is bordered by a background pixel, and by default are colored black.
They can also function as another static pixel if explicitly included in the template, though they will not 
generate borders around themselves like bit/accent/fill pixels do.  Background pixels are default transparent.  
Let's look at what a compliant template may look like:

[MS Paint Window with an open template file](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/ExampleTemplateWorkflow.png)

This example shows a template image at 800% zoom, with each square of color representing one pixel.
Notice that there are 5 colors used in the 'face' image, all of which correspond to our control scheme. 

If we run 
'''
BitSprite.exe -template=face
'''
we can check the GenerationDirectory and find this:

[greyscale face spritesheet](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/FaceSpriteSheet.png)

It's a little small.  The transparent background also does no favors for visibility with this basic grayscale.
Instead let's return a large image with a better background.  We'll set the background to magenta and we'll scale
the image by 4, so that our original pixel is written as 4 pixels in a square.

'''
BitSprite.Exe -template=face -background=#ff00ff -upscale=4
'''

faceSpriteSheet.png should now look like this:
[greyscale face spritesheet with magenta background, 4 times scale](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/FaceSpriteSheetUp.png)

We can now see that our bit pixels are counting up in binary based on their position when read.  Our fill and accent
pixels remain the same through each iteration, with distinct colors.  The activated bit pixels also influence the
drawing of our outlines, though the red 'smile' from our template stays static.  

#### Now this is all well and good, but why did you use half a face?
Well, another nifty feature of bitsprite is the ability to reflect our templates while writing the new images.
if we decide to run:

'''
BitSprite.exe -template=face -background=#ff00ff -upscale=4 -fold=odd
'''

we return:

[greyscale face spritesheet with magenta background, 4 times scale unfolded](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/FaceSpriteSheetUpEx.png)

Now we have a symettrical image, and also one that allows your imagination to wander a little bit. 
Is it the many hairstyles of Tommy Pickles? Maybe it is a happy slime, expanding and contracting 
across the sprite sheet? You decide!

#### It's still a little bland
Fine, BitSprite also supports some cool ways to control colors.  Let's say for the time being we've decided it's
a face with facial hair.  (Values supplied by google)

'''
BitSprite.exe -template=face -upscale=4 -fold=odd -color=#9a3300 -accent=#3d671d -fill=#F1C27D
'''

[colored in face spritesheet 4 times scale, unfolded](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/FaceSpriteSheetColor.png)

With color, we create a passable face.  But we're not done yet, we can also do a blend of tones
across the entire sheet, for any of those values. 

'''
BitSprite.exe -template=face -upscale=4 -fold=odd -color=#9a3300:#4f1a00 -accent=#3d671d:#497665 -fill=#F1C27D:#503335
'''

[blended tones face spritesheet 4 times scale, unfolded](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/FaceSpritSheetblend.png)

Now we can represent our little buddy across a variety of tones.  Finally, as an Easter Egg,
we have the original rainbow gradient I used when I first developed the program.  While it
lacks granular control, I think it looks neat.

'''
BitSprite.exe -template=face -upcalse=4 -fold=odd -legacy=true
'''

[Legacy gradient face spritesheet 4 times scale, unfolded](https://github.com/PhiloTFarnsworth/BitSprite/blob/main/docs/FaceSpriteSheetlegacy.png)

It's a candy colored nightmare to be sure, but it can be useful when you're unsure what color you might want
to base your sprites on.  This is created by running across two gradients at .5 lumia on the YCbCr color scheme,
which stands in for our activated bit pixels, while fills and accents are lighter and darker lumias of
the same gradient.

### Flag Commands
After creating the .PNG template and placing it in the Templates folder, the user
can then use the command prompt, to create a template, based on the following flags:

-template (required)
    -This flag looks for its associated string in the templates folder, and if successful opens up the template
    .PNG for parsing.
-fold
    -fold controls whether the template should be reflected across the rightmost bounds.  User has the option to
    choose to fold even, and write the last column twice, or fold odd and have the last column of the template
    only represented once in the output.  Acceptable Values: odd = odd, o; even = even, e. (Not case sensitive) 
-color
    -color designates the color of activated bit pixels, can be expressed as both a single Hex value or two Hex
    values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the
    across the images of the sprite sheet. Acceptable Values: Hex or Hex:Hex (i.e. #FFFFFF, #FFFFFF:#000000).
-accent
    -accent designates the color of accent pixels, can be expressed as both a single Hex value or two Hex
    values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the
    across the images of the sprite sheet. Acceptable Values: Hex or Hex:Hex (i.e. #FFFFFF, #FFFFFF:#000000).
-fill
    -fill designates the color of fill pixels, can be expressed as both a single Hex value or two Hex
    values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the
    across the images of the sprite sheet. Acceptable Values: Hex or Hex:Hex (i.e. #FFFFFF, #FFFFFF:#000000).
-background
    -background designates the color of background pixels, can be expressed as both a single Hex value or two Hex
    values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the
    across the images of the sprite sheet. Acceptable Values: Hex or Hex:Hex (i.e. #FFFFFF, #FFFFFF:#000000).
-outcolor
    -outcolor designates the color of outlines and deactivated bit pixels, can be expressed as both a single Hex value or two Hex
    values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the
    across the images of the sprite sheet. Acceptable Values: Hex or Hex:Hex (i.e. #FFFFFF, #FFFFFF:#000000).
-outline
    -outline can be used to toggle whether BitSprite draws outlines around bit, accent and fill pixels.  Does not 
    effect explicitly designated outline pixels.  Acceptable Values: True = true, t; false = false, f. (Not case
    sensitive, accepts all Golang Bool values.)
-upscale
    -upscale controls the scale of the output images.  Keep in mind that 1 pixel -> 4 -> 9 as you scale in this
    program, so excessively large values will start to chug.  Acceptable Values:  Positive integer (integers < 1
    will automatically be set at 1).
-sheetwidth
    -sheetwidth controls the number of columns in an output Sprite Sheet.  Must be a factor of 256, otherwise
    will default to 16 columns.  Acceptable Values: Positive factor of 256.
-outname
    -outname controls the naming of the output directory and spritesheet.  Acceptable Values: Any string that doesn't
    anger your OS.
-legacy
    -legacy uses the original YCbCr gradient for coloring sprites. Acceptable Values: True = true, t; false = false, f. 
    (Not case sensitive, accepts all Golang Bool values.)

### Why?
BitSprite is a quick and dirty approach to 'art', which allows the user to generate large amounts of assets in a 
short time.  These assets are then stored for easy usage as both a sprite sheet and individual images.  In the
spirit of garage bands that overcame lack of talent with a surplus of noise, so too does BitSprite argue that
quantity has a quality all its own.  Beyond the fact that it has potential for small sprites, icons, ect. there
are also more interesting ways you can composite multiple sprites together to create even more complex and interesting
shapes.

### Further Reading
http://web.archive.org/web/20080228054405/http://www.davebollinger.com/works/pixelrobots/
-This project is inspired by this blog post, which gave some interesting ideas that were adapted either partially
or full cloth.  Still need to adapt the color scheme.

### Todos
Beyond further futzing with colors and how to define them, I think I've cleared the low hanging fruit.  That being
said, please feel free to add issues or pull requests, I would love to get input on how to make even more interesting 
generative art or less cluttered Go code.  