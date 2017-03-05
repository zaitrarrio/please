go_library(
    name = 'worker',
    srcs = ['worker.go'],
    deps = [
        '//src/test/proto:worker',
        '//third_party/go:context',
        '//third_party/go:grpc',
    ],
)

go_binary(
    name = 'please_test_worker',
    deps = [
        ':worker',
        '//src/cli',
    ],
)
