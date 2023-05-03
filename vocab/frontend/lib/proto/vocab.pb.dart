///
//  Generated code. Do not modify.
//  source: vocab.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class VocabAdditional extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'VocabAdditional', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'vocab'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'info')
    ..hasRequiredFields = false
  ;

  VocabAdditional._() : super();
  factory VocabAdditional({
    $core.String? info,
  }) {
    final _result = create();
    if (info != null) {
      _result.info = info;
    }
    return _result;
  }
  factory VocabAdditional.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VocabAdditional.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VocabAdditional clone() => VocabAdditional()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VocabAdditional copyWith(void Function(VocabAdditional) updates) => super.copyWith((message) => updates(message as VocabAdditional)) as VocabAdditional; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static VocabAdditional create() => VocabAdditional._();
  VocabAdditional createEmptyInstance() => create();
  static $pb.PbList<VocabAdditional> createRepeated() => $pb.PbList<VocabAdditional>();
  @$core.pragma('dart2js:noInline')
  static VocabAdditional getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VocabAdditional>(create);
  static VocabAdditional? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get info => $_getSZ(0);
  @$pb.TagNumber(1)
  set info($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasInfo() => $_has(0);
  @$pb.TagNumber(1)
  void clearInfo() => clearField(1);
}

class Vocab extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'Vocab', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'vocab'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'word')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'description')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'translation')
    ..pPS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'info')
    ..hasRequiredFields = false
  ;

  Vocab._() : super();
  factory Vocab({
    $core.String? word,
    $core.String? description,
    $core.String? translation,
    $core.Iterable<$core.String>? info,
  }) {
    final _result = create();
    if (word != null) {
      _result.word = word;
    }
    if (description != null) {
      _result.description = description;
    }
    if (translation != null) {
      _result.translation = translation;
    }
    if (info != null) {
      _result.info.addAll(info);
    }
    return _result;
  }
  factory Vocab.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Vocab.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Vocab clone() => Vocab()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Vocab copyWith(void Function(Vocab) updates) => super.copyWith((message) => updates(message as Vocab)) as Vocab; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static Vocab create() => Vocab._();
  Vocab createEmptyInstance() => create();
  static $pb.PbList<Vocab> createRepeated() => $pb.PbList<Vocab>();
  @$core.pragma('dart2js:noInline')
  static Vocab getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Vocab>(create);
  static Vocab? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get word => $_getSZ(0);
  @$pb.TagNumber(1)
  set word($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWord() => $_has(0);
  @$pb.TagNumber(1)
  void clearWord() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(2)
  set description($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(2)
  void clearDescription() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get translation => $_getSZ(2);
  @$pb.TagNumber(3)
  set translation($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasTranslation() => $_has(2);
  @$pb.TagNumber(3)
  void clearTranslation() => clearField(3);

  @$pb.TagNumber(4)
  $core.List<$core.String> get info => $_getList(3);
}

class VocabListResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'VocabListResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'vocab'), createEmptyInstance: create)
    ..pc<Vocab>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'vocab', $pb.PbFieldType.PM, subBuilder: Vocab.create)
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalCount', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  VocabListResponse._() : super();
  factory VocabListResponse({
    $core.Iterable<Vocab>? vocab,
    $core.int? totalCount,
  }) {
    final _result = create();
    if (vocab != null) {
      _result.vocab.addAll(vocab);
    }
    if (totalCount != null) {
      _result.totalCount = totalCount;
    }
    return _result;
  }
  factory VocabListResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VocabListResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VocabListResponse clone() => VocabListResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VocabListResponse copyWith(void Function(VocabListResponse) updates) => super.copyWith((message) => updates(message as VocabListResponse)) as VocabListResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static VocabListResponse create() => VocabListResponse._();
  VocabListResponse createEmptyInstance() => create();
  static $pb.PbList<VocabListResponse> createRepeated() => $pb.PbList<VocabListResponse>();
  @$core.pragma('dart2js:noInline')
  static VocabListResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VocabListResponse>(create);
  static VocabListResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<Vocab> get vocab => $_getList(0);

  @$pb.TagNumber(2)
  $core.int get totalCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalCount($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalCount() => clearField(2);
}

class VocabListRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'VocabListRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'vocab'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'pageNumber', $pb.PbFieldType.O3)
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'pageSize', $pb.PbFieldType.O3)
    ..aOM<Pagination>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'pagination', subBuilder: Pagination.create)
    ..hasRequiredFields = false
  ;

  VocabListRequest._() : super();
  factory VocabListRequest({
    $core.int? pageNumber,
    $core.int? pageSize,
    Pagination? pagination,
  }) {
    final _result = create();
    if (pageNumber != null) {
      _result.pageNumber = pageNumber;
    }
    if (pageSize != null) {
      _result.pageSize = pageSize;
    }
    if (pagination != null) {
      _result.pagination = pagination;
    }
    return _result;
  }
  factory VocabListRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VocabListRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VocabListRequest clone() => VocabListRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VocabListRequest copyWith(void Function(VocabListRequest) updates) => super.copyWith((message) => updates(message as VocabListRequest)) as VocabListRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static VocabListRequest create() => VocabListRequest._();
  VocabListRequest createEmptyInstance() => create();
  static $pb.PbList<VocabListRequest> createRepeated() => $pb.PbList<VocabListRequest>();
  @$core.pragma('dart2js:noInline')
  static VocabListRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VocabListRequest>(create);
  static VocabListRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get pageNumber => $_getIZ(0);
  @$pb.TagNumber(1)
  set pageNumber($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPageNumber() => $_has(0);
  @$pb.TagNumber(1)
  void clearPageNumber() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get pageSize => $_getIZ(1);
  @$pb.TagNumber(2)
  set pageSize($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPageSize() => $_has(1);
  @$pb.TagNumber(2)
  void clearPageSize() => clearField(2);

  @$pb.TagNumber(3)
  Pagination get pagination => $_getN(2);
  @$pb.TagNumber(3)
  set pagination(Pagination v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasPagination() => $_has(2);
  @$pb.TagNumber(3)
  void clearPagination() => clearField(3);
  @$pb.TagNumber(3)
  Pagination ensurePagination() => $_ensure(2);
}

class Pagination extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'Pagination', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'vocab'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'start', $pb.PbFieldType.O3)
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'end', $pb.PbFieldType.O3)
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'resultsPerPage', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  Pagination._() : super();
  factory Pagination({
    $core.int? start,
    $core.int? end,
    $core.int? resultsPerPage,
  }) {
    final _result = create();
    if (start != null) {
      _result.start = start;
    }
    if (end != null) {
      _result.end = end;
    }
    if (resultsPerPage != null) {
      _result.resultsPerPage = resultsPerPage;
    }
    return _result;
  }
  factory Pagination.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Pagination.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Pagination clone() => Pagination()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Pagination copyWith(void Function(Pagination) updates) => super.copyWith((message) => updates(message as Pagination)) as Pagination; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static Pagination create() => Pagination._();
  Pagination createEmptyInstance() => create();
  static $pb.PbList<Pagination> createRepeated() => $pb.PbList<Pagination>();
  @$core.pragma('dart2js:noInline')
  static Pagination getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Pagination>(create);
  static Pagination? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get start => $_getIZ(0);
  @$pb.TagNumber(1)
  set start($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasStart() => $_has(0);
  @$pb.TagNumber(1)
  void clearStart() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get end => $_getIZ(1);
  @$pb.TagNumber(2)
  set end($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasEnd() => $_has(1);
  @$pb.TagNumber(2)
  void clearEnd() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get resultsPerPage => $_getIZ(2);
  @$pb.TagNumber(3)
  set resultsPerPage($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasResultsPerPage() => $_has(2);
  @$pb.TagNumber(3)
  void clearResultsPerPage() => clearField(3);
}

