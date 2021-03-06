![BitSprite Logo](BitSpriteCard.png)

# BitSprite
## A Free 8-Bit Sprite Generator.

##### Flower Version:
| | | | | | | | | | | | |
|----|----|----|----|----|----|----|----|----|----|----|----|
|![Flower 1](docs/example/Individuals/0.png)|![Flower 2](docs/example/Individuals/1.png)|![Flower 3](docs/example/Individuals/2.png)|![Flower 4](docs/example/Individuals/3.png)|![Flower 5](docs/example/Individuals/4.png)|![Flower 6](docs/example/Individuals/5.png)|![Flower 7](docs/example/Individuals/6.png)|![Flower 8](docs/example/Individuals/7.png)|![Flower 9](docs/example/Individuals/8.png)|![Flower 10](docs/example/Individuals/9.png)|![Flower 11](docs/example/Individuals/10.png)|![Flower 12](docs/example/Individuals/11.png)|


### What?
BitSprite is a program that creates sprite sheets of 256 images based on a single template image.  By taking the index of the variant as it lays on the sprite sheet, Bitsprite uses that index number to activate or deactivate certain designated pixels based on its representation as a binary number. After downloading the package, or the .exe, the user can use the command prompt to run the program, using flags to control the results.  The user has several options to designate colors, as well as 'unfold' the template to create symmetrical images.


### Install
Download the package and place it in an easy to reach place.  You will need to run all commands from inside the Bitsprite directory.  

![Flowers?](docs/FlowersorSkullsHeader.png)

### Using BitSprite

#### Creating a Template
To create a template, open your favorite image editor.  In this example, we'll just use the venerable Microsoft Paint, but any editor should be fine, so long as you encode the template as a png.  Templates have a special color coding to determine how the pixels are read and later written.

#### Pixel Legend
|Pixel|Template Color|
|------|---------------|
|Bit|Black (RGB 0,0,0)|
|Accent|Green (RGB 0,255,0)|
|Fill|Blue (RGB 0,0,255)|
|Outline|Red (RGB 255,0,0)|
|Background|White (RGB 255,255,255)|
|Delimiter|Magenta (RGB 255,0,255)|

Any deviation from these specific colors on a template will result in the offending color being treated as background.  So if your output is entirely transparent, check that your pixels are correctly colored on the template.

'Bit' pixels are how our generator creates the variety between images, as these pixels switch between being active and displaying a set color, or inactive and becoming an outline pixel.  While BitSprite expects 8 bit pixels in a template, you can do less, which results in fewer unique combinations
on the final sprite sheet.  You can also include more than 8 bit pixels in a template, in which case, the bit pattern will repeat, with 9th pixel getting the value
of the first assigned pixel, and so on until it repeats again.  

Accent and fill pixels are static, allowing the user to set pixels that are always active.  

Outline pixels are generated any time an active (colored) pixel is bordered by a background pixel, and by default are colored black.  They can also function as another static pixel if explicitly included in the template, though they will not generate borders around themselves like bit/accent/fill pixels do. 

Background pixels are default transparent.  

#### Back To The Template
Let's look at what a compliant template may look like:

![MS Paint Window with an open template file](docs/ExampleTemplateWorkflow.png)

This example shows a template image at 800% zoom, with each square of color representing one pixel. Notice that there are 5 colors used in the 'face' image, all of which correspond to our control scheme. 

If we run 
```
    BitSprite.exe -template=face
```

we can check the GenerationDirectory and find FaceSpriteSheet.png:

![greyscale face sprite sheet](docs/FaceSpriteSheet.png)

It's a little small.  The transparent background also does no favors for visibility with this basic gray scale. Instead let's return a large image with a better background.We'll set the background to magenta and we'll scale the image by 4, so that our original pixel is written as 4 pixels in a square.

```
    BitSprite.exe -template=face -background=#ff00ff -upscale=4
```

FaceSpriteSheet.png should now look like this:

![greyscale face sprite sheet with magenta background, 4 times scale](docs/FaceSpriteSheetUp.png)

We can now see that our bit pixels are counting up in binary based on their position when read.  Our fill and accent pixels remain the same through each iteration, with distinct colors.  The activated bit pixels also influence the drawing of our outlines, though the red 'smile' from our template stays static.  

#### Now this is all well and good, but why did you use half a face?
Well, a nifty feature of BitSprite is the ability to reflect our templates while writing the new images.

if we decide to run:

```
    BitSprite.exe -template=face -background=#ff00ff -upscale=4 -fold=odd
```

![greyscale face sprite sheet with magenta background, 4 times scale unfolded](docs/FaceSpriteSheetUpscaleEx.png)

Now we have a symmetrical image, and also one that allows your imagination to wander a little bit. Is it the many hairstyles of Tommy Pickles? Maybe it is a happy slime, expanding and contracting across the sprite sheet? You decide!

We can also reflect our images across the bottom of the image as well.  Let's check out the triangle template provided:

```
    Bitsprite.exe -template=triangle -upscale=4
```

