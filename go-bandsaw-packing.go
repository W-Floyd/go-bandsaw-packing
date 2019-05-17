package packing

import (
	"fmt"
)

var verbose = true

type Rectangle struct {
	Width  int
	Height int
}

type Position struct {
	X int
	Y int
}

type PlacedRectangle struct { // doubles as a boundary
	Rectangle Rectangle
	Position  Position
}

type Solution struct {
	Remainder []Rectangle
	Placed    []PlacedRectangle
}

func Pack(boundaries []PlacedRectangle, rectangles []Rectangle, allowRotation bool, currentSolution Solution, offset Position) (solutions []Solution) {

	announce("Pack")

	varInfo("boundaries", boundaries)
	varInfo("rectangles", rectangles)
	varInfo("allowRotation", allowRotation)
	varInfo("currentSolution", currentSolution)
	varInfo("offset", offset)

	validRectangles, oversizedRectangles := cullRectangles(boundaries, rectangles, allowRotation)

	info("")
	varInfo("validRectangles", validRectangles)
	varInfo("oversizedRectangles", oversizedRectangles)

	currentSolution.Remainder = append(currentSolution.Remainder, oversizedRectangles...)

	varInfo("currentSolution", currentSolution)

	announce("Boundry Range")

	for bID, b := range boundaries {

		excluded := false
		dB := make([]PlacedRectangle, len(boundaries))
		k := 0
		for _, n := range boundaries {
			if n != b || excluded { // filter
				dB[k] = n
				k++
			} else {
				excluded = true
			}
		}
		dB = dB[:k]

		announce("Boundary Instance")
		varInfo("index", bID)
		varInfo("boundary", b)
		varInfo("Boundaries Left", dB)

		announce("Rectangle Range")

		for rID, r := range validRectangles {

			excluded := false
			dR := make([]Rectangle, len(validRectangles))
			k := 0
			for _, n := range validRectangles {
				if n != r || excluded { // filter
					dR[k] = n
					k++
				} else {
					excluded = true
				}
			}
			dR = dR[:k]

			announce("Rectangle Instance")
			varInfo("index", rID)
			varInfo("rectangle ID", rID)
			varInfo("rectangle", r)
			varInfo("Rectangles Left", dR)

			var rOptions []Rectangle

			rOptions = append(rOptions, r)

			if allowRotation {
				rOptions = append(rOptions, rotate(r))
			}

			varInfo("rOptions", rOptions)

			announce("Rectangle Option Range")

			for _, r := range rOptions {
				announce("Rectangle Option Instance")
				if checkRectangleFit(r, b.Rectangle) {
					newSolution := currentSolution
					announce("Rectangle Option Fit")
					placeRectangle(offset, r, &newSolution)

				}
			}

		}
	}

	announce("End Pack")

	return solutions

}

// rotate will swap a rectangles width/height
func rotate(rectangle Rectangle) (output Rectangle) {
	output.Height = rectangle.Width
	output.Width = rectangle.Height
	return output
}

func placeRectangle(offset Position, rectangle Rectangle, currentSolution *Solution) {
	fmt.Println("Placing: ", rectangle)
	fmt.Println("Position: ", offset)

	(*currentSolution).Placed = append((*currentSolution).Placed, PlacedRectangle{Rectangle: rectangle, Position: offset})
}

func cullRectangles(boundaries []PlacedRectangle, rectangles []Rectangle, allowRotation bool) (validRectangles []Rectangle, oversizedRectangles []Rectangle) {

	announce("Culling")
	varInfo("rectangles", rectangles)
	varInfo("boundaries", boundaries)
	info("")

	for _, r := range rectangles {
		willFit := false
		varInfo("Rectangle", r)
		for _, b := range boundaries {
			if checkRectangleFit(r, b.Rectangle) {
				info("   Fits Normal")
				info("")
				willFit = true
				break
			}
			if allowRotation {
				if checkRectangleFit(rotate(r), b.Rectangle) {
					info("   Fits Rotated")
					info("")
					willFit = true
					break
				}
			}
			if willFit {

				willFit = true
				break
			}
		}
		if willFit {
			validRectangles = append(validRectangles, r)
		} else {
			info("   Does Not fit")
			oversizedRectangles = append(oversizedRectangles, r)
		}
	}
	announce("End Culling")

	return

}

func checkRectangleFit(rectangle Rectangle, boundary Rectangle) bool {

	if rectangle.Height <= boundary.Height && rectangle.Width <= boundary.Width {
		return true
	}
	return false

}

func varInfo(title string, a ...interface{}) {
	if verbose {
		fmt.Println(" " + title)
		fmt.Printf("    %+v\n", a...)
	}
}

func announce(a ...interface{}) {
	if verbose {
		fmt.Println("")
		fmt.Println(a...)
	}
}

func info(a ...interface{}) {
	if verbose {
		fmt.Print(" ")
		fmt.Println(a...)
	}
}
