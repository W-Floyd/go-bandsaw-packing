package packing

import (
	"fmt"
	"reflect"
	"sort"
)

var iteration = 0

var Verbose = false
var forceVerbose = false

var globalState []State

type State struct {
	Boundaries []Rectangle
	Rectangles []RectangleQuantity
	Iteration  int
}
type Rectangle struct {
	Width  int
	Height int
}

type RectangleQuantity struct {
	Rectangle Rectangle
	Quantity  int
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
	Remainder []RectangleQuantity
	Placed    []PlacedRectangle
}

func areaUsed(solution Solution) (area int) {
	for _, r := range solution.Placed {
		area += r.Rectangle.Width * r.Rectangle.Height
	}
	return
}

func FilteredPack(boundaries []PlacedRectangle, rectangles []RectangleQuantity, allowRotation bool, currentSolution Solution, offset Position) (output []Solution) {

	allSolutions := Pack(boundaries, rectangles, allowRotation, currentSolution, offset)

	options := map[int]Solution{}

	for _, s := range allSolutions {
		options[areaUsed(s)] = s
	}

	keys := make([]int, len(options))

	i := 0
	for k := range options {
		keys[i] = k
		i++
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	for _, a := range keys {
		output = append(output, options[a])
	}

	return output

}

func normalizeRectangles(rectangles []PlacedRectangle, allowRotation bool) (output []Rectangle) {
	for _, r := range rectangles {
		if r.Rectangle.Height > r.Rectangle.Width && allowRotation {
			output = append(output, Rectangle{Width: r.Rectangle.Height, Height: r.Rectangle.Width})
		} else {
			output = append(output, r.Rectangle)
		}
	}

	// sort each Parent in the parents slice by Id
	sort.Slice(output, func(i, j int) bool {
		if output[i].Width == output[j].Width {
			return output[i].Height < output[j].Height
		}
		return output[i].Width < output[j].Width
	})

	return output
}

func stateDuplicate(givenState State) bool {

	for _, state := range globalState {

		if reflect.DeepEqual(state.Boundaries, givenState.Boundaries) && reflect.DeepEqual(state.Rectangles, givenState.Rectangles) {
			Verbose = true
			Announce("Equal")
			VarInfo("state", state)
			VarInfo("givenState", givenState)
			Info("")
			Verbose = false
			return true
		}
	}
	return false
}

func Pack(boundaries []PlacedRectangle, rectangles []RectangleQuantity, allowRotation bool, currentSolution Solution, offset Position) (solutions []Solution) {

	currentState := State{
		Rectangles: rectangles,
		Boundaries: normalizeRectangles(boundaries, allowRotation),
		Iteration:  iteration,
	}

	Verbose = true
	Info("")
	VarInfo("iteration", iteration)
	VarInfo("currentState", currentState)
	VarInfo("currentSolution", currentSolution)

	Verbose = false

	iteration++

	validRectangles, oversizedRectangles := cullRectangles(boundaries, rectangles, allowRotation)

	Verbose = true

	Info("foo")

	if len(validRectangles) == 0 || len(boundaries) == 0 {
		solutions = append(solutions, currentSolution)
		Announce("Final Solution")
		VarInfo("currentSolution", currentSolution)
		return solutions
	}

	Verbose = false

	skipPack := false

	if stateDuplicate(currentState) {
		skipPack = true
	}

	globalState = append(globalState, currentState)

	if skipPack {
		Verbose = true
		Announce("Skipping")
		Verbose = false

	} else {

		Verbose = true
		Announce("Pack")

		VarInfo("boundaries", boundaries)
		VarInfo("rectangles", rectangles)
		VarInfo("allowRotation", allowRotation)
		VarInfo("currentSolution", currentSolution)
		Info("placed: ", len(currentSolution.Placed))
		VarInfo("offset", offset)

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
				dR := make([]RectangleQuantity, len(validRectangles))
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

				if r.Quantity > 1 {
					dR = append(dR, RectangleQuantity{Rectangle: r.Rectangle, Quantity: r.Quantity - 1})
				}

				Announce("Rectangle Instance")
				VarInfo("index", rID)
				VarInfo("rectangle ID", rID)
				VarInfo("rectangle", r)
				VarInfo("Rectangles Left", dR)

				var rOptions []Rectangle

				rOptions = append(rOptions, r.Rectangle)

				if allowRotation && !reflect.DeepEqual(rotate(r.Rectangle), r.Rectangle) {
					rOptions = append(rOptions, rotate(r.Rectangle))
				}

				VarInfo("rOptions", rOptions)

				Announce("Rectangle Option Range")

				for _, r := range rOptions {
					Announce("Rectangle Option Instance")
					VarInfo("r", r)

					if checkRectangleFit(r, b.Rectangle) {
						newSolution := currentSolution
						Announce("Rectangle Option Fit")
						placeRectangle(b.Position, r, &newSolution)
						boundaryOptionSets := provideBoundaryOptions(b, r, b.Position)
						VarInfo("newSolution", newSolution)
						VarInfo("dR", dR)
						VarInfo("boundaryOptionSets", boundaryOptionSets)
						if len(boundaryOptionSets) == 0 {
							boundaryOptionSets = append(boundaryOptionSets, dB)
						}

						for _, option := range boundaryOptionSets {
							for _, b := range dB {
								option = append(option, b)
							}
							VarInfo("option", option)
							Verbose = false
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

		Announce("End Pack")

	}

	return solutions

}

// rotate will swap a rectangles width/height
func rotate(rectangle Rectangle) (output Rectangle) {
	output.Height = rectangle.Width
	output.Width = rectangle.Height
	return output
}

func placeRectangle(offset Position, rectangle Rectangle, currentSolution *Solution) {
	Info("Placing: ", rectangle)
	Info("Position: ", offset)

	(*currentSolution).Placed = append((*currentSolution).Placed, PlacedRectangle{Rectangle: rectangle, Position: offset})
}

func cullRectangles(boundaries []PlacedRectangle, rectangles []RectangleQuantity, allowRotation bool) (validRectangles []RectangleQuantity, oversizedRectangles []RectangleQuantity) {

	Announce("Culling")
	VarInfo("rectangles", rectangles)
	VarInfo("boundaries", boundaries)
	Info("")

	if len(boundaries) == 0 {
		Info("   No Boundaries")
		return []RectangleQuantity{}, rectangles
	}

	for _, r := range rectangles {
		willFit := false
		VarInfo("Rectangle", r)
		for _, b := range boundaries {
			if r.Quantity > 0 {
				if checkRectangleFit(r.Rectangle, b.Rectangle) {
					Info("   Fits Normal")
					Info("")
					willFit = true
					break
				}
				if allowRotation {
					if checkRectangleFit(rotate(r.Rectangle), b.Rectangle) {
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
	if Verbose && forceVerbose {
		fmt.Println(" " + title)
		fmt.Printf("    %+v\n", a...)
	}
}

func Announce(a ...interface{}) {
	if Verbose && forceVerbose {
		fmt.Println("")
		fmt.Println(a...)
	}
}

func Info(a ...interface{}) {
	if Verbose && forceVerbose {
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
