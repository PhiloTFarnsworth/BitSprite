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
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSVanilla.png"
	testArgs := []string{"cmd", "-template=triangle"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOddFold(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSOdd.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=odd"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestEvenFold(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSEven.png"
	testArgs := []string{"cmd", "-template=triangle", "-fold=even"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestScale(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSScale.png"
	testArgs := []string{"cmd", "-template=triangle", "-upscale=4", "-fold=none"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColor(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSRed.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-upscale=1"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsBA(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBA.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-accent=#00ff00"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsBF(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBF.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-accent=", "-fill=#0000ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsFA(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSFA.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=", "-accent=#00ff00", "-fill=#0000ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorsBAF(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBAF.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000", "-accent=#00ff00", "-fill=#0000ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestBlends(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBlend.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=#ff0000:#00ff00", "-accent=#00ff00:#0000ff", "-fill=#0000ff:#ff0000"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestColorOutline(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSOutRed.png"
	testArgs := []string{"cmd", "-template=triangle", "-color=", "-accent=", "-fill=", "-outcolor=#ff0000"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestOutlineBool(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSOutFalse.png"
	testArgs := []string{"cmd", "-template=triangle", "-outline=f", "-outcolor="}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestBackgroundColor(t *testing.T) {
	gotFileName := "GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "testResources/TriangleSSBack.png"
	testArgs := []string{"cmd", "-template=triangle", "-outline=t", "-background=#ff00ff"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func TestLegacy(t *testing.T) {
	gotFileName := "/GenerationDirectory/Triangle/TriangleSpriteSheet.png"
	wantFileName := "/testResources/TriangleSSLegacy.png"
	testArgs := []string{"cmd", "-template=triangle", "-background=", "-legacy=t"}
	Compare(t, gotFileName, wantFileName, testArgs)
}

func Compare(t *testing.T, gotFileName, wantFileName string, testArgs []string) {
	os.Args = testArgs
	for i := range os.Args {
		fmt.Printf(os.Args[i])
		fmt.Printf("\n")
	}

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
