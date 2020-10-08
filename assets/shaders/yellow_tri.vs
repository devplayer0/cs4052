#version 430

in vec3 vPosition;

uniform mat4 transform;

void main() {
    gl_Position = transform * vec4(vPosition.x, vPosition.y, vPosition.z, 1.0);
}
