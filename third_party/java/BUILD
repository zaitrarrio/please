package(default_visibility = ['PUBLIC'])

maven_jar(
    name = 'junit',
    id = 'junit:junit:4.12',
    hash = 'sha1: a791201ac8a3d2a251045a52e264de01343ad2df',
)

maven_jar(
    name = 'hamcrest',
    id = 'org.hamcrest:hamcrest-all:1.3',
    hash = 'sha1: 6e72502584462b2d35a7a6a18fc5541c412f08cc',
)

maven_jar(
    name = 'guava',
    id = 'com.google.guava:guava:20.0',
    hash = 'sha1: 8908f57125132a4468ea9977705a83f61b94fa89',
    deps = [
        ':error_prone_annotations',
        ':jsr305',
    ],
)

maven_jar(
    name = 'error_prone_annotations',
    id = 'com.google.errorprone:error_prone_annotations:2.0.14',
    hash = 'c737ea66a3d80c93de9966834da1b3d7c9c0a73e',
)

maven_jar(
    name = 'jsr305',
    id = 'com.google.code.findbugs:jsr305:3.0.1',
    hash = '33d84dc1711681f51259bbcc1c2bde260a723039',
)

maven_jar(
    name = 'logback-classic',
    id = 'ch.qos.logback:logback-classic:1.1.3',
    hash = 'sha1: 245ace75b7a18931f103a79e74d53c66f1f90a9e',
    deps = [
        ':logback-core',
        ':slf4j-api',
    ],
)

maven_jar(
    name = 'logback-core',
    id = 'ch.qos.logback:logback-core:1.1.3',
    hash = 'sha1: 53732a6ab1f058ad76808646519801e30cc4a780',
)

maven_jar(
    name = 'slf4j-api',
    id = 'org.slf4j:slf4j-api:1.7.10',
    hash = 'sha1: 0ff5ea6b9f644c97ef1de83eb75e7cc304103549',
)

maven_jar(
    name = 'protobuf',
    id = 'com.google.protobuf:protobuf-java:3.0.0-beta-3',
    hash = 'ed118aa0276005ca49b18a3a2b77eaacdec95680',
    sources = False,
)

maven_jar(
    name = 'proguard',
    id = 'net.sf.proguard:proguard-base:5.3.1',
    hash = '259cab7182906260547ebb195d44bbf41943a8b9',
    sources = False,
)
