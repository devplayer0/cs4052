#version 430

in vec4 frag_pos;

#define N_LAMPS {{len .Lamps}}
uniform vec3 lamp_positions[N_LAMPS];
uniform float far_plane;

void main() {
    // gl_Layer represents the current face of the current cubemap array element
    // Divide by 6 to get the lamp index
    vec3 lamp_pos = lamp_positions[gl_Layer / 6];
    float lamp_distance = length(frag_pos.xyz - lamp_pos);

    // map to [0;1] range by dividing by far_plane
    lamp_distance /= far_plane;

    // write this as modified depth
    gl_FragDepth = lamp_distance;
}
