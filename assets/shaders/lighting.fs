#version 430

#define N_LAMPS {{len .Lamps}}

struct material {
    float spec_exponent;
};

struct attenuation_params {
    float constant;
    float linear;
    float quadratic;
};

struct lamp {
    attenuation_params attenuation;

    vec3 position;
    vec3 ambient, diffuse, specular;
};

in vec3 world_pos;
in vec3 world_normal;

out vec4 out_color;

// global
uniform vec3 view_pos;
uniform lamp lamps[N_LAMPS];

// per object
uniform vec3 in_color;
uniform material mat;

float get_attenuation(attenuation_params p, float dist) {
    return 1.0 / (p.constant + p.linear * dist + p.quadratic * (dist*dist));
}

vec3 lamp_phong(lamp l, vec3 world_pos, vec3 normal, vec3 view_dir) {
    vec3 lamp_dir = normalize(l.position - world_pos);

    // diffuse
    float diffuse = max(dot(normal, lamp_dir), 0.0);

    // specular
    vec3 reflect_dir = reflect(-lamp_dir, normal);
    float specular = pow(max(dot(view_dir, reflect_dir), 0.0), mat.spec_exponent);

    // attenuation
    float dist = length(l.position - world_pos);
    float attenuation = get_attenuation(l.attenuation, dist);

    vec3 result;
    result += l.ambient * in_color * attenuation;
    result += l.diffuse * diffuse * in_color * attenuation;
    result += l.specular * specular * in_color * attenuation;

    return result;
}

void main() {
    vec3 normal = normalize(world_normal);
    vec3 view_dir = normalize(view_pos - world_pos);

    vec3 result;
    for (int i = 0; i < N_LAMPS; i++) {
        result += lamp_phong(lamps[i], world_pos, normal, view_dir);
    }

    out_color = vec4(result, 1.0);
}
