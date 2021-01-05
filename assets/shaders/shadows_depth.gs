#version 430

#define N_LAMPS {{len .Lamps}}
layout (triangles) in;
layout (triangle_strip, max_vertices={{mul 6 3 (len .Lamps)}}) out;

// Only emit geometry for lamps we actually need to update
uniform bool update_lamps[N_LAMPS];
// Set of transforms for each lamp (one for each face of the cubemap)
uniform mat4 shadow_transforms[N_LAMPS*6];

out vec4 frag_pos; // frag_pos from GS (output per emitvertex)

void main() {
    for (int lamp = 0; lamp < N_LAMPS; lamp++) {
        if (!update_lamps[lamp]) {
            continue;
        }

        for (int face = 0; face < 6; face++) {
            // We're using a cubemap array - each layer is a single face of an
            // element in the array
            gl_Layer = lamp*6 + face;
            // for each triangle's vertices
            for (int i = 0; i < 3; i++) {
                frag_pos = gl_in[i].gl_Position;
                gl_Position = shadow_transforms[lamp*6 + face] * frag_pos;

                EmitVertex();
            }

            EndPrimitive();
        }
    }
}
