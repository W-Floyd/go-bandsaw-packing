package packing

import (
	"fmt"
)

var Verbose = true

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

	Announce("Pack")

	VarInfo("boundaries", boundaries)
	VarInfo("rectangles", rectangles)
	VarInfo("allowRotation", allowRotation)
	VarInfo("currentSolution", currentSolution)
	VarInfo("offset", offset)

	validRectangles, oversizedRectangles := cullRectangles(boundaries, rectangles, allowRotation)

	Info("")
	VarInfo("validRectangles", validRectangles)
	VarInfo("oversizedRectangles", oversizedRectangles)

	currentSolution.Remainder = append(currentSolution.Remainder, oversizedRectangles...)

	VarInfo("currentSolution", currentSolution)

	Announce("Boundry Range")

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

		Announce("Boundary Instance")
		VarInfo("index", bID)
		VarInfo("boundary", b)
		VarInfo("Boundaries Left", dB)

		Announce("Rectangle Range")

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

			Announce("Rectangle Instance")
			VarInfo("index", rID)
			VarInfo("rectangle ID", rID)
			VarInfo("rectangle", r)
			VarInfo("Rectangles Left", dR)

			var rOptions []Rectangle

			rOptions = append(rOptions, r)

			if allowRotation {
				rOptions = append(rOptions, rotate(r))
			}

			VarInfo("rOptions", rOptions)

			Announce("Rectangle Option Range")

			for _, r := range rOptions {
				Announce("Rectangle Option Instance")
				VarInfo("r", r)
				if checkRectangleFit(r, b.Rectangle) {
					newSolution := currentSolution
					Announce("Rectangle Option Fit")
					placeRectangle(offset, r, &newSolution)
					boundaryOptionSets := provideBoundaryOptions(b, r, offset)
					VarInfo("newSolution", newSolution)
					VarInfo("dR", dR)
					VarInfo("boundaryOptionSets", boundaryOptionSets)

					for _, option := range boundaryOptionSets {
						for _, b := range dB {
							option = append(option, b)
						}
						VarInfo("option", option)
						newSolutions := Pack(option, dR, allowRotation, newSolution, b.Position)
						solutions = append(solutions, newSolutions...)
					}

				} else {
					Info("   Did not fit.")
					VarInfo("b", b)
				}
			}

		}
	}

	if len(validRectangles) == 0 {
		solutions = append(solutions, currentSolution)
	}

	Announce("End Pack")

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

	Announce("Culling")
	VarInfo("rectangles", rectangles)
	VarInfo("boundaries", boundaries)
	Info("")

	for _, r := range rectangles {
		willFit := false
		VarInfo("Rectangle", r)
		for _, b := range boundaries {
			if checkRectangleFit(r, b.Rectangle) {
				Info("   Fits Normal")
				Info("")
				willFit = true
				break
			}
			if allowRotation {
				if checkRectangleFit(rotate(r), b.Rectangle) {
					Info("   Fits Rotated")
					Info("")
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
			Info("   Does Not fit")
			oversizedRectangles = append(oversizedRectangles, r)
		}
	}
	Announce("End Culling")

	return

}

func checkRectangleFit(rectangle Rectangle, boundary Rectangle) bool {

	if rectangle.Height <= boundary.Height && rectangle.Width <= boundary.Width {
		return true
	}
	return false

}

func VarInfo(title string, a ...interface{}) {
	if Verbose {
		fmt.Println(" " + title)
		fmt.Printf("    %+v\n", a...)
	}
}

func Announce(a ...interface{}) {
	if Verbose {
		fmt.Println("")
		fmt.Println(a...)
	}
}

func Info(a ...interface{}) {
	if Verbose {
		fmt.Print(" ")
		fmt.Println(a...)
	}
}

func provideBoundaryOptions(b PlacedRectangle, r Rectangle, offset Position) (boundaryOptionSets [][]PlacedRectangle) {

	if b.Rectangle.Height == r.Height && b.Rectangle.Width == r.Width {

	} else if b.Rectangle.Width == r.Width {

		boundaryOptionSets = [][]PlacedRectangle{
			[]PlacedRectangle{
				PlacedRectangle{
					Rectangle: Rectangle{
						Width:  b.Rectangle.Width,
						Height: b.Rectangle.Height - r.Height,
					},
					Position: Position{
						X: offset.X,
						Y: r.Height + offset.Y,
					},
				},
			},
		}

	} else if b.Rectangle.Height == r.Height {

		boundaryOptionSets = [][]PlacedRectangle{
			[]PlacedRectangle{
				PlacedRectangle{
					Rectangle: Rectangle{
						Height: b.Rectangle.Height,
						Width:  b.Rectangle.Width - r.Width,
					},
					Position: Position{
						Y: offset.Y,
						X: r.Width + offset.X,
					},
				},
			},
		}

	} else {

		boundaryOptionSets = [][]PlacedRectangle{
			[]PlacedRectangle{
				PlacedRectangle{
					Rectangle: Rectangle{
						Width:  b.Rectangle.Width - r.Width,
						Height: b.Rectangle.Height,
					},
					Position: Position{
						X: r.Width + offset.X,
						Y: offset.Y,
					},
				},
				PlacedRectangle{
					Rectangle: Rectangle{
						Width:  r.Width,
						Height: b.Rectangle.Height - r.Height,
					},
					Position: Position{
						X: offset.X,
						Y: r.Height + offset.Y,
					},
				},
			},

			[]PlacedRectangle{
				PlacedRectangle{
					Rectangle: Rectangle{
						Width:  b.Rectangle.Width - r.Width,
						Height: r.Height,
					},
					Position: Position{
						X: r.Width + offset.X,
						Y: offset.Y,
					},
				},
				PlacedRectangle{
					Rectangle: Rectangle{
						Width:  b.Rectangle.Width,
						Height: b.Rectangle.Height - r.Height,
					},
					Position: Position{
						X: offset.X,
						Y: r.Height + offset.Y,
					},
				},
			},
		}

	}

	return
}
