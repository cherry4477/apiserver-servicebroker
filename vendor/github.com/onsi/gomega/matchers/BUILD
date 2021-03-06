load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "and.go",
        "assignable_to_type_of_matcher.go",
        "be_a_directory.go",
        "be_a_regular_file.go",
        "be_an_existing_file.go",
        "be_closed_matcher.go",
        "be_empty_matcher.go",
        "be_equivalent_to_matcher.go",
        "be_false_matcher.go",
        "be_identical_to.go",
        "be_nil_matcher.go",
        "be_numerically_matcher.go",
        "be_sent_matcher.go",
        "be_temporally_matcher.go",
        "be_true_matcher.go",
        "be_zero_matcher.go",
        "consist_of.go",
        "contain_element_matcher.go",
        "contain_substring_matcher.go",
        "equal_matcher.go",
        "have_cap_matcher.go",
        "have_key_matcher.go",
        "have_key_with_value_matcher.go",
        "have_len_matcher.go",
        "have_occurred_matcher.go",
        "have_prefix_matcher.go",
        "have_suffix_matcher.go",
        "match_error_matcher.go",
        "match_json_matcher.go",
        "match_regexp_matcher.go",
        "match_xml_matcher.go",
        "match_yaml_matcher.go",
        "not.go",
        "or.go",
        "panic_matcher.go",
        "receive_matcher.go",
        "succeed_matcher.go",
        "type_support.go",
        "with_transform.go",
    ],
    importpath = "github.com/onsi/gomega/matchers",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/onsi/gomega/format:go_default_library",
        "//vendor/github.com/onsi/gomega/internal/oraclematcher:go_default_library",
        "//vendor/github.com/onsi/gomega/matchers/support/goraph/bipartitegraph:go_default_library",
        "//vendor/github.com/onsi/gomega/types:go_default_library",
        "//vendor/golang.org/x/net/html/charset:go_default_library",
        "//vendor/gopkg.in/yaml.v2:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = [
        "and_test.go",
        "assignable_to_type_of_matcher_test.go",
        "be_a_directory_test.go",
        "be_a_regular_file_test.go",
        "be_an_existing_file_test.go",
        "be_closed_matcher_test.go",
        "be_empty_matcher_test.go",
        "be_equivalent_to_matcher_test.go",
        "be_false_matcher_test.go",
        "be_identical_to_test.go",
        "be_nil_matcher_test.go",
        "be_numerically_matcher_test.go",
        "be_sent_matcher_test.go",
        "be_temporally_matcher_test.go",
        "be_true_matcher_test.go",
        "be_zero_matcher_test.go",
        "consist_of_test.go",
        "contain_element_matcher_test.go",
        "contain_substring_matcher_test.go",
        "equal_matcher_test.go",
        "have_cap_matcher_test.go",
        "have_key_matcher_test.go",
        "have_key_with_value_matcher_test.go",
        "have_len_matcher_test.go",
        "have_occurred_matcher_test.go",
        "have_prefix_matcher_test.go",
        "have_suffix_matcher_test.go",
        "match_error_matcher_test.go",
        "match_json_matcher_test.go",
        "match_regexp_matcher_test.go",
        "match_xml_matcher_test.go",
        "match_yaml_matcher_test.go",
        "matcher_tests_suite_test.go",
        "not_test.go",
        "or_test.go",
        "panic_matcher_test.go",
        "receive_matcher_test.go",
        "succeed_matcher_test.go",
        "with_transform_test.go",
    ],
    importpath = "github.com/onsi/gomega/matchers_test",
    deps = [
        ":go_default_library",
        "//vendor/github.com/onsi/ginkgo:go_default_library",
        "//vendor/github.com/onsi/gomega:go_default_library",
        "//vendor/github.com/onsi/gomega/types:go_default_library",
    ],
)
