package app

import (
	"github.com/devplayer0/cs4052/pkg/object"
	"github.com/devplayer0/cs4052/pkg/util"
	"github.com/go-gl/mathgl/mgl32"
)

func makeScorpion(p *util.Program, dp *util.Program) (*object.Object, error) {
	obj := object.NewObject(&object.Joint{
		Keyframes: []mgl32.Mat4{
			mgl32.Ident4(),
		},
		Children: map[string]*object.Joint{
			"hat": {
				Keyframes: []mgl32.Mat4{
					mgl32.Translate3D(0, 0.3, -1.3).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(5))),
					mgl32.Translate3D(0, 0.3, -1.3),
					mgl32.Translate3D(0, 0.3, -1.3).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-5))),
					mgl32.Translate3D(0, 0.3, -1.3),
					mgl32.Translate3D(0, 0.3, -1.3).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(5))),
				},
			},

			"stinger": {
				Keyframes: []mgl32.Mat4{
					mgl32.Translate3D(0, 0, 1.3).Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(50))),
					mgl32.Translate3D(0, 0, 1.3).Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(55))).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(5))),
					mgl32.Translate3D(0, 0, 1.3).Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(60))),
					mgl32.Translate3D(0, 0, 1.3).Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(55))).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-5))),
					mgl32.Translate3D(0, 0, 1.3).Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(50))),
				},
				Children: map[string]*object.Joint{
					"end": {
						Keyframes: []mgl32.Mat4{
							mgl32.Translate3D(0, 1, 0),
						},
						Children: map[string]*object.Joint{
							"seg": {
								Keyframes: []mgl32.Mat4{
									mgl32.HomogRotate3DX(mgl32.DegToRad(-30)),
									mgl32.HomogRotate3DX(mgl32.DegToRad(-25)),
									mgl32.HomogRotate3DX(mgl32.DegToRad(-30)),
								},
								Children: map[string]*object.Joint{
									"end": {
										Keyframes: []mgl32.Mat4{
											mgl32.Translate3D(0, 1, 0),
										},
										Children: map[string]*object.Joint{
											"seg": {
												Keyframes: []mgl32.Mat4{
													mgl32.HomogRotate3DX(mgl32.DegToRad(-60)),
												},
												Children: map[string]*object.Joint{
													"end": {
														Keyframes: []mgl32.Mat4{
															mgl32.Translate3D(0, 1, 0),
														},
														Children: map[string]*object.Joint{
															"seg": {
																Keyframes: []mgl32.Mat4{
																	mgl32.HomogRotate3DX(mgl32.DegToRad(-70)),
																},
																Children: map[string]*object.Joint{
																	"end": {
																		Keyframes: []mgl32.Mat4{
																			mgl32.Translate3D(0, 1, 0),
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"back_left": {
				Keyframes: []mgl32.Mat4{
					mgl32.Translate3D(-0.8, 0, 1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(30))),
					mgl32.Translate3D(-0.8, 0, 1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(30))).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(30))),
					mgl32.Translate3D(-0.8, 0, 1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(30))),
				},
				Children: map[string]*object.Joint{
					"end": {
						Keyframes: []mgl32.Mat4{
							mgl32.Translate3D(-1, 0, 0),
						},
						Children: map[string]*object.Joint{
							"seg": {
								Keyframes: []mgl32.Mat4{
									mgl32.Translate3D(-0.2, -0.02, 0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(20))),
								},
								Children: map[string]*object.Joint{
									"end": {
										Keyframes: []mgl32.Mat4{
											mgl32.Translate3D(-1, 0, 0),
										},
									},
								},
							},
						},
					},
				},
			},
			"front_left": {
				Keyframes: []mgl32.Mat4{
					mgl32.Translate3D(-0.8, 0, -1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(30))),
					mgl32.Translate3D(-0.8, 0, -1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(30))).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-30))),
					mgl32.Translate3D(-0.8, 0, -1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(30))),
				},
				Children: map[string]*object.Joint{
					"end": {
						Keyframes: []mgl32.Mat4{
							mgl32.Translate3D(-1, 0, 0),
						},
						Children: map[string]*object.Joint{
							"seg": {
								Keyframes: []mgl32.Mat4{
									mgl32.Translate3D(-0.2, -0.02, 0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(20))),
								},
								Children: map[string]*object.Joint{
									"end": {
										Keyframes: []mgl32.Mat4{
											mgl32.Translate3D(-1, 0, 0),
										},
									},
								},
							},
						},
					},
				},
			},

			"back_right": {
				Keyframes: []mgl32.Mat4{
					mgl32.Translate3D(0.8, 0, 1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-30))),
					mgl32.Translate3D(0.8, 0, 1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-30))).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(-30))),
					mgl32.Translate3D(0.8, 0, 1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-30))),
				},
				Children: map[string]*object.Joint{
					"end": {
						Keyframes: []mgl32.Mat4{
							mgl32.Translate3D(1, 0, 0),
						},
						Children: map[string]*object.Joint{
							"seg": {
								Keyframes: []mgl32.Mat4{
									mgl32.Translate3D(0.2, -0.02, 0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-20))),
								},
								Children: map[string]*object.Joint{
									"end": {
										Keyframes: []mgl32.Mat4{
											mgl32.Translate3D(1, 0, 0),
										},
									},
								},
							},
						},
					},
				},
			},
			"front_right": {
				Keyframes: []mgl32.Mat4{
					mgl32.Translate3D(0.8, 0, -1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-30))),
					mgl32.Translate3D(0.8, 0, -1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-30))).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(30))),
					mgl32.Translate3D(0.8, 0, -1.0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-30))),
				},
				Children: map[string]*object.Joint{
					"end": {
						Keyframes: []mgl32.Mat4{
							mgl32.Translate3D(1, 0, 0),
						},
						Children: map[string]*object.Joint{
							"seg": {
								Keyframes: []mgl32.Mat4{
									mgl32.Translate3D(0.2, -0.02, 0).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(-20))),
								},
								Children: map[string]*object.Joint{
									"end": {
										Keyframes: []mgl32.Mat4{
											mgl32.Translate3D(1, 0, 0),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, mgl32.Translate3D(-3, 2, 1), dp)

	body, err := util.NewOBJMeshFile("assets/meshes/body.obj", mgl32.Ident4())
	body.Upload(p)

	legObj, err := util.ReadOBJFile("assets/meshes/leg.obj")
	if err != nil {
		return nil, err
	}

	stingerSegObj, err := util.ReadOBJFile("assets/meshes/stinger_seg.obj")
	if err != nil {
		return nil, err
	}
	stingerObj, err := util.ReadOBJFile("assets/meshes/stinger.obj")
	if err != nil {
		return nil, err
	}

	hat, err := util.NewOBJMeshFile("assets/meshes/hat.obj", mgl32.Ident4())
	if err != nil {
		return nil, err
	}
	hat.Upload(p)

	obj.Meshes = map[string]*object.Mesh{
		"body": {
			Mesh: body,
			VertexWeights: map[string][]float32{
				"root": {1},
			},
		},

		"hat": {
			Mesh: hat,
			VertexWeights: map[string][]float32{
				"root.hat": {1},
			},
		},

		"stinger_seg1": {
			Mesh: util.NewOBJMesh(stingerSegObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.stinger":     {1},
				"root.stinger.end": {1},
			},
		},
		"stinger_seg2": {
			Mesh: util.NewOBJMesh(stingerSegObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.stinger.end.seg":     {1},
				"root.stinger.end.seg.end": {1},
			},
		},
		"stinger_seg3": {
			Mesh: util.NewOBJMesh(stingerSegObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.stinger.end.seg.end.seg":     {1},
				"root.stinger.end.seg.end.seg.end": {1},
			},
		},
		"stinger_seg4": {
			Mesh: util.NewOBJMesh(stingerSegObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.stinger.end.seg.end.seg.end.seg":     {1},
				"root.stinger.end.seg.end.seg.end.seg.end": {1},
			},
		},
		"stinger": {
			Mesh: util.NewOBJMesh(stingerObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.stinger.end.seg.end.seg.end.seg.end": {1},
			},
		},

		"back_left_seg1": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.back_left":     {1},
				"root.back_left.end": {1},
			},
		},
		"back_left_seg2": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.back_left.end.seg":     {1},
				"root.back_left.end.seg.end": {1},
			},
		},

		"front_left_seg1": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.front_left":     {1},
				"root.front_left.end": {1},
			},
		},
		"front_left_seg2": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.front_left.end.seg":     {1},
				"root.front_left.end.seg.end": {1},
			},
		},

		"back_right_seg1": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.back_right":     {1},
				"root.back_right.end": {1},
			},
		},
		"back_right_seg2": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.back_right.end.seg":     {1},
				"root.back_right.end.seg.end": {1},
			},
		},

		"front_right_seg1": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.front_right":     {1},
				"root.front_right.end": {1},
			},
		},
		"front_right_seg2": {
			Mesh: util.NewOBJMesh(legObj, mgl32.Ident4()).Upload(p),
			VertexWeights: map[string][]float32{
				"root.front_right.end.seg":     {1},
				"root.front_right.end.seg.end": {1},
			},
		},
	}

	return obj, nil
}
