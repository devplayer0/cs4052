#version 430

in vec3 frag_pos;

uniform mat4 model;

void main() {
    gl_Position = model * vec4(frag_pos, 1);
}
