#version 430

#define MAX_JOINTS 256

in vec3 frag_pos;
in vec3 normal;
in vec2 uv;

in ivec4 joint_ids_a;
in ivec4 joint_ids_b;
in vec4 weights_a;
in vec4 weights_b;

out vec3 world_pos;
out vec3 world_normal;

uniform mat4 projection, camera, model;
uniform mat4 joints[MAX_JOINTS];

void main() {
    mat4 skinning  = joints[joint_ids_a[0]] * weights_a[0];
         skinning += joints[joint_ids_a[1]] * weights_a[1];
         skinning += joints[joint_ids_a[2]] * weights_a[2];
         skinning += joints[joint_ids_a[3]] * weights_a[3];
         skinning += joints[joint_ids_b[0]] * weights_b[0];
         skinning += joints[joint_ids_b[1]] * weights_b[1];
         skinning += joints[joint_ids_b[2]] * weights_b[2];
         skinning += joints[joint_ids_b[3]] * weights_b[3];

    world_pos = vec3(skinning * model * vec4(frag_pos, 1.0));
    world_normal = mat3(transpose(inverse(skinning * model))) * normal;
    gl_Position = projection * camera * model * skinning * vec4(frag_pos, 1.0);
}
