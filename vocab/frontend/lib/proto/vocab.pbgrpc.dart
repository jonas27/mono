///
//  Generated code. Do not modify.
//  source: vocab.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'vocab.pb.dart' as $0;
export 'vocab.pb.dart';

class VocabServiceClient extends $grpc.Client {
  static final _$listVocabs =
      $grpc.ClientMethod<$0.VocabListRequest, $0.VocabListResponse>(
          '/vocab.VocabService/ListVocabs',
          ($0.VocabListRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $0.VocabListResponse.fromBuffer(value));

  VocabServiceClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options, interceptors: interceptors);

  $grpc.ResponseStream<$0.VocabListResponse> listVocabs(
      $0.VocabListRequest request,
      {$grpc.CallOptions? options}) {
    return $createStreamingCall(
        _$listVocabs, $async.Stream.fromIterable([request]),
        options: options);
  }
}

abstract class VocabServiceBase extends $grpc.Service {
  $core.String get $name => 'vocab.VocabService';

  VocabServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.VocabListRequest, $0.VocabListResponse>(
        'ListVocabs',
        listVocabs_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.VocabListRequest.fromBuffer(value),
        ($0.VocabListResponse value) => value.writeToBuffer()));
  }

  $async.Stream<$0.VocabListResponse> listVocabs_Pre($grpc.ServiceCall call,
      $async.Future<$0.VocabListRequest> request) async* {
    yield* listVocabs(call, await request);
  }

  $async.Stream<$0.VocabListResponse> listVocabs(
      $grpc.ServiceCall call, $0.VocabListRequest request);
}
