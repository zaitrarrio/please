java_library(
    name = 'jacoco',
    exported_deps = [
        ':asm',
        ':jacoco.agent',
        ':jacoco.core',
        ':jacoco.report',
    ],
    visibility = ['PUBLIC'],
)

maven_jar(
    name = 'jacoco.core',
    id = 'org.jacoco:org.jacoco.core:0.7.7.201606060606',
    hash = 'sha256: 99a22e9acf1c939a71c479677f8cb1c7e37e24ab5eb674e07eb7598a195f1dd1',
)

maven_jar(
    name = 'jacoco.report',
    id = 'org.jacoco:org.jacoco.report:0.7.7.201606060606',
    hash = 'sha256: 456f62bcff7f374809a341eaedbd717c4a53ac6a30be45002ac4fe2d422e69b0',
)

maven_jar(
    name = 'jacoco.agent',
    id = 'org.jacoco:org.jacoco.agent:0.7.7.201606060606',
    hash = 'sha256: 3d80628381c241913a4b14a26a0da8e8fca0415b4e889dd6a44c4352b5bbb5fc',
)

maven_jar(
    name = 'asm',
    id = 'org.ow2.asm:asm-debug-all:5.0.4',
    hash = 'sha256: 4fbf7eee1236389e1257b37d3453199de425bec3410998cc20ed632d05193112',
)