![greyscale triangle sprite sheet, 4 times scale](docs/triangleSSvanilla.png)

Now if we want to fold it vertically:

```
    BitSprite.exe -template=triangle -upscale=4 -vertfold=odd
```

![greyscale triangle sprite sheet, 4 times scale unfolded vertically](docs/triangleSSvert.png)

Notice the vertfold uses the same inputs as fold, and even better, we can combine the two:

```
    BitSprite.exe -template=triangle -upscale=4 -vertfold=even -fold=even
```

![greyscale triangle sprite sheet, 4 times scale unfolded vertically and horizontally](docs/triangleSSvertEfoldE.png)


#### It's still a little bland
Fine, BitSprite also supports some cool ways to control colors.

```
    BitSprite.exe -template=face -upscale=4 -fold=odd -color=#9a3300 -accent=#3d671d -fill=#F1C27D
```

![colored in face sprite sheet 4 times scale, unfolded](docs/FaceSpriteSheetColor.png)

With color, we create a passable face.  But we're not done yet, we can also do a blend of tones across the entire sheet, for any of those values.

```
    BitSprite.exe -template=face -upscale=4 -fold=odd -color=#9a3300:#4f1a00 -accent=#3d671d:#497665 -fill=#F1C27D:#503335
```

![blended tones face sprite sheet 4 times scale, unfolded](docs/FaceSpriteSheetblend.png)

Now we can represent our little buddy across a variety of tones.  Finally, as an Easter Egg, we have the original rainbow gradient I used when I first developed the program. While it lacks granular control, I think it looks neat.

```
    BitSprite.exe -template=face -upscale=4 -fold=odd -legacy=true
```

![Legacy gradient face sprite sheet 4 times scale, unfolded](docs/FaceSpriteSheetlegacy.png)

It's a candy colored nightmare to be sure, but it can be useful when you're unsure what color you might want to base your sprites on. This is created by running across two gradients at .5 lumia on the YCbCr color scheme, which stands in for our activated bit pixels, while fills and accents are lighter and darker lumias of the same gradient.

![Robot heads](docs/RobotsHeader.png)

#### But isn't 8 pixels a little limiting?

Yes, yes it is.  However, this limit is also a strength;  Let's consider creating a flower.  We could probably create a passable flower using a single template, but if we break down the flower into parts and then sew them together, we can create huge amount of diversity in a specific asset in a short time.  We can use BitSprite to preview what this would look like.

![An example of 4 templates representing a flower, a stem, leaves and a base of plant](docs/FlowerIntro.png)

We can combine these into one template and run it through BitSprite.

![The 4 templates combined together into a composite template, no delimiters](docs/FlowerUnlimited.png)


```
    BitSprite.exe -template=flower -upscale=4 -fold=o -legacy=t
```

This results in something a little interesting, but something feels a little off:

![the composite template from the previous image rendered, 256 "flowers"](docs/FlowerUnlimitedOutput.png)

Some are definitely more plant-like than the others, but this is only a single combination of 256^4, and it's rather unlikely if going to want 0th base image, 0th leaves image, 0th stem and 0th flower.  Instead, when using this in a game engine you might pick random numbers for all 4 images and put them together.  So what do we do to the template to return 256 samples?  

![Enter the magenta delimiter pixel; the composite template now has pink delimiter pixels at the beginning of each image](docs/FlowerDelimiter.png)

First, we add our magenta 'delimiter' pixels to what would be the first pixel of each sub-template (reading rows from left to right).  Let's 
run that command again


```
    BitSprite.exe -template=flowerdelimited -upscale=4 -fold=o -legacy=t
```


![Sprite Sheet with 256 randomized flowers, with all sub templates randomly rendered](docs/FlowerDelimiterOutput.png)

Now we have a cross sample of what we can expect from randomized rendering of individual parts.  Now this feature isn't very useful for production, but when you're prototyping a composite sprite, like the flower above, this can give you a good idea of whether your shapes work together.

#### A Final Note
You might be wondering, what if I'm not using templates with 8 'Bit' pixels?  You'll find the 'Bit' pattern repeats every 8 'Bit' pixels you have in your template.  There's no upper bound for 'Bit' pixels, but large images with more complexity generally don't look great.

![Triangles](docs/TriangleHeader.png)

### Flag Commands
After creating the .PNG template and placing it in the Templates folder, the user can then use the command prompt, to create a sprite sheet, based on the following flags:
```
-template (required)    Expected Values: Any string that doesn't anger your OS.
```
Template looks in the templates folder for a file named after the provided string, and if successful opens up the template .PNG for parsing.

```
-fold   Expected Values: odd = odd, o; even = even, e. (Not case sensitive) 
```
Fold controls whether the template should be reflected across the rightmost bounds.  User has the option to choose to fold even, and write the last column twice or fold odd and have the last column of the template only represented once in the output.  

