package(default_visibility = ['PUBLIC'])

maven_jar(
    name = 'junit',
    id = 'junit:junit:4.12',
    hash = 'sha256: 39285bdaa21ebf4594e1da91627bc08717a8dfadebb7d1a18f7aa44fce0187f9',
)

maven_jar(
    name = 'hamcrest',
    id = 'org.hamcrest:hamcrest-all:1.3',
    hash = 'sha256: 9bbac02785535fa353da408c873a53188e66f78e04f6b5010598d0591aff1b9d',
)

maven_jar(
    name = 'logback-classic',
    id = 'ch.qos.logback:logback-classic:1.1.3',
    hash = 'sha256: 413c6076c59527efe893d054355d410e411f6e4838921447c766ebadbb9c797e',
    test_only = True,  # Only used in one test.
    deps = [
        ':logback-core',
        ':slf4j-api',
    ],
)

maven_jar(
    name = 'logback-core',
    id = 'ch.qos.logback:logback-core:1.1.3',
    hash = 'sha256: 48b2580ae9e0d1b4fd8d0e67a10261122bf8d66fe90d1068201784f2cd7d6440',
    test_only = True,
)

maven_jar(
    name = 'slf4j-api',
    id = 'org.slf4j:slf4j-api:1.7.10',
    hash = 'sha256: c4f94794f57f4efee9cd9154b60b55dde9c8036f41d671d05008c7ba5b967582',
    test_only = True,
)

maven_jar(
    name = 'protobuf',
    id = 'com.google.protobuf:protobuf-java:3.0.0',
    hash = 'sha256: c1d2cea2398896ba7b64d63360b02358f979695d863ea37cb5c54f593b9b2fc1',
    sources = False,
)
