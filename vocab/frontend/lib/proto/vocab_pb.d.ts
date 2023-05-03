import * as jspb from 'google-protobuf'



export class VocabAdditional extends jspb.Message {
  getInfo(): string;
  setInfo(value: string): VocabAdditional;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VocabAdditional.AsObject;
  static toObject(includeInstance: boolean, msg: VocabAdditional): VocabAdditional.AsObject;
  static serializeBinaryToWriter(message: VocabAdditional, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VocabAdditional;
  static deserializeBinaryFromReader(message: VocabAdditional, reader: jspb.BinaryReader): VocabAdditional;
}

export namespace VocabAdditional {
  export type AsObject = {
    info: string,
  }
}

export class Vocab extends jspb.Message {
  getWord(): string;
  setWord(value: string): Vocab;

  getDescription(): string;
  setDescription(value: string): Vocab;

  getTranslation(): string;
  setTranslation(value: string): Vocab;

  getInfoList(): Array<string>;
  setInfoList(value: Array<string>): Vocab;
  clearInfoList(): Vocab;
  addInfo(value: string, index?: number): Vocab;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Vocab.AsObject;
  static toObject(includeInstance: boolean, msg: Vocab): Vocab.AsObject;
  static serializeBinaryToWriter(message: Vocab, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Vocab;
  static deserializeBinaryFromReader(message: Vocab, reader: jspb.BinaryReader): Vocab;
}

export namespace Vocab {
  export type AsObject = {
    word: string,
    description: string,
    translation: string,
    infoList: Array<string>,
  }
}

export class VocabListResponse extends jspb.Message {
  getVocabList(): Array<Vocab>;
  setVocabList(value: Array<Vocab>): VocabListResponse;
  clearVocabList(): VocabListResponse;
  addVocab(value?: Vocab, index?: number): Vocab;

  getTotalCount(): number;
  setTotalCount(value: number): VocabListResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VocabListResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VocabListResponse): VocabListResponse.AsObject;
  static serializeBinaryToWriter(message: VocabListResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VocabListResponse;
  static deserializeBinaryFromReader(message: VocabListResponse, reader: jspb.BinaryReader): VocabListResponse;
}

export namespace VocabListResponse {
  export type AsObject = {
    vocabList: Array<Vocab.AsObject>,
    totalCount: number,
  }
}

export class VocabListRequest extends jspb.Message {
  getPageNumber(): number;
  setPageNumber(value: number): VocabListRequest;

  getPageSize(): number;
  setPageSize(value: number): VocabListRequest;

  getPagination(): Pagination | undefined;
  setPagination(value?: Pagination): VocabListRequest;
  hasPagination(): boolean;
  clearPagination(): VocabListRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VocabListRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VocabListRequest): VocabListRequest.AsObject;
  static serializeBinaryToWriter(message: VocabListRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VocabListRequest;
  static deserializeBinaryFromReader(message: VocabListRequest, reader: jspb.BinaryReader): VocabListRequest;
}

export namespace VocabListRequest {
  export type AsObject = {
    pageNumber: number,
    pageSize: number,
    pagination?: Pagination.AsObject,
  }

  export enum PaginationCase { 
    _PAGINATION_NOT_SET = 0,
    PAGINATION = 3,
  }
}

export class Pagination extends jspb.Message {
  getStart(): number;
  setStart(value: number): Pagination;

  getEnd(): number;
  setEnd(value: number): Pagination;

  getResultsPerPage(): number;
  setResultsPerPage(value: number): Pagination;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Pagination.AsObject;
  static toObject(includeInstance: boolean, msg: Pagination): Pagination.AsObject;
  static serializeBinaryToWriter(message: Pagination, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Pagination;
  static deserializeBinaryFromReader(message: Pagination, reader: jspb.BinaryReader): Pagination;
}

export namespace Pagination {
  export type AsObject = {
    start: number,
    end: number,
    resultsPerPage: number,
  }
}

