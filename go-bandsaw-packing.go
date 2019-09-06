package packing

// Rectangle holds a rectangle height, width, and ability to rotate
type Rectangle struct {
	Width      float32
	Height     float32
	Rotateable bool
}

// Position holds X and Y coordinates
type Position struct {
	X float32
	Y float32
}

// PlacedObject holds a rectangle and it's position, according to the center of the rectangle
type PlacedObject struct {
	Rectangle Rectangle
	Position  Position
	Rotated   bool
}

// Solution holds a full set of placed rectangles
type Solution struct {
	PlacedObjects []PlacedObject
	Boundary      Rectangle
}

// GetSolutions finds all possible solutions given boundaries and rectangles
func GetSolutions(boundaries []Rectangle, rectangles []Rectangle, solutions *[]Solution, currentSolution *Solution) {
	if canRectanglesFit(&boundaries, &rectangles) {

	} else {
		*solutions = append(*solutions, *currentSolution)
	}
}

// canRectanglesFit checks if a given set of rectangles and boundaries allows at least on more rectangle to fit.
func canRectanglesFit(boundaries *[]Rectangle, rectangles *[]Rectangle) bool {
	for _, rectangle := range *rectangles {
		for _, boundary := range *boundaries {
			if rectangle.Height <= boundary.Height && rectangle.Width <= boundary.Width {
				return true
			} else if rectangle.Rotateable {
				if rectangle.Height <= boundary.Width && rectangle.Width <= boundary.Height {
					return true
				}
			}
		}
	}
	return false
}
