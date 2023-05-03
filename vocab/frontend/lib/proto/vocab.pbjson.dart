///
//  Generated code. Do not modify.
//  source: vocab.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use vocabAdditionalDescriptor instead')
const VocabAdditional$json = const {
  '1': 'VocabAdditional',
  '2': const [
    const {'1': 'info', '3': 1, '4': 1, '5': 9, '10': 'info'},
  ],
};

/// Descriptor for `VocabAdditional`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vocabAdditionalDescriptor = $convert.base64Decode('Cg9Wb2NhYkFkZGl0aW9uYWwSEgoEaW5mbxgBIAEoCVIEaW5mbw==');
@$core.Deprecated('Use vocabDescriptor instead')
const Vocab$json = const {
  '1': 'Vocab',
  '2': const [
    const {'1': 'word', '3': 1, '4': 1, '5': 9, '10': 'word'},
    const {'1': 'description', '3': 2, '4': 1, '5': 9, '10': 'description'},
    const {'1': 'translation', '3': 3, '4': 1, '5': 9, '10': 'translation'},
    const {'1': 'info', '3': 4, '4': 3, '5': 9, '10': 'info'},
  ],
};

/// Descriptor for `Vocab`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vocabDescriptor = $convert.base64Decode('CgVWb2NhYhISCgR3b3JkGAEgASgJUgR3b3JkEiAKC2Rlc2NyaXB0aW9uGAIgASgJUgtkZXNjcmlwdGlvbhIgCgt0cmFuc2xhdGlvbhgDIAEoCVILdHJhbnNsYXRpb24SEgoEaW5mbxgEIAMoCVIEaW5mbw==');
@$core.Deprecated('Use vocabListResponseDescriptor instead')
const VocabListResponse$json = const {
  '1': 'VocabListResponse',
  '2': const [
    const {'1': 'vocab', '3': 1, '4': 3, '5': 11, '6': '.vocab.Vocab', '10': 'vocab'},
    const {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `VocabListResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vocabListResponseDescriptor = $convert.base64Decode('ChFWb2NhYkxpc3RSZXNwb25zZRIiCgV2b2NhYhgBIAMoCzIMLnZvY2FiLlZvY2FiUgV2b2NhYhIfCgt0b3RhbF9jb3VudBgCIAEoBVIKdG90YWxDb3VudA==');
@$core.Deprecated('Use vocabListRequestDescriptor instead')
const VocabListRequest$json = const {
  '1': 'VocabListRequest',
  '2': const [
    const {'1': 'page_number', '3': 1, '4': 1, '5': 5, '10': 'pageNumber'},
    const {'1': 'page_size', '3': 2, '4': 1, '5': 5, '10': 'pageSize'},
    const {'1': 'pagination', '3': 3, '4': 1, '5': 11, '6': '.vocab.Pagination', '9': 0, '10': 'pagination', '17': true},
  ],
  '8': const [
    const {'1': '_pagination'},
  ],
};

/// Descriptor for `VocabListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vocabListRequestDescriptor = $convert.base64Decode('ChBWb2NhYkxpc3RSZXF1ZXN0Eh8KC3BhZ2VfbnVtYmVyGAEgASgFUgpwYWdlTnVtYmVyEhsKCXBhZ2Vfc2l6ZRgCIAEoBVIIcGFnZVNpemUSNgoKcGFnaW5hdGlvbhgDIAEoCzIRLnZvY2FiLlBhZ2luYXRpb25IAFIKcGFnaW5hdGlvbogBAUINCgtfcGFnaW5hdGlvbg==');
@$core.Deprecated('Use paginationDescriptor instead')
const Pagination$json = const {
  '1': 'Pagination',
  '2': const [
    const {'1': 'start', '3': 1, '4': 1, '5': 5, '10': 'start'},
    const {'1': 'end', '3': 2, '4': 1, '5': 5, '10': 'end'},
    const {'1': 'results_per_page', '3': 3, '4': 1, '5': 5, '10': 'resultsPerPage'},
  ],
};

/// Descriptor for `Pagination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List paginationDescriptor = $convert.base64Decode('CgpQYWdpbmF0aW9uEhQKBXN0YXJ0GAEgASgFUgVzdGFydBIQCgNlbmQYAiABKAVSA2VuZBIoChByZXN1bHRzX3Blcl9wYWdlGAMgASgFUg5yZXN1bHRzUGVyUGFnZQ==');
