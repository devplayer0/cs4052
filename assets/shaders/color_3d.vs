#version 430

in vec3 vPosition;
in vec4 vColor;

uniform mat4 projection, camera, model;

out vec4 color;

void main() {
    gl_Position = projection * camera * model * vec4(vPosition, 1);
    color = vColor;
}
