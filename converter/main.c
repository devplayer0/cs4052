#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

#include <assimp/cimport.h>
#include <assimp/cfileio.h>
#include <assimp/scene.h>
#include <assimp/postprocess.h>

void convert(const struct aiScene *scene) {
    for (size_t i = 0; i < scene->mNumAnimations; i++) {
        const struct aiMesh *mesh = scene->mMeshes[i];
        const struct aiAnimation *anim = scene->mAnimations[i];
        for (size_t j = 0; j < anim->mNumChannels; j++) {
            fprintf(stderr, "channel: %s\n", anim->mChannels[j]->mNodeName.data);
        }
    }
}

int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "usage: %s <file>\n", argv[0]);
        return 1;
    }

    const struct aiScene *scene = aiImportFile(argv[1],
        aiProcess_Triangulate           |
        aiProcess_JoinIdenticalVertices |
        aiProcess_SortByPType           |
        aiProcess_CalcTangentSpace);
    if (!scene) {
        fprintf(stderr, "import failed: %s\n", aiGetErrorString());
        return -1;
    }

    convert(scene);

    aiReleaseImport(scene);
    return 0;
}
