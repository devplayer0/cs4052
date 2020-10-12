#version 430

in vec3 vPosition;

uniform mat4 projection, camera, model;

void main() {
    gl_Position = projection * camera * model * vec4(vPosition, 1);
}
