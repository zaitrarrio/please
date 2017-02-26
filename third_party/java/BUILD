package(default_visibility = ['PUBLIC'])

maven_jar(
    name = 'junit',
    id = 'junit:junit:4.12',
    hash = 'sha256: 3a810d6ccf7806f54dd5499d1e2b475a49101054d0e2218c9071c0208e6e3f8a',
)

maven_jar(
    name = 'hamcrest',
    id = 'org.hamcrest:hamcrest-all:1.3',
    hash = 'sha256: 826e761897468838697dbb45210b885522c2905915a4badfbcc9beba23d04034',
)

maven_jar(
    name = 'logback-classic',
    id = 'ch.qos.logback:logback-classic:1.1.3',
    hash = 'sha256: 6286b50a19ec4aa9b5b540e91179ad87f5f7cff2413a10697b8c6693b25ae183',
    test_only = True,  # Only used in one test.
    deps = [
        ':logback-core',
        ':slf4j-api',
    ],
)

maven_jar(
    name = 'logback-core',
    id = 'ch.qos.logback:logback-core:1.1.3',
    hash = 'sha256: 7f3fc941dfe55bf18b426e8d7d9ef30f3d731c471f5b60c936cf8aca1934cc8a',
    test_only = True,
)

maven_jar(
    name = 'slf4j-api',
    id = 'org.slf4j:slf4j-api:1.7.10',
    hash = 'sha256: 6cdfca529defba9de3b464ed8b631ed5eb81ad7c8e1731be06f265ba69893c0a',
    test_only = True,
)

maven_jar(
    name = 'protobuf',
    id = 'com.google.protobuf:protobuf-java:3.0.0',
    hash = '88d92d14d257f10971fbb734ffcdd290817045c0',
    sources = False,
)