```
-vertfold   Expected Values: odd = odd, o; even = even, e. (Not case sensitive)
```
Vertfold controls whether the template should be reflected across the bottom bounds.  User has the option to choose to fold even, and write the last row twice or fold odd and have the last row of the template only represented once in the output.  
```
-color   Expected Values: Hex or Hex:Hex (#FFFFFF or #FFFFFF:#000000).
```
Color designates the color of activated bit pixels, can be expressed as both a single Hex value or two Hex values with a ':' in between.Passing two Hex values will result in 'blended' shades between the designated colors across the images of the sprite sheet. 

```
-accent   Expected Values: Hex or Hex:Hex (#FFFFFF or #FFFFFF:#000000).
```
Accent designates the color of accent pixels, can be expressed as both a single Hex value or two Hex values with a ':' in between.Passing two Hex values will result in 'blended' shades between the designated colors across the images of the sprite sheet. 
```
-fill    Expected Values: Hex or Hex:Hex (#FFFFFF or #FFFFFF:#000000).
```
Fill designates the color of fill pixels, can be expressed as both a single Hex value or two Hex values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the designated colors across the images of the sprite sheet. 
```
-background    Expected Values: Hex or Hex:Hex (IE #FFFFFF, #FFFFFF:#000000).
```
Background designates the color of background pixels, can be expressed as both a single Hex value or two Hex values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the designated colors across the images of the sprite sheet. 
```
-outcolor    Expected Values: Hex or Hex:Hex (IE #FFFFFF,#FFFFFF:#000000).
```
Outcolor designates the color of outlines and deactivated bit pixels, can be expressed as both a single Hex value or two Hex values with a ':' in between.  Passing two Hex values will result in 'blended' shades between the designated colors across the images of the sprite sheet. 
```
-outline    Expected Values: True = true, t; false = false, f. (Not case sensitive, accepts all Golang Bool values.)
```
Outline can be used to toggle whether BitSprite draws outlines around bit, accent and fill pixels.  Does not effect explicitly designated outline pixels.  
```
-upscale    Expected Values:  Positive integer (integers < 1 will automatically be set at 1).
```
Upscale controls the scale of the output images.  Keep in mind that 1 pixel -> 4 -> 9 as you scale in this program.
```
-sheetwidth    Expected Values: Positive factor integer of 256.
```
Sheetwidth controls the number of columns in an output Sprite Sheet.  Must be a factor of 256, otherwise will default to 16 columns.  
```
-outname    Expected Values: Any string that doesn't anger your OS.
```
Outname controls the naming of the output directory and sprite sheet.  'docs' is reserved and will output to docs/example.
```
-individuals    Expected Values: True = true, t; false = false, f. (Not case sensitive, accepts all Golang Bool values.
```
Indivduals controls whether BitSprite creates a directory with individual .png files of all the images that constitute the spritesheet.  Defaults
to false. 
```
-legacy    Expected Values: True = true, t; false = false, f. (Not case sensitive, accepts all Golang Bool values.)
```
Legacy uses the original YCbCr gradient for coloring sprites. 

![Dog with hat](docs/DogwHatHeader.png)

### Why?
BitSprite is a quick and dirty approach to 'art', which allows the user to generate large amounts of assets in a short time.  These assets are then stored for easy usage as both a sprite sheet and individual images.  In the spirit of garage bands that overcame lack of talent with a surplus of noise, so too does BitSprite argue that quantity has a quality all its own.  Beyond the fact that it has potential for small sprites, icons, etc. there are also more interesting ways you can composite multiple sprites together to create even more complex and interesting shapes.

### FAQs
#### All my pixels are coming out black, why do?
This is usually a problem with passing bad color values to any of the color flags.  Please ensure it is a valid Hex Code and begins with #.

#### When I produce a new sprite sheet, it's the wrong shape/completely transparent!
This problem usually indicates that you are not using the correct colors in the template file.  Ensure that they are as indicated in the pixel legend at the beginning of the readme.    

### Further Reading
[Pixel Robots](http://web.archive.org/web/20080228054405/http://www.davebollinger.com/works/pixelrobots/) by Dave Bollinger

This project is inspired by this blog post, which I think is a very interesting dive into generated art and some nifty things you can do with it.

### Todos
Beyond further futzing with colors and how to define them, I think I've cleared the low hanging fruit.  That being said, please feel free to add issues or pull requests,I would love to get input on how to make even more interesting generative art or less cluttered Go code. 

| | | | | | | | | | | | |
|----|----|----|----|----|----|----|----|----|----|----|----|
|![Flower 245](docs/example/Individuals/244.png)|![Flower 246](docs/example/Individuals/245.png)|![Flower 247](docs/example/Individuals/246.png)|![Flower 248](docs/example/Individuals/247.png)|![Flower 249](docs/example/Individuals/248.png)|![Flower 250](docs/example/Individuals/249.png)|![Flower 251](docs/example/Individuals/250.png)|![Flower 252](docs/example/Individuals/251.png)|![Flower 253](docs/example/Individuals/252.png)|![Flower 254](docs/example/Individuals/253.png)|![Flower 255](docs/example/Individuals/254.png)|![Flower 256](docs/example/Individuals/255.png)| 
