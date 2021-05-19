package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

//We're going to compare with fresh images with ones previously prepared.  One annoying thing
//is that flags appear to persist between test functions, so there are redundant Args.
func TestDefault(t *testing.T) {
	fmt.Printf("TestDefault\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSVanilla.png"
	testArgs := []string{"cmd", "-template=triangle"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestIndividuals(t *testing.T) {
	fmt.Printf("TestIndividuals\n")
	gotFileName := "GenerationDirectory/Triangle/Individuals/127.png"
	wantFileName := "testResources/127Test.png"
	testArgs := []string{"cmd", "-template=triangle", "-individuals=t"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOddFold(t *testing.T) {
	fmt.Printf("TestOddFold\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSOdd.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=odd", "-individuals=f"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestEvenFold(t *testing.T) {
	fmt.Printf("TestEvenFold\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSEven.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=even"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestEvenVert(t *testing.T) {
	fmt.Printf("TestEvenVert\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSVE.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=", "-vertfold=even"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOddVert(t *testing.T) {
	fmt.Printf("TestOddVert\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSVO.png"
	testArgs := []string{"cmd", "-template=triangle", "-vertfold=odd"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestEvenFolds(t *testing.T) {
	fmt.Printf("TestEvenFolds\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSFEVE.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=even", "-vertfold=even"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOddFolds(t *testing.T) {
	fmt.Printf("TestOddFolds\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSFOVO.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=odd", "-vertfold=odd"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOddFoldEvenVert(t *testing.T) {
	fmt.Printf("TestOddFoldEvenVert\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSFOVE.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=odd", "-vertfold=even"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestEvenFoldOddVert(t *testing.T) {
	fmt.Printf("TestEvenFoldOddVert\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSFEVO.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=even", "-vertfold=odd"}
	Compare(t, gotFileName, wantFileName, testArgs)
}
func TestScale(t *testing.T) {
	fmt.Printf("TestScale\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSScale.png"
	testArgs := []string{"cmd", "-template=triangle", "-upscale=4", "-fold=", "-vertfold="}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColor(t *testing.T) {
	fmt.Printf("TestColor\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSRed.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-upscale=1"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsBA(t *testing.T) {
	fmt.Printf("TestColorsBA\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBA.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-accent=#00ff00"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsBF(t *testing.T) {
	fmt.Printf("TestColorsBF\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBF.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-accent=", "-fill=#0000ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsFA(t *testing.T) {
	fmt.Printf("TestColorsFA\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSFA.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=", "-accent=#00ff00", "-fill=#0000ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsBAF(t *testing.T) {
	fmt.Printf("TestColorsBAF\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBAF.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-accent=#00ff00", "-fill=#0000ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestBlends(t *testing.T) {
	fmt.Printf("TestBlends\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBlend.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000:#00ff00", "-accent=#00ff00:#0000ff", "-fill=#0000ff:#ff0000"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorOutline(t *testing.T) {
	fmt.Printf("TestColorOutline\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSOutRed.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=", "-accent=", "-fill=", "-outcolor=#ff0000"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOutlineBool(t *testing.T) {
	fmt.Printf("TestOutlineBool\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSOutFalse.png"
	testArgs := []string{"cmd", "-template=triangle", "-outline=f", "-outcolor="}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestBackgroundColor(t *testing.T) {
	fmt.Printf("TestBackgroundColor\n")
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBack.png"
	testArgs := []string{"cmd", "-template=triangle", "-outline=t", "-background=#ff00ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestLegacy(t *testing.T) {
	fmt.Printf("TestLegacy\n")
	gotFileName := "/GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "/testResources/TriangleSSLegacy.png"
	testArgs := []string{"cmd", "-template=triangle", "-background=", "-legacy=t"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestUnlimitedComposite(t *testing.T) {
	fmt.Printf("TestUnlimitedComposite\n")
	gotFileName := "/GenerationDirectory/Flower/FlowerSpriteSheet.png"
	wantFileName := "/testResources/VanillaFlower.png"
	testArgs := []string{"cmd", "-template=flower", "-fold=o", "-legacy=t"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestDelimitedComposite(t *testing.T) {
	fmt.Printf("TestDelimitedComposite\n")
	gotFileName := "/GenerationDirectory/Flower/FlowerSpriteSheet.png"
	wantFileName := "/testResources/DelimitedFlower.png"
	testArgs := []string{"cmd", "-template=flowerDelimited", "-fold=o", "-legacy=t", "-outname=Flower", "-randseed=f"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestFace(t *testing.T) {
	fmt.Printf("TestFace\n")
	gotFileName := "/GenerationDirectory/face/faceSpriteSheet.png"
	wantFileName := "/testResources/faceSS.png"
	testArgs := []string{"cmd", "-template=face", "-legacy=f", "-color=#ff0000", "-accent=#00ff00", "-fill=#0000ff", "-fold=odd", "-outname=", "randseed=t"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

//This also tests the reading of red template pixels (outlines), which I forgot to consider.  We'll
//use the example face.png template to have that included.  Do this test last otherwise you need to reset
//all the flags set here.
func TestInputSanitizers(t *testing.T) {
	fmt.Printf("TestInputSanitizers\n")
	gotFileName := "/GenerationDirectory/BadInput/BadInputSpriteSheet.png"
	wantFileName := "/testResources/faceSSbadinput.png"
	testArgs := []string{"cmd", "-template=face", "-legacy=f", "-color=badInput", "-accent=BadInput", "-fill=BadInput", "-sheetwidth=1000", "-upscale=0", "-outname=BadInput", "-fold=none"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func BenchmarkDefault(b *testing.B) {
	os.Args = []string{"cmd", "-template=triangle", "-color=", "-accent=", "-fill=", "-sheetwidth=16", "-upscale=1", "-outname="}
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkIndividuals(b *testing.B) {
	os.Args = []string{"cmd", "-template=triangle", "individuals=t"}
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkTenScale(b *testing.B) {
	os.Args = []string{"cmd", "-template=triangle", "-upscale=10", "individuals=f"}
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkBlends(b *testing.B) {
	os.Args = []string{"cmd", "-template=triangle", "-upscale=1", "-color=#ff0000:#00ff00", "-accent=#00ff00:#0000ff", "-fill=#0000ff:#ff0000"}
	for i := 0; i < b.N; i++ {
		main()
	}
}

func BenchmarkLegacy(b *testing.B) {
	os.Args = []string{"cmd", "-template=triangle", "-upscale=1", "-color=#", "-accent=", "-fill="}
	for i := 0; i < b.N; i++ {
		main()
	}
}

func Compare(t *testing.T, gotFileName, wantFileName string, testArgs []string) {
	os.Args = testArgs

	main()

	currentDir, err := filepath.Abs("")
	if err != nil {
		t.Fatal(err)
	}

	gotFile, err := os.Open(filepath.Join(currentDir + "/" + gotFileName))
	if err != nil {
		t.Fatal(err)
	}

	wantFile, err := os.Open(filepath.Join(currentDir + "/" + wantFileName))
	if err != nil {
		t.Fatal(err)
	}

	gotStream, err := png.Decode(gotFile)
	if err != nil {
		t.Fatal(err)
	}
	wantStream, err := png.Decode(wantFile)
	if err != nil {
		t.Fatal(err)
	}
	gotFile.Seek(0, 0)
	wantFile.Seek(0, 0)
	gotConfig, err := png.DecodeConfig(gotFile)
	if err != nil {
		t.Fatal(err)
	}
	wantConfig, err := png.DecodeConfig(wantFile)
	if err != nil {
		t.Fatal(err)
	}

	// Dimension Check
	if gotConfig.Width != wantConfig.Width {
		t.Fatalf("Got width of %v, wanted width of %v", gotConfig.Width, wantConfig.Width)
	}
	if gotConfig.Height != wantConfig.Height {
		t.Fatalf("Got height of %v, wanted height of %v", gotConfig.Height, wantConfig.Height)
	}

	//Compare the streams.
	for i := 0; i < wantConfig.Width*wantConfig.Height; i++ {
		//translate our increment to (x,y)
		p := image.Point{(i % wantConfig.Width), int(i / wantConfig.Width)}
		if wantStream.At(p.X, p.Y) != gotStream.At(p.X, p.Y) {
			t.Fatalf("Wanted color %v at point %v,%v; got color %v", wantStream.At(p.X, p.Y), p.X, p.Y, gotStream.At(p.X, p.Y))
		}
	}

	wantFile.Close()
	gotFile.Close()
}
