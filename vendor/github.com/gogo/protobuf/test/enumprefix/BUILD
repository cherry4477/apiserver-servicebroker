load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["enumprefix.pb.go"],
    importpath = "github.com/gogo/protobuf/test/enumprefix",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/gogo/protobuf/gogoproto:go_default_library",
        "//vendor/github.com/gogo/protobuf/proto:go_default_library",
        "//vendor/github.com/gogo/protobuf/test:go_default_library",
    ],
)

proto_library(
    name = "enumprefix_proto",
    srcs = ["enumprefix.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//github.com/gogo/protobuf/gogoproto:gogoproto_proto",
        "//github.com/gogo/protobuf/test:test_proto",
    ],
)

go_proto_library(
    name = "enumprefix_go_proto",
    importpath = "github.com/gogo/protobuf/test/enumprefix",
    proto = ":enumprefix_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//github.com/gogo/protobuf/gogoproto:go_default_library",
        "//github.com/gogo/protobuf/test:go_default_library",
    ],
)
