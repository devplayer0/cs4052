#version 430

#define MAX_JOINTS 256

in vec3 frag_pos;
{{if not .DepthPass}}
in vec3 normal;
in vec2 uv_in;
in vec3 tangent;
in vec3 bitangent;
{{end}}

in ivec4 joint_ids_a;
in ivec4 joint_ids_b;
in vec4 weights_a;
in vec4 weights_b;

{{if not .DepthPass}}
out vec3 world_pos;
out vec3 world_normal;
out mat3 TBN;
out vec2 uv;

uniform mat4 projection, camera, model;
{{else}}
uniform mat4 model;
{{end}}
uniform mat4 joints[MAX_JOINTS];

{{if not .DepthPass}}
uniform bool normal_map;
{{end}}

void main() {
    mat4 skinning  = joints[joint_ids_a[0]] * weights_a[0];
         skinning += joints[joint_ids_a[1]] * weights_a[1];
         skinning += joints[joint_ids_a[2]] * weights_a[2];
         skinning += joints[joint_ids_a[3]] * weights_a[3];
         skinning += joints[joint_ids_b[0]] * weights_b[0];
         skinning += joints[joint_ids_b[1]] * weights_b[1];
         skinning += joints[joint_ids_b[2]] * weights_b[2];
         skinning += joints[joint_ids_b[3]] * weights_b[3];

{{if not .DepthPass}}
    world_pos = vec3(model * skinning * vec4(frag_pos, 1.0));

    mat3 normal_matrix = transpose(inverse(mat3(model * skinning)));
    if (normal_map) {
        vec3 T = normalize(normal_matrix * tangent);
        vec3 N = normalize(normal_matrix * normal);
        T = normalize(T - dot(T, N) * N);
        vec3 B = cross(N, T);
        TBN = transpose(mat3(T, B, N));
    } else {
        world_normal = normalize(normal_matrix * normal);
    }

    uv = uv_in;

    gl_Position = projection * camera * model * skinning * vec4(frag_pos, 1.0);
{{else}}
    gl_Position = model * skinning * vec4(frag_pos, 1.0);
{{end}}
}
