package packing

import (
	"fmt"
)

// Rectangle holds a rectangle
type Rectangle struct {
	Width, Height float32
}

// Boundary holds the bounding box that rectangles must fit into
type Boundary struct {
	Dimensions    Rectangle
	Offset        Position
	AllowRotation bool //  In case a material has some lateral property that means rotation would be undesirable.
}

// Position provides x/y coordinates
type Position struct {
	X, Y float32
}

// PlacedRectangle holds information as to the size, position, and rotation of a rectangle
type PlacedRectangle struct {
	Block    Rectangle
	Position Position
}

// Solution holds a single solution to the problem, and any blocks that could not be placed
type Solution struct {
	Blocks    []PlacedRectangle
	Leftovers []Rectangle
}

// rotate will swap a rectangles width/height
func rotate(rectangle Rectangle) (output Rectangle) {
	width := rectangle.Width
	height := rectangle.Height
	output.Height = width
	output.Width = height
	return output
}

// Pack tries to pack a set of given rectangles.
func Pack(boundary Boundary, rectangles []Rectangle) (solutions []Solution) {

	fmt.Println("Pack...")

	solutions = packRecurse([]Boundary{boundary}, rectangles, Solution{}, Position{X: 0, Y: 0})

	return solutions

}

func packRecurse(boundaries []Boundary, rectangles []Rectangle, currentSolution Solution, offset Position) []Solution {

	var solutions []Solution

	// fmt.Println("Recurse...")
	// fmt.Println("Offset: (", offset.X, ",", offset.Y, ")")
	// fmt.Println("")

	// fmt.Println("Culling...")
	cullRectangles(boundaries, &rectangles, &currentSolution)

	for boundaryID, boundary := range boundaries {

		boundariesExtra := boundaries

		boundariesExtra = append(boundariesExtra[:boundaryID], boundariesExtra[boundaryID+1:]...) // Holds the rest of the boundaries

		// fmt.Println("Boundary ", boundaryID)
		// fmt.Println("Width: ", boundary.Dimensions.Width)
		// fmt.Println("Height: ", boundary.Dimensions.Height)
		// fmt.Println(len(boundariesExtra), "left")
		// fmt.Println("")

		for rectangleID, rectangle := range rectangles {

			var rectangleOptions []Rectangle

			rectangleOptions = append(rectangleOptions, rectangle)

			if boundary.AllowRotation {
				rectangleOptions = append(rectangleOptions, rotate(rectangle))
			}

			rectangleExtra := rectangles

			rectangleExtra = append(rectangleExtra[:rectangleID], rectangleExtra[rectangleID+1:]...) // Holds the rest of the rectangles

			for _, rectangleOption := range rectangleOptions {

				newSolution := currentSolution

				fmt.Println("Placing Rectangle")
				placeRectangle(boundary.Offset, rectangleOption, &newSolution)

				// for _, b := range newSolution.Blocks {
				// 	fmt.Println("X-Offset: ", b.Position.X)
				// 	fmt.Println("Y-Offset: ", b.Position.Y)
				// 	fmt.Println("Width: ", b.Block.Width)
				// 	fmt.Println("Height: ", b.Block.Height)
				// 	fmt.Println("")
				// }

				var boundaryOptions [][]Boundary

				newBoundaries := boundariesExtra

				if rectanglesEqual(boundary, rectangleOption) {

					boundaryOptions = [][]Boundary{}

				} else if boundary.Dimensions.Width == rectangleOption.Width {

					boundaryOptions = [][]Boundary{
						[]Boundary{
							Boundary{
								Dimensions: Rectangle{
									Width:  boundary.Dimensions.Width,
									Height: boundary.Dimensions.Height - rectangleOption.Height,
								},
								Offset: Position{
									X: boundary.Offset.X + offset.X,
									Y: boundary.Offset.Y + rectangleOption.Height + offset.Y,
								},
								AllowRotation: boundary.AllowRotation,
							},
						},
					}

				} else if boundary.Dimensions.Height == rectangleOption.Height {

					boundaryOptions = [][]Boundary{
						[]Boundary{
							Boundary{
								Dimensions: Rectangle{
									Height: boundary.Dimensions.Height,
									Width:  boundary.Dimensions.Width - rectangleOption.Width,
								},
								Offset: Position{
									Y: boundary.Offset.Y + offset.Y,
									X: boundary.Offset.X + rectangleOption.Width + offset.X,
								},
								AllowRotation: boundary.AllowRotation,
							},
						},
					}

				} else {

					boundaryOptions = [][]Boundary{
						[]Boundary{
							Boundary{
								Dimensions: Rectangle{
									Width:  boundary.Dimensions.Width - rectangleOption.Width,
									Height: boundary.Dimensions.Height,
								},
								Offset: Position{
									X: boundary.Offset.X + rectangleOption.Width + offset.X,
									Y: boundary.Offset.Y + offset.Y,
								},
								AllowRotation: boundary.AllowRotation,
							},
							Boundary{
								Dimensions: Rectangle{
									Width:  rectangleOption.Width,
									Height: boundary.Dimensions.Height - rectangleOption.Height,
								},
								Offset: Position{
									X: boundary.Offset.X + offset.X,
									Y: boundary.Offset.Y + rectangleOption.Height + offset.Y,
								},
								AllowRotation: boundary.AllowRotation,
							},
						},

						[]Boundary{
							Boundary{
								Dimensions: Rectangle{
									Width:  boundary.Dimensions.Width - rectangleOption.Width,
									Height: rectangleOption.Height,
								},
								Offset: Position{
									X: boundary.Offset.X + rectangleOption.Width + offset.X,
									Y: boundary.Offset.Y + offset.Y,
								},
								AllowRotation: boundary.AllowRotation,
							},
							Boundary{
								Dimensions: Rectangle{
									Width:  boundary.Dimensions.Width,
									Height: boundary.Dimensions.Height - rectangleOption.Height,
								},
								Offset: Position{
									X: boundary.Offset.X + offset.X,
									Y: boundary.Offset.Y + rectangleOption.Height + offset.Y,
								},
								AllowRotation: boundary.AllowRotation,
							},
						},
					}

				}

				for _, boundaryOptionSet := range boundaryOptions {

					for _, boundaryOption := range boundaryOptionSet {

						newBoundaries = append(newBoundaries, boundaryOption)
					}

					for _, solutionOption := range packRecurse(newBoundaries, rectangleExtra, newSolution, boundary.Offset) {

						fmt.Println("Appending solution")

						solutions = append(solutions, solutionOption)

					}

				}

			}

		}

	}

	if len(rectangles) == 0 || len(boundaries) == 0 {
		// fmt.Println("0 Rectangles, adding solution...")
		solutions = append(solutions, currentSolution)
		// fmt.Println("")

	}

	return solutions

}

