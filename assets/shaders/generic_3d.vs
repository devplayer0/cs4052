#version 430

in vec3 frag_pos;

uniform mat4 projection, camera, model;

void main() {
    gl_Position = projection * camera * model * vec4(frag_pos, 1);
}
