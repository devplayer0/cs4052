#version 430

in vec3 vPosition;
in vec3 Normal;

uniform mat4 projection, camera, model;

//flat out vec3 color;

void main() {
    gl_Position = projection * camera * model * vec4(vPosition, 1.0);
    //color = Normal;
}
