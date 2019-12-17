package packing

// Rectangle ...
type Rectangle struct {
	Width  int
	Height int
}

// Position ...
type Position struct {
	X int
	Y int
}

// PlacedObject ...
type PlacedObject struct {
	Rectangle Rectangle
	Position  Position
}

// Solution ...
type Solution struct {
	PlacedObjects []PlacedObject
	Boundary      Rectangle
}

// ObjectSet ...
type ObjectSet struct {
	Objects []Rectangle
}

// PlacedObjectSet ...
type PlacedObjectSet struct {
	PlacedObjects []PlacedObject
}

// BoundarySet ...
type BoundarySet struct {
	Boundaries []PlacedObject
}

// ParameterSet ...
type ParameterSet struct {
	Objects      ObjectSet
	Boundaries   BoundarySet
	CurrentState State
}

// State ...
type State struct {
	PlacedObjects PlacedObjectSet
}

// Solve ...
func Solve(boundaries BoundarySet, rectangles ObjectSet, currentState State) []State {
	parameterSets := []ParameterSet{}
	for b, boundary := range boundaries.Boundaries {
		for r, rectangle := range rectangles.Objects {
			if rectangleFits(rectangle, boundary) {
				objects := append([]Rectangle(nil), rectangles.Objects...)
				objects[r] = objects[len(objects)-1]
				objects = objects[:len(objects)-1]

				bounds := append([]PlacedObject(nil), boundaries.Boundaries...)
				bounds[b] = bounds[len(bounds)-1]
				bounds = bounds[:len(bounds)-1]

				state := currentState
				state.PlacedObjects.PlacedObjects = append(state.PlacedObjects.PlacedObjects,
					PlacedObject{
						Rectangle: rectangle,
						Position:  boundary.Position,
					},
				)

				newBounds := SplitBoundary(boundary, rectangle)

				if len(newBounds) == 0 {
					parameterSets = append(parameterSets,
						ParameterSet{
							Objects: ObjectSet{
								Objects: objects,
							},
							Boundaries:   BoundarySet{Boundaries: bounds},
							CurrentState: state,
						},
					)
				} else {
					for _, boundarySet := range newBounds {
						parameterSets = append(parameterSets,
							ParameterSet{
								Objects: ObjectSet{
									Objects: objects,
								},
								Boundaries:   BoundarySet{Boundaries: append(boundarySet.Boundaries, bounds...)},
								CurrentState: state,
							},
						)
					}
				}
			}

		}
	}

	if len(parameterSets) > 0 {

		stateSet := []State{}

		for _, parameterSet := range parameterSets {
			stateSet = append(stateSet, Solve(parameterSet.Boundaries, parameterSet.Objects, parameterSet.CurrentState)...)
		}

		return stateSet

	}

	return []State{currentState}

}

func rectangleFits(candidate Rectangle, boundary PlacedObject) bool {
	if candidate.Height <= boundary.Rectangle.Height && candidate.Width <= boundary.Rectangle.Width {
		return true
	}
	return false
}

// SplitBoundary ...
func SplitBoundary(boundary PlacedObject, rectangle Rectangle) []BoundarySet {
	if boundary.Rectangle.Width == rectangle.Width && boundary.Rectangle.Height == rectangle.Height {
		return []BoundarySet{}
	} else if boundary.Rectangle.Width == rectangle.Width && boundary.Rectangle.Height != rectangle.Height {
		return []BoundarySet{
			BoundarySet{
				Boundaries: []PlacedObject{
					PlacedObject{
						Rectangle: Rectangle{
							Width:  boundary.Rectangle.Width,
							Height: boundary.Rectangle.Height - rectangle.Height,
						},
						Position: Position{
							X: boundary.Position.X,
							Y: boundary.Position.Y + rectangle.Height,
						},
					},
				},
			},
		}
	} else if boundary.Rectangle.Height == rectangle.Height && boundary.Rectangle.Width != rectangle.Width {
		return []BoundarySet{
			BoundarySet{
				Boundaries: []PlacedObject{
					PlacedObject{
						Rectangle: Rectangle{
							Width:  boundary.Rectangle.Width - rectangle.Width,
							Height: boundary.Rectangle.Height,
						},
						Position: Position{
							X: boundary.Position.X + rectangle.Width,
							Y: boundary.Position.Y,
						},
					},
				},
			},
		}
	}

	return []BoundarySet{
		BoundarySet{
			Boundaries: []PlacedObject{
				PlacedObject{
					Rectangle{
						Width:  boundary.Rectangle.Width - rectangle.Width,
						Height: boundary.Rectangle.Height,
					},
					Position{
						X: boundary.Position.X + rectangle.Width,
						Y: boundary.Position.Y,
					},
				},
				PlacedObject{
					Rectangle{
						Width:  rectangle.Width,
						Height: boundary.Rectangle.Height - rectangle.Height,
					},
					Position{
						X: boundary.Position.X,
						Y: boundary.Position.Y + rectangle.Height,
					},
				},
			},
		},
		BoundarySet{
			Boundaries: []PlacedObject{
				PlacedObject{
					Rectangle{
						Width:  boundary.Rectangle.Width,
						Height: boundary.Rectangle.Height - rectangle.Height,
					},
					Position{
						X: boundary.Position.X,
						Y: boundary.Position.Y + rectangle.Height,
					},
				},
				PlacedObject{
					Rectangle{
						Width:  boundary.Rectangle.Width - rectangle.Width,
						Height: rectangle.Height,
					},
					Position{
						X: boundary.Position.X + rectangle.Width,
						Y: boundary.Position.Y,
					},
				},
			},
		},
	}

}
