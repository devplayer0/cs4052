#version 430

in vec3 frag_pos;

out vec3 tex_coords;

uniform mat4 projection, camera;

void main() {
    tex_coords = frag_pos;
    vec4 pos = projection * camera * vec4(frag_pos, 1);
    // Always fail the depth test if there's something else to be rendered
    gl_Position = pos.xyww;
}
