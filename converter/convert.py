#!/usr/bin/env python
import sys
import os

import numpy as np
from google.protobuf import text_format
import pyassimp
from pyassimp.postprocess import *
import object_pb2 as pbobj

class ConvertException(Exception):
    pass

def np_to_vec2(v2, np):
    v2.x = np[0]
    v2.y = np[1]
def np_to_vec3(v3, np):
    v3.x = np[0]
    v3.y = np[1]
    v3.z = np[2]
def np_to_vec4(v4, np):
    v4.x = np[0]
    v4.y = np[1]
    v4.z = np[2]
    v4.w = np[3]
def np_to_quat(v4, np):
    # assimp orders quats differently...
    v4.x = np[3]
    v4.y = np[0]
    v4.z = np[1]
    v4.w = np[2]

def np_to_mat4(m4, np):
    np_to_vec4(m4.a, np[0])
    np_to_vec4(m4.b, np[1])
    np_to_vec4(m4.c, np[2])
    np_to_vec4(m4.d, np[3])

class Converter:
    def __init__(self, s):
        self.scene = s

    def _flatten_all_nodes(self, n):
        self.simp_nodes[n.name] = n
        for child in n.children:
            self._flatten_all_nodes(child)

    def _add_node(self, n):
        n_id = len(self.obj.hierarchy)

        cn = self.obj.hierarchy.add()
        cn.name = n.name
        np_to_mat4(cn.transform, n.transformation)

        self.node_name_id[cn.name] = n_id
        return cn

    # Make sure all of the parent nodes of this converted node are in the hierarchy
    # and wired up to their children
    def _ensure_hierarchy(self, cn):
        n = self.simp_nodes[cn.name]
        if n.parent == self.scene:
            return

        # Our parent isn't in the hierarchy yet
        if n.parent.name not in self.node_name_id:
            self._add_node(n.parent)

        cn_id = self.node_name_id[cn.name]
        parent_cn = self.obj.hierarchy[self.node_name_id[n.parent.name]]
        if cn_id in parent_cn.children:
            # Parent already has this node as a child, we can tell that the
            # hierarchy already exists at this point
            return

        parent_cn.children.append(cn_id)
        self._ensure_hierarchy(parent_cn)

    # Find meshes in the hierarchy and compile their transform
    def _find_instances(self, n, transform):
        final = transform * n.transformation
        for m in n.meshes:
            i = self.obj.instances.add()
            i.meshID = self.mesh_name_id[m.name]
            np_to_mat4(i.transform, n.transformation)

        for child in n.children:
            self._find_instances(child, final)

    def convert(self):
        self.obj = pbobj.Object()

        self.mesh_name_id = {}
        for i, m in enumerate(self.scene.meshes):
            self.mesh_name_id[m.name] = i
        self._find_instances(self.scene.rootnode, np.identity(4))

        self.simp_nodes = {}
        self._flatten_all_nodes(self.scene.rootnode)

        self.node_name_id = {}
        # Ensure the root node gets ID 0
        self._add_node(self.scene.rootnode)
        # TODO: Why does the root node's transform seem to be messed up?
        #np_to_mat4(self.obj.hierarchy[0].transform, np.identity(4))

        joint_name_id = {}
        for m in self.scene.meshes:
            if m.texturecoords is not None:
                # Only support a single channel of U + V texture coordinates
                if m.numuvcomponents[0] != 2:
                    raise ConvertException('Only a 2 channel texture coordinate system is supported')

            cm = self.obj.meshes.add()
            cm.name = m.name
            for i, v in enumerate(m.vertices):
                cv = cm.vertices.add()
                np_to_vec3(cv.position, v)
                np_to_vec3(cv.normal, m.normals[i])

                if m.texturecoords is not None:
                    np_to_vec2(cv.uv, m.texturecoords[0][i])

            for f in m.faces:
                if len(f) != 3:
                    raise ConvertException('Non-triangular face! (this should be impossible...)')
                cf = cm.faces.add()
                cf.a = f[0]
                cf.b = f[1]
                cf.c = f[2]

            for b in m.bones:
                # There might be multiple meshes referring to the same bone!
                if b.name in joint_name_id:
                    j_id = joint_name_id[b.name]
                else:
                    j_id = len(self.obj.joints)
                    joint_name_id[b.name] = j_id

                j = self.obj.joints.add()
                np_to_mat4(j.inverseBind, b.offsetmatrix)

                # It's possible the node could already be in the hierarchy!
                if b.name in self.node_name_id:
                    cn = self.obj.hierarchy[self.node_name_id[b.name]]
                else:
                    cn = self._add_node(self.simp_nodes[b.name])
                cn.jointID = j_id

                self._ensure_hierarchy(cn)

                for w in b.weights:
                    cw = cm.weights[j_id].weights.add()
                    cw.vertex = w.vertexid
                    cw.weight = w.weight

        for a in self.scene.animations:
            ca = self.obj.animations.add()
            ca.name = a.name
            ca.duration = a.duration
            ca.tps = a.tickspersecond

            for c in a.channels:
                node_name = c.nodename.data.decode('utf8')
                if node_name not in self.node_name_id:
                    print(f'warning: animation {ca.name} operates on unneeded node {node_name}', file=sys.stderr)
                    continue

                cc = ca.channels.add()
                cc.nodeID = self.node_name_id[node_name]

                for p in c.positionkeys:
                    cp = cc.posFrames.add()
                    cp.time = p.time
                    np_to_vec3(cp.value, p.value)
                for r in c.rotationkeys:
                    cr = cc.rotFrames.add()
                    cr.time = r.time
                    np_to_quat(cr.value, r.value)
                for s in c.scalingkeys:
                    cs = cc.scaleFrames.add()
                    cs.time = s.time
                    np_to_vec3(cs.value, s.value)

        return self.obj

def main():
    if len(sys.argv) != 2:
        print(f'{sys.argv[0]} <file type>', file=sys.stderr)
        return 1

    with pyassimp.load(sys.stdin.buffer, file_type=sys.argv[1], processing=
            aiProcess_Triangulate           |
            aiProcess_JoinIdenticalVertices |
            aiProcess_GenSmoothNormals      |
            aiProcess_SortByPType           |
            aiProcess_CalcTangentSpace
        ) as scene:
        converter = Converter(scene)
        obj = converter.convert()

        if os.getenv('SOBJ_TEXT'):
            text_format.PrintMessage(obj, sys.stdout)
        else:
            sys.stdout.buffer.write(obj.SerializeToString())

    return 0

if __name__ == '__main__':
    sys.exit(main())