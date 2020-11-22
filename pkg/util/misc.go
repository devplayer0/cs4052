package util

import "github.com/go-gl/mathgl/mgl32"

// LineBuffer is a list of vertices for a line
var LineBuffer = []mgl32.Vec3{
	{-1, -1, 0},
	{1, 1, 0},
}

// CubeVertices is a list of vertices for a cube
var CubeVertices = []mgl32.Vec3{
	{-1, -1, -1},
	{1, -1, -1},
	{1, 1, -1},
	{1, 1, -1},
	{-1, 1, -1},
	{-1, -1, -1},

	{-1, -1, 1},
	{1, -1, 1},
	{1, 1, 1},
	{1, 1, 1},
	{-1, 1, 1},
	{-1, -1, 1},

	{-1, 1, 1},
	{-1, 1, -1},
	{-1, -1, -1},
	{-1, -1, -1},
	{-1, -1, 1},
	{-1, 1, 1},

	{1, 1, 1},
	{1, 1, -1},
	{1, -1, -1},
	{1, -1, -1},
	{1, -1, 1},
	{1, 1, 1},

	{-1, -1, -1},
	{1, -1, -1},
	{1, -1, 1},
	{1, -1, 1},
	{-1, -1, 1},
	{-1, -1, -1},

	{-1, 1, -1},
	{1, 1, -1},
	{1, 1, 1},
	{1, 1, 1},
	{-1, 1, 1},
	{-1, 1, -1},
}
