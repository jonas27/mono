import * as grpcWeb from 'grpc-web';

import * as vocab_pb from './vocab_pb';


export class VocabServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listVocabs(
    request: vocab_pb.VocabListRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<vocab_pb.VocabListResponse>;

}

export class VocabServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  listVocabs(
    request: vocab_pb.VocabListRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<vocab_pb.VocabListResponse>;

}

