package util

import (
	"github.com/devplayer0/cs4052/pkg/pb"
	"github.com/go-gl/mathgl/mgl32"
)

// PBVec2 converts an Object protobuf Vec2 to an mgl32.Vec2
func PBVec2(i *pb.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{i.X, i.Y}
}

// PBVec3 converts an Object protobuf Vec23to an mgl32.Vec3
func PBVec3(i *pb.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{i.X, i.Y, i.Z}
}

// PBVec4 converts an Object protobuf Vec4 to an mgl32.Vec4
func PBVec4(i *pb.Vec4) mgl32.Vec4 {
	return mgl32.Vec4{i.X, i.Y, i.Z, i.W}
}

// PBQuat converts an Object protobuf Vec4 to an mgl32.Quat
func PBQuat(i *pb.Vec4) mgl32.Quat {
	return mgl32.Quat{W: i.W, V: mgl32.Vec3{i.X, i.Y, i.Z}}
}

// PBMat4 converts an Object protobuf Mat4 to an mgl32.Mat4
func PBMat4(i *pb.Mat4) mgl32.Mat4 {
	return mgl32.Mat4FromRows(
		PBVec4(i.A),
		PBVec4(i.B),
		PBVec4(i.C),
		PBVec4(i.D),
	)
}
