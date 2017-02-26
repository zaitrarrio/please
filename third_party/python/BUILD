package(default_visibility = ['PUBLIC'])

pip_library(
    name = 'pex',
    version = '1.1.20',
    hashes = ['sha256: d9103f602dc9d915a6e2c2d2a3d8a67ee3c49e1b2e0855b0008b534b50b3b176'],
    patch = [
        'never_cache_pex.patch',
        'bytecode_timestamps.patch',
    ],
)

pip_library(
    name = 'xmlrunner',
    package_name = 'unittest-xml-reporting',
    version = '1.11.0',
    hashes = ['sha256: 4ba8380f99daf6fe951dea05f05977bc4129ca71da8303971b05e7e8ce917f97'],
    deps = [':six'],
)

pip_library(
    name = 'six',
    version = '1.10.0',
    outs = ['six.py'],
    hashes = ['sha256: b5ab7239fa0c1e2b4c2ce8ff34b17b4eb10b99078c104ca6b084bfb82b3d0e55'],
)

pip_library(
    name = 'pkg_resources',
    package_name = 'setuptools',
    version = '33.1.1',
    hashes = ['sha256: 1176f44404a33b10662a117131a72627a623189fb6b6b649b9023f0e5124b443'],
)

pip_library(
    name = 'requests',
    version = '2.5.0',
    hashes = [
        'sha256: bcad599053aed5aaa6079760171e6c04cfade75d80d1362f40d0b959d53e0337',
    ],
    test_only = True,  # Not used by plz itself.
)

pip_library(
    name = 'dateutil',
    package_name = 'python-dateutil',
    version = '2.4.0',
    hashes = [
        'sha256: 94234e68e7cf65f3ab564cc4016451f6027feddc3a02bd922dfb5489569efcfc',
    ],
    test_only = True,  # Not used by plz itself.
    deps = [':six'],
)

pip_library(
    name = 'protobuf',
    version = '3.0.0',
    outs = ['google'],
    deps = [':six'],
)
