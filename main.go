package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// type Color interface {
// 	RGBA() (r, g, b, a uint32)
// }

// type RGBA struct {
// 	R, G, B, A uint8
// }

// type Image interface {
// 	ColorModel() color.Model
// 	Bounds() Rectangle
// 	At(x, y int) color.Color
// }

// // top-left (min), bottom-right (max) rectangle
// type Rectangle struct {
//     Min, Max Point
// }

// type Point struct {
//     X, Y int
// }

// // image.NRGBA
// // RGBA = Red Green Blue Alpha
// // NRGBA = No-Alpha Red Green Blue
// type NRGBA struct {
// 	Pix []uint8 // byte array contains pixels in image
// 	Stride int // distance between 2 verticallt adjacent pixels
// 	Rect Rectangle // dimension of image
// }

// 1 Load Image From File
func load(filePath string) *image.NRGBA {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot read file: ", err)
	}

	// decode returns: imageValue, imageFormat, error
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Cannot decode file: ", err)
	}

	// Image is interface, NRGBA is actual implementation
	// aka, you can manipulate NRBGA, but NOT Image
	return img.(*image.NRGBA)
}

// 2 Saving Image to File
func saveImage(filePath string, img image.Image) {
	f, _ := os.Create(filePath)
	png.Encode(f, img)
}

func calculateMeanAverageColourWithRect(
	img image.Image,
	rect image.Rectangle,
	useSquaredAverage bool,
) (red, green, blue uint8) {
	var redSum float64
	var greenSum float64
	var blueSum float64

	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y++ {
			pixel := img.At(x, y)
			col := color.RGBAModel.Convert(pixel).(color.RGBA)

			if useSquaredAverage {
				redSum += float64(col.R) * float64(col.R)
				greenSum += float64(col.G) * float64(col.G)
				blueSum += float64(col.B) * float64(col.B)
			} else {
				redSum += float64(col.R)
				greenSum += float64(col.G)
				blueSum += float64(col.B)
			}
		}
	}

	rectArea := float64((rect.Dx() + 1) * (rect.Dy() + 1))

	if useSquaredAverage {
		red = uint8(math.Round(math.Sqrt(redSum / rectArea)))
		green = uint8(math.Round(math.Sqrt(greenSum / rectArea)))
		blue = uint8(math.Round(math.Sqrt(blueSum / rectArea)))
	} else {
		red = uint8(math.Round(redSum / rectArea))
		green = uint8(math.Round(greenSum / rectArea))
		blue = uint8(math.Round(blueSum / rectArea))
	}

	return
}

func pixelate(img image.Image, size int, useSquaredAverage bool) image.Image {
	imgSize := img.Bounds().Size()
	width := imgSize.X
	height := imgSize.Y

	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x += size {
		for y := 0; y < height; y += size {
			rect := image.Rect(x, y, x+size, y+size)

			if rect.Max.X > width {
				rect.Max.X = width
			}
			if rect.Max.Y > height {
				rect.Max.Y = height
			}

			r, g, b := calculateMeanAverageColourWithRect(img, rect, useSquaredAverage)
			col := color.RGBA{r, g, b, 255}

			for x2 := rect.Min.X; x2 < rect.Max.X; x2++ {
				for y2 := rect.Min.Y; y2 < rect.Max.Y; y2++ {
					newImage.Set(x2, y2, col)
				}
			}
		}
	}
	return newImage
}

func printCrochetInstructions(img image.Image, size int) string {

	// start tile = bottom right
	startX := img.Bounds().Max.X
	startY := img.Bounds().Max.Y

	// middle tile = top right
	middleX := img.Bounds().Max.X
	middleY := img.Bounds().Min.Y

	// end tile == top left
	endX := img.Bounds().Min.X
	endY := img.Bounds().Min.Y

	fmt.Printf(
		"start (x,y): (%d, %d), \n"+
			"middle (x, y): (%d, %d), \n"+
			"end (x,y): (%d, %d) \n",
		startX, startY, middleX, middleY, endX, endY)

	startColor := img.At(startX, startY)
	fmt.Printf("Start tile color: ", startColor)

	oddRow := true
	rowLength := 1
	totalRows := startX / size
	// middleRow := totalRows / 2

	var row = Row{}
	var pattern = Pattern{Name: "My Crochet Pattern"}
	c := Coordinate{startX, startY}

	// increase until middle
	for rowLength <= totalRows {
		if oddRow {
			row = increaseOddRow(&c, rowLength, img, size)
		} else {
			row = increaseEvenRow(&c, rowLength, img, size)
		}
		fmt.Printf("Row %d: %v", rowLength, row)
		// fmt.Printf("New x, y: %d, %d", x, y)

		pattern.AddRow(row)
		rowLength += 1
		oddRow = !oddRow
	}

	fmt.Println("Reached middle, Pattern looks like: ", pattern)

	// decrease until end
	colorCounts := pattern.GetColorCounts()
	fmt.Println("Pattern Color Counts: ", colorCounts)
	return "instructions in progress..."

}

