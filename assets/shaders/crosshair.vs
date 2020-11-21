#version 430

in vec2 frag_pos;

uniform mat4 projection, model;

void main() {
    gl_Position = projection * model * vec4(frag_pos, 0, 1);
}
