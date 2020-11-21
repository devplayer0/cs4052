#version 430

in vec3 frag_pos;
in vec3 normal;
in vec2 uv;

out vec3 world_pos;
out vec3 world_normal;

uniform mat4 projection, camera, model;

void main() {
    world_pos = vec3(model * vec4(frag_pos, 1.0));
    world_normal = mat3(transpose(inverse(model))) * normal;
    gl_Position = projection * camera * model * vec4(frag_pos, 1.0);
}