// COORDINATE
type Coordinate struct {
	X, Y int
}

// TILE
type Tile struct {
	Coordinate Coordinate
	Color      color.Color
}

func (tile *Tile) AddCoordinate(coordinate Coordinate) {
	tile.Coordinate = coordinate
}
func (tile *Tile) AddColor(color color.Color) {
	tile.Color = color
}

// ROW
type Row struct {
	Tiles []Tile
}

func (row *Row) AddTile(tile Tile) []Tile {
	row.Tiles = append(row.Tiles, tile)
	return row.Tiles
}

// PATTERN
type Pattern struct {
	Rows []Row
	Name string
}

func (p *Pattern) AddRow(r Row) []Row {
	p.Rows = append(p.Rows, r)
	return p.Rows
}
func (p *Pattern) AddName(n string) {
	p.Name = n
}
func (p Pattern) GetColorCounts() map[string]int {
	var colorCounts map[string]int
	colorCounts = make(map[string]int)

	for _, rV := range p.Rows {
		for _, tV := range rV.Tiles {
			R, G, B, _ := tV.Color.RGBA()
			hex := fmt.Sprintf("#%02x%02x%02x", uint8(R), uint8(G), uint8(B))
			colorCounts[hex] += 1
		}
	}
	return colorCounts
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func increaseOddRow(c *Coordinate, rowLength int, img image.Image, size int) (row Row) {

	tiles := []Tile{}
	row = Row{tiles}

	minX := img.Bounds().Min.X
	maxY := img.Bounds().Max.Y

	r := make([]int, rowLength)
	for _ = range r {
		// fmt.Println("i: %d", i)

		// new tile
		color := img.At(c.X, c.Y)
		tile := Tile{}
		tile.AddCoordinate(Coordinate{c.X, c.Y})
		tile.AddColor(color)
		row.AddTile(tile)

		if c.X > minX {
			c.X -= 1 * size
		}
		if c.Y < maxY {
			c.Y += 1 * size
		}
	}
	return
}

func increaseEvenRow(c *Coordinate, rowLength int, img image.Image, size int) (row Row) {
	tiles := []Tile{}
	row = Row{tiles}

	maxX := img.Bounds().Max.X
	minY := img.Bounds().Min.Y

	r := make([]int, rowLength)
	for _ = range r {
		// fmt.Println("i: %d", i)

		// new tile
		color := img.At(c.X, c.Y)
		tile := Tile{}
		tile.AddCoordinate(Coordinate{c.X, c.Y})
		tile.AddColor(color)
		row.AddTile(tile)

		if c.X < maxX {
			c.X += 1 * size
		}
		if c.Y > minY {
			c.Y -= 1 * size
		}
	}
	return
}

func main() {
	// Parse Flags
	pixelSizePtr := flag.Int("p", 120, "more size, more pixelated")
	// patternName := flag.String("n", "My Crochet Pattern", "pattern name")
	flag.Parse()

	// Open image, decode
	file, err := os.Open("pika.png")
	if err != nil {
		log.Fatalln("Problem opening image: ", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalln("Problem decoding image: ", err)
	}

	// Average Image
	// averageColorImage := averageColorImage(img)
	// saveImage("averagePika.png", averageColorImage)

	// Pixelate Image
	pixelatedImage1 := pixelate(img, *pixelSizePtr, false)
	saveImage("pixelPika1.png", pixelatedImage1)
	pixelatedImage2 := pixelate(img, *pixelSizePtr, true)
	saveImage("pixelPika2.png", pixelatedImage2)

	// create instructions
	crochetInstructions := printCrochetInstructions(pixelatedImage1, *pixelSizePtr)
	fmt.Println(crochetInstructions)
}
