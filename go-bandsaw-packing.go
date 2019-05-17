package packing

import "fmt"

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

	fmt.Println("")
	fmt.Println("Pack")

	for _, b := range boundaries {
		fmt.Println("Boundary:")
		rectangleInfo(b.Rectangle)
		fmt.Println("")
	}

	for _, r := range rectangles {
		fmt.Println("Rectangle:")
		rectangleInfo(r)
		fmt.Println("")
	}

	validRectangles, oversizedRectangles := cullRectangles(boundaries, rectangles)

	currentSolution.Remainder = append(currentSolution.Remainder, oversizedRectangles...)

	if len(validRectangles) == 0 {
		// fmt.Println("0 valid rectangles left")
		solutions = append(solutions, currentSolution)
	} else {
		for rID, r := range validRectangles {
			rOptions := []Rectangle{r}

			if allowRotation {
				rOptions = append(rOptions, rotate(r))
			}

			newRectangles := validRectangles
			newRectangles = append(newRectangles[:rID], newRectangles[rID+1:]...) // Holds the rest of the rectangles

			fmt.Println("rOptions: ", len(rOptions))

			for _, r := range rOptions {

				for bID, b := range boundaries {

					var boundaryOptionSets [][]PlacedRectangle

					newBoundaries := boundaries
					newBoundaries = append(newBoundaries[:bID], newBoundaries[bID+1:]...) // Holds the rest of the boundaries

					fmt.Println("newBoundaries: ", len(newBoundaries))

					if b.Rectangle.Height == r.Height && b.Rectangle.Width == r.Width {

						if len(newBoundaries) > 0 {
							boundaryOptionSets = [][]PlacedRectangle{
								newBoundaries,
							}
						} else {
							fmt.Println("Exact fit", offset)

							placeRectangle(b.Position, r, &currentSolution)
							solutions = append(solutions, currentSolution)
						}

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
							newBoundaries,
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
							newBoundaries,
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
							newBoundaries,
						}

					}

					fmt.Println("boundaryOptionSets: ", len(boundaryOptionSets))

					for _, bOptionSet := range boundaryOptionSets {

						fmt.Println("bOptionSet: ", len(bOptionSet))

						for _, b := range bOptionSet {

							newSolution := currentSolution

							var newSolutions []Solution

							if checkRectangleFit(r, b.Rectangle) {

								fmt.Println("Fit - Recurse")

								placeRectangle(offset, r, &newSolution)

								newSolutions = Pack(bOptionSet, newRectangles, allowRotation, newSolution, Position{X: b.Position.X, Y: b.Position.Y})

							} else {
								fmt.Println("No fit - Rectangle")
								rectangleInfo(r)
								fmt.Println("Boundary")
								rectangleInfo(b.Rectangle)
								fmt.Println("")
							}

							solutions = append(solutions, newSolutions...)

						}

					}
				}
			}
		}
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
	fmt.Println("Placing: ", rectangle)
	fmt.Println("Position: ", offset)

	(*currentSolution).Placed = append((*currentSolution).Placed, PlacedRectangle{Rectangle: rectangle, Position: offset})
}

func cullRectangles(boundaries []PlacedRectangle, rectangles []Rectangle) (validRectangles []Rectangle, oversizedRectangles []Rectangle) {

	fmt.Println("Culling")
	fmt.Println("Rectangles: ", len(rectangles))
	fmt.Println("Boundaries: ", len(boundaries))

	for _, r := range rectangles {
		willFit := false
		for _, b := range boundaries {
			if checkRectangleFit(r, b.Rectangle) {

				willFit = true
				break
			}
		}
		if willFit {
			// fmt.Println("Fits")
			validRectangles = append(validRectangles, r)
		} else {
			// fmt.Println("Oversized:")
			// rectangleInfo(r)
			// fmt.Println("")

			// for _, b := range boundaries {
			// 	fmt.Print("Boundary: w:")
			// 	fmt.Print(b.Rectangle.Width)
			// 	fmt.Print(" h: ")
			// 	fmt.Print(b.Rectangle.Height)
			// 	fmt.Println("")
			// }

			oversizedRectangles = append(oversizedRectangles, r)
		}
	}

	fmt.Println("validRectangles: ", len(validRectangles))
	fmt.Println("oversizedRectangles: ", len(oversizedRectangles))
	fmt.Println("")

	return validRectangles, oversizedRectangles

}

func checkRectangleFit(rectangle Rectangle, boundary Rectangle) bool {

	if rectangle.Height <= boundary.Height && rectangle.Width <= boundary.Width {
		return true
	}
	return false

}

func rectangleInfo(r Rectangle) {
	fmt.Println("Width: ", r.Width)
	fmt.Println("Height: ", r.Height)
}
