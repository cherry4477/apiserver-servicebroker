load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["fuzz.pb.go"],
    importpath = "github.com/gogo/protobuf/test/fuzztests",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/gogo/protobuf/gogoproto:go_default_library",
        "//vendor/github.com/gogo/protobuf/proto:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["fuzz_test.go"],
    importpath = "github.com/gogo/protobuf/test/fuzztests",
    library = ":go_default_library",
    deps = ["//vendor/github.com/gogo/protobuf/proto:go_default_library"],
)

proto_library(
    name = "fuzztests_proto",
    srcs = ["fuzz.proto"],
    visibility = ["//visibility:public"],
    deps = ["//github.com/gogo/protobuf/gogoproto:gogoproto_proto"],
)

go_proto_library(
    name = "fuzztests_go_proto",
    importpath = "github.com/gogo/protobuf/test/fuzztests",
    proto = ":fuzztests_proto",
    visibility = ["//visibility:public"],
    deps = ["//github.com/gogo/protobuf/gogoproto:go_default_library"],
)
