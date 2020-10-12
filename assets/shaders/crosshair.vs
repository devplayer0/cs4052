#version 430

in vec2 vPosition;

uniform mat4 projection, model;

void main() {
    gl_Position = projection * model * vec4(vPosition, 0, 1);
}