func rectangleFitsBoundary(boundary Boundary, rectangle Rectangle) bool {

	if rectangle.Height <= boundary.Dimensions.Height && rectangle.Width <= boundary.Dimensions.Width {

		return true
	}

	return false
}

func rectanglesEqual(boundary Boundary, rectangle Rectangle) bool {

	if rectangle == boundary.Dimensions { // If block is equal to the boundary

		return true
	}

	return false
}

func cullRectangles(boundaries []Boundary, rectangles *[]Rectangle, currentSolution *Solution) {

	fmt.Println("Culling")

	for rectangleID, rectangle := range *rectangles {

		willFit := false

		for _, boundary := range boundaries {

			if rectangleFitsBoundary(boundary, rectangle) { // If block is larger than boundary
				fmt.Println("Will fit")
				willFit = true
				break
			}

		}

		if !willFit {

			fmt.Println("No Fit")

			(*rectangles) = append((*rectangles)[:rectangleID], (*rectangles)[rectangleID+1:]...) // Holds the rest of the rectangles

			(*currentSolution).Leftovers = append((*currentSolution).Leftovers, rectangle)
		}

	}

}

func placeRectangle(offset Position, rectangle Rectangle, currentSolution *Solution) {
	(*currentSolution).Blocks = append((*currentSolution).Blocks, PlacedRectangle{Block: rectangle, Position: offset})
}
