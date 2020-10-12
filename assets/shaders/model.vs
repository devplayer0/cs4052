#version 430

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vPosition;
in vec3 Normal;

flat out vec3 color;

void main() {
    gl_Position = projection * camera * model * vec4(vPosition, 1.0);
    color = Normal;
}
