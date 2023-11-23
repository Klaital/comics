/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import Long = require("long");

export const protobufPackage = "";

export interface UpdateComicRequest {
  UserId: number;
  NewState: ComicState | undefined;
}

export interface GetComicsRequest {
  UserId: number;
}

export interface ComicState {
  ComicId: number;
  Name: string;
  HomeUrl: string;
  BookmarkUrl: string;
  LastReadTimestamp: string;
  RssFeedUrl: string;
}

export interface GetComicsResponse {
  Comics: ComicState[];
}

export interface MarkReadRequest {
  UserId: number;
  ComicId: number;
  /** If omitted, use current time */
  ReadAt?: string | undefined;
}

export interface MarkItemReadRequest {
  UserId: number;
  ComicId: number;
  /** If omitted, use current time */
  ReadAt?: string | undefined;
  ItemID: number;
}

export interface MarkReadResponse {
}

export interface DescribeComicsRequest {
  UserId: number;
  ComicIDs: number[];
}

export interface ComicDetails {
  ComicId: number;
  Name: string;
  HomeUrl: string;
  BookmarkUrl: string;
  LastReadTimestamp: string;
  RssFeedUrl: string;
  Feed: RssItem[];
}

export interface RssItem {
  ID: number;
  Link: string;
  Guid: string;
  IsRead: boolean;
  Title: string;
}

export interface DescribeComicsResponse {
  Comics: ComicDetails[];
}

export interface RefreshRssFeedResponse {
}

function createBaseUpdateComicRequest(): UpdateComicRequest {
  return { UserId: 0, NewState: undefined };
}

export const UpdateComicRequest = {
  encode(message: UpdateComicRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UserId !== 0) {
      writer.uint32(8).uint64(message.UserId);
    }
    if (message.NewState !== undefined) {
      ComicState.encode(message.NewState, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateComicRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateComicRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.UserId = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.NewState = ComicState.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateComicRequest {
    return {
      UserId: isSet(object.UserId) ? globalThis.Number(object.UserId) : 0,
      NewState: isSet(object.NewState) ? ComicState.fromJSON(object.NewState) : undefined,
    };
  },

  toJSON(message: UpdateComicRequest): unknown {
    const obj: any = {};
    if (message.UserId !== 0) {
      obj.UserId = Math.round(message.UserId);
    }
    if (message.NewState !== undefined) {
      obj.NewState = ComicState.toJSON(message.NewState);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateComicRequest>, I>>(base?: I): UpdateComicRequest {
    return UpdateComicRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateComicRequest>, I>>(object: I): UpdateComicRequest {
    const message = createBaseUpdateComicRequest();
    message.UserId = object.UserId ?? 0;
    message.NewState = (object.NewState !== undefined && object.NewState !== null)
      ? ComicState.fromPartial(object.NewState)
      : undefined;
    return message;
  },
};

function createBaseGetComicsRequest(): GetComicsRequest {
  return { UserId: 0 };
}

export const GetComicsRequest = {
  encode(message: GetComicsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UserId !== 0) {
      writer.uint32(8).uint64(message.UserId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetComicsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetComicsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.UserId = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetComicsRequest {
    return { UserId: isSet(object.UserId) ? globalThis.Number(object.UserId) : 0 };
  },

  toJSON(message: GetComicsRequest): unknown {
    const obj: any = {};
    if (message.UserId !== 0) {
      obj.UserId = Math.round(message.UserId);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetComicsRequest>, I>>(base?: I): GetComicsRequest {
    return GetComicsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetComicsRequest>, I>>(object: I): GetComicsRequest {
    const message = createBaseGetComicsRequest();
    message.UserId = object.UserId ?? 0;
    return message;
  },
};

function createBaseComicState(): ComicState {
  return { ComicId: 0, Name: "", HomeUrl: "", BookmarkUrl: "", LastReadTimestamp: "", RssFeedUrl: "" };
}

export const ComicState = {
  encode(message: ComicState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ComicId !== 0) {
      writer.uint32(8).uint64(message.ComicId);
    }
    if (message.Name !== "") {
      writer.uint32(18).string(message.Name);
    }
    if (message.HomeUrl !== "") {
      writer.uint32(26).string(message.HomeUrl);
    }
    if (message.BookmarkUrl !== "") {
      writer.uint32(34).string(message.BookmarkUrl);
    }
    if (message.LastReadTimestamp !== "") {
      writer.uint32(42).string(message.LastReadTimestamp);
    }
    if (message.RssFeedUrl !== "") {
      writer.uint32(50).string(message.RssFeedUrl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ComicState {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComicState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.ComicId = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.HomeUrl = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.BookmarkUrl = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.LastReadTimestamp = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.RssFeedUrl = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ComicState {
    return {
      ComicId: isSet(object.ComicId) ? globalThis.Number(object.ComicId) : 0,
      Name: isSet(object.Name) ? globalThis.String(object.Name) : "",
      HomeUrl: isSet(object.HomeUrl) ? globalThis.String(object.HomeUrl) : "",
      BookmarkUrl: isSet(object.BookmarkUrl) ? globalThis.String(object.BookmarkUrl) : "",
      LastReadTimestamp: isSet(object.LastReadTimestamp) ? globalThis.String(object.LastReadTimestamp) : "",
      RssFeedUrl: isSet(object.RssFeedUrl) ? globalThis.String(object.RssFeedUrl) : "",
    };
  },

  toJSON(message: ComicState): unknown {
    const obj: any = {};
    if (message.ComicId !== 0) {
      obj.ComicId = Math.round(message.ComicId);
    }
    if (message.Name !== "") {
      obj.Name = message.Name;
    }
    if (message.HomeUrl !== "") {
      obj.HomeUrl = message.HomeUrl;
    }
    if (message.BookmarkUrl !== "") {
      obj.BookmarkUrl = message.BookmarkUrl;
    }
    if (message.LastReadTimestamp !== "") {
      obj.LastReadTimestamp = message.LastReadTimestamp;
    }
    if (message.RssFeedUrl !== "") {
      obj.RssFeedUrl = message.RssFeedUrl;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ComicState>, I>>(base?: I): ComicState {
    return ComicState.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ComicState>, I>>(object: I): ComicState {
    const message = createBaseComicState();
    message.ComicId = object.ComicId ?? 0;
    message.Name = object.Name ?? "";
    message.HomeUrl = object.HomeUrl ?? "";
    message.BookmarkUrl = object.BookmarkUrl ?? "";
    message.LastReadTimestamp = object.LastReadTimestamp ?? "";
    message.RssFeedUrl = object.RssFeedUrl ?? "";
    return message;
  },
};

function createBaseGetComicsResponse(): GetComicsResponse {
  return { Comics: [] };
}

export const GetComicsResponse = {
  encode(message: GetComicsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Comics) {
      ComicState.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetComicsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetComicsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Comics.push(ComicState.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetComicsResponse {
    return {
      Comics: globalThis.Array.isArray(object?.Comics) ? object.Comics.map((e: any) => ComicState.fromJSON(e)) : [],
    };
  },

  toJSON(message: GetComicsResponse): unknown {
    const obj: any = {};
    if (message.Comics?.length) {
      obj.Comics = message.Comics.map((e) => ComicState.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetComicsResponse>, I>>(base?: I): GetComicsResponse {
    return GetComicsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetComicsResponse>, I>>(object: I): GetComicsResponse {
    const message = createBaseGetComicsResponse();
    message.Comics = object.Comics?.map((e) => ComicState.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMarkReadRequest(): MarkReadRequest {
  return { UserId: 0, ComicId: 0, ReadAt: undefined };
}

export const MarkReadRequest = {
  encode(message: MarkReadRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UserId !== 0) {
      writer.uint32(8).uint64(message.UserId);
    }
    if (message.ComicId !== 0) {
      writer.uint32(16).uint64(message.ComicId);
    }
    if (message.ReadAt !== undefined) {
      writer.uint32(26).string(message.ReadAt);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarkReadRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarkReadRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.UserId = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.ComicId = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.ReadAt = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarkReadRequest {
    return {
      UserId: isSet(object.UserId) ? globalThis.Number(object.UserId) : 0,
      ComicId: isSet(object.ComicId) ? globalThis.Number(object.ComicId) : 0,
      ReadAt: isSet(object.ReadAt) ? globalThis.String(object.ReadAt) : undefined,
    };
  },

  toJSON(message: MarkReadRequest): unknown {
    const obj: any = {};
    if (message.UserId !== 0) {
      obj.UserId = Math.round(message.UserId);
    }
    if (message.ComicId !== 0) {
      obj.ComicId = Math.round(message.ComicId);
    }
    if (message.ReadAt !== undefined) {
      obj.ReadAt = message.ReadAt;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MarkReadRequest>, I>>(base?: I): MarkReadRequest {
    return MarkReadRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MarkReadRequest>, I>>(object: I): MarkReadRequest {
    const message = createBaseMarkReadRequest();
    message.UserId = object.UserId ?? 0;
    message.ComicId = object.ComicId ?? 0;
    message.ReadAt = object.ReadAt ?? undefined;
    return message;
  },
};

function createBaseMarkItemReadRequest(): MarkItemReadRequest {
  return { UserId: 0, ComicId: 0, ReadAt: undefined, ItemID: 0 };
}

export const MarkItemReadRequest = {
  encode(message: MarkItemReadRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UserId !== 0) {
      writer.uint32(8).uint64(message.UserId);
    }
    if (message.ComicId !== 0) {
      writer.uint32(16).uint64(message.ComicId);
    }
    if (message.ReadAt !== undefined) {
      writer.uint32(26).string(message.ReadAt);
    }
    if (message.ItemID !== 0) {
      writer.uint32(32).uint64(message.ItemID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarkItemReadRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarkItemReadRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.UserId = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.ComicId = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.ReadAt = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.ItemID = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MarkItemReadRequest {
    return {
      UserId: isSet(object.UserId) ? globalThis.Number(object.UserId) : 0,
      ComicId: isSet(object.ComicId) ? globalThis.Number(object.ComicId) : 0,
      ReadAt: isSet(object.ReadAt) ? globalThis.String(object.ReadAt) : undefined,
      ItemID: isSet(object.ItemID) ? globalThis.Number(object.ItemID) : 0,
    };
  },

  toJSON(message: MarkItemReadRequest): unknown {
    const obj: any = {};
    if (message.UserId !== 0) {
      obj.UserId = Math.round(message.UserId);
    }
    if (message.ComicId !== 0) {
      obj.ComicId = Math.round(message.ComicId);
    }
    if (message.ReadAt !== undefined) {
      obj.ReadAt = message.ReadAt;
    }
    if (message.ItemID !== 0) {
      obj.ItemID = Math.round(message.ItemID);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MarkItemReadRequest>, I>>(base?: I): MarkItemReadRequest {
    return MarkItemReadRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MarkItemReadRequest>, I>>(object: I): MarkItemReadRequest {
    const message = createBaseMarkItemReadRequest();
    message.UserId = object.UserId ?? 0;
    message.ComicId = object.ComicId ?? 0;
    message.ReadAt = object.ReadAt ?? undefined;
    message.ItemID = object.ItemID ?? 0;
    return message;
  },
};

function createBaseMarkReadResponse(): MarkReadResponse {
  return {};
}

export const MarkReadResponse = {
  encode(_: MarkReadResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarkReadResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarkReadResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): MarkReadResponse {
    return {};
  },

  toJSON(_: MarkReadResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<MarkReadResponse>, I>>(base?: I): MarkReadResponse {
    return MarkReadResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MarkReadResponse>, I>>(_: I): MarkReadResponse {
    const message = createBaseMarkReadResponse();
    return message;
  },
};

function createBaseDescribeComicsRequest(): DescribeComicsRequest {
  return { UserId: 0, ComicIDs: [] };
}

export const DescribeComicsRequest = {
  encode(message: DescribeComicsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UserId !== 0) {
      writer.uint32(8).uint64(message.UserId);
    }
    writer.uint32(18).fork();
    for (const v of message.ComicIDs) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DescribeComicsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDescribeComicsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.UserId = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag === 16) {
            message.ComicIDs.push(longToNumber(reader.uint64() as Long));

            continue;
          }

          if (tag === 18) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.ComicIDs.push(longToNumber(reader.uint64() as Long));
            }

            continue;
          }

          break;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DescribeComicsRequest {
    return {
      UserId: isSet(object.UserId) ? globalThis.Number(object.UserId) : 0,
      ComicIDs: globalThis.Array.isArray(object?.ComicIDs) ? object.ComicIDs.map((e: any) => globalThis.Number(e)) : [],
    };
  },

  toJSON(message: DescribeComicsRequest): unknown {
    const obj: any = {};
    if (message.UserId !== 0) {
      obj.UserId = Math.round(message.UserId);
    }
    if (message.ComicIDs?.length) {
      obj.ComicIDs = message.ComicIDs.map((e) => Math.round(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DescribeComicsRequest>, I>>(base?: I): DescribeComicsRequest {
    return DescribeComicsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DescribeComicsRequest>, I>>(object: I): DescribeComicsRequest {
    const message = createBaseDescribeComicsRequest();
    message.UserId = object.UserId ?? 0;
    message.ComicIDs = object.ComicIDs?.map((e) => e) || [];
    return message;
  },
};

function createBaseComicDetails(): ComicDetails {
  return { ComicId: 0, Name: "", HomeUrl: "", BookmarkUrl: "", LastReadTimestamp: "", RssFeedUrl: "", Feed: [] };
}

export const ComicDetails = {
  encode(message: ComicDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ComicId !== 0) {
      writer.uint32(8).uint64(message.ComicId);
    }
    if (message.Name !== "") {
      writer.uint32(18).string(message.Name);
    }
    if (message.HomeUrl !== "") {
      writer.uint32(26).string(message.HomeUrl);
    }
    if (message.BookmarkUrl !== "") {
      writer.uint32(34).string(message.BookmarkUrl);
    }
    if (message.LastReadTimestamp !== "") {
      writer.uint32(42).string(message.LastReadTimestamp);
    }
    if (message.RssFeedUrl !== "") {
      writer.uint32(50).string(message.RssFeedUrl);
    }
    for (const v of message.Feed) {
      RssItem.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ComicDetails {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComicDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.ComicId = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.HomeUrl = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.BookmarkUrl = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.LastReadTimestamp = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.RssFeedUrl = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.Feed.push(RssItem.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ComicDetails {
    return {
      ComicId: isSet(object.ComicId) ? globalThis.Number(object.ComicId) : 0,
      Name: isSet(object.Name) ? globalThis.String(object.Name) : "",
      HomeUrl: isSet(object.HomeUrl) ? globalThis.String(object.HomeUrl) : "",
      BookmarkUrl: isSet(object.BookmarkUrl) ? globalThis.String(object.BookmarkUrl) : "",
      LastReadTimestamp: isSet(object.LastReadTimestamp) ? globalThis.String(object.LastReadTimestamp) : "",
      RssFeedUrl: isSet(object.RssFeedUrl) ? globalThis.String(object.RssFeedUrl) : "",
      Feed: globalThis.Array.isArray(object?.Feed) ? object.Feed.map((e: any) => RssItem.fromJSON(e)) : [],
    };
  },

  toJSON(message: ComicDetails): unknown {
    const obj: any = {};
    if (message.ComicId !== 0) {
      obj.ComicId = Math.round(message.ComicId);
    }
    if (message.Name !== "") {
      obj.Name = message.Name;
    }
    if (message.HomeUrl !== "") {
      obj.HomeUrl = message.HomeUrl;
    }
    if (message.BookmarkUrl !== "") {
      obj.BookmarkUrl = message.BookmarkUrl;
    }
    if (message.LastReadTimestamp !== "") {
      obj.LastReadTimestamp = message.LastReadTimestamp;
    }
    if (message.RssFeedUrl !== "") {
      obj.RssFeedUrl = message.RssFeedUrl;
    }
    if (message.Feed?.length) {
      obj.Feed = message.Feed.map((e) => RssItem.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ComicDetails>, I>>(base?: I): ComicDetails {
    return ComicDetails.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ComicDetails>, I>>(object: I): ComicDetails {
    const message = createBaseComicDetails();
    message.ComicId = object.ComicId ?? 0;
    message.Name = object.Name ?? "";
    message.HomeUrl = object.HomeUrl ?? "";
    message.BookmarkUrl = object.BookmarkUrl ?? "";
    message.LastReadTimestamp = object.LastReadTimestamp ?? "";
    message.RssFeedUrl = object.RssFeedUrl ?? "";
    message.Feed = object.Feed?.map((e) => RssItem.fromPartial(e)) || [];
    return message;
  },
};

function createBaseRssItem(): RssItem {
  return { ID: 0, Link: "", Guid: "", IsRead: false, Title: "" };
}

export const RssItem = {
  encode(message: RssItem, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ID !== 0) {
      writer.uint32(8).uint64(message.ID);
    }
    if (message.Link !== "") {
      writer.uint32(18).string(message.Link);
    }
    if (message.Guid !== "") {
      writer.uint32(26).string(message.Guid);
    }
    if (message.IsRead === true) {
      writer.uint32(32).bool(message.IsRead);
    }
    if (message.Title !== "") {
      writer.uint32(42).string(message.Title);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RssItem {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRssItem();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.ID = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Link = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.Guid = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.IsRead = reader.bool();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.Title = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RssItem {
    return {
      ID: isSet(object.ID) ? globalThis.Number(object.ID) : 0,
      Link: isSet(object.Link) ? globalThis.String(object.Link) : "",
      Guid: isSet(object.Guid) ? globalThis.String(object.Guid) : "",
      IsRead: isSet(object.IsRead) ? globalThis.Boolean(object.IsRead) : false,
      Title: isSet(object.Title) ? globalThis.String(object.Title) : "",
    };
  },

  toJSON(message: RssItem): unknown {
    const obj: any = {};
    if (message.ID !== 0) {
      obj.ID = Math.round(message.ID);
    }
    if (message.Link !== "") {
      obj.Link = message.Link;
    }
    if (message.Guid !== "") {
      obj.Guid = message.Guid;
    }
    if (message.IsRead === true) {
      obj.IsRead = message.IsRead;
    }
    if (message.Title !== "") {
      obj.Title = message.Title;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RssItem>, I>>(base?: I): RssItem {
    return RssItem.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RssItem>, I>>(object: I): RssItem {
    const message = createBaseRssItem();
    message.ID = object.ID ?? 0;
    message.Link = object.Link ?? "";
    message.Guid = object.Guid ?? "";
    message.IsRead = object.IsRead ?? false;
    message.Title = object.Title ?? "";
    return message;
  },
};

function createBaseDescribeComicsResponse(): DescribeComicsResponse {
  return { Comics: [] };
}

export const DescribeComicsResponse = {
  encode(message: DescribeComicsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Comics) {
      ComicDetails.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DescribeComicsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDescribeComicsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Comics.push(ComicDetails.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DescribeComicsResponse {
    return {
      Comics: globalThis.Array.isArray(object?.Comics) ? object.Comics.map((e: any) => ComicDetails.fromJSON(e)) : [],
    };
  },

  toJSON(message: DescribeComicsResponse): unknown {
    const obj: any = {};
    if (message.Comics?.length) {
      obj.Comics = message.Comics.map((e) => ComicDetails.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DescribeComicsResponse>, I>>(base?: I): DescribeComicsResponse {
    return DescribeComicsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DescribeComicsResponse>, I>>(object: I): DescribeComicsResponse {
    const message = createBaseDescribeComicsResponse();
    message.Comics = object.Comics?.map((e) => ComicDetails.fromPartial(e)) || [];
    return message;
  },
};

function createBaseRefreshRssFeedResponse(): RefreshRssFeedResponse {
  return {};
}

export const RefreshRssFeedResponse = {
  encode(_: RefreshRssFeedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshRssFeedResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshRssFeedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): RefreshRssFeedResponse {
    return {};
  },

  toJSON(_: RefreshRssFeedResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<RefreshRssFeedResponse>, I>>(base?: I): RefreshRssFeedResponse {
    return RefreshRssFeedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RefreshRssFeedResponse>, I>>(_: I): RefreshRssFeedResponse {
    const message = createBaseRefreshRssFeedResponse();
    return message;
  },
};

export interface Comics {
  GetComics(request: GetComicsRequest): Promise<GetComicsResponse>;
  GetComicFeeds(request: DescribeComicsRequest): Promise<DescribeComicsResponse>;
  UpdateComic(request: UpdateComicRequest): Promise<GetComicsResponse>;
  MarkAsRead(request: MarkReadRequest): Promise<MarkReadResponse>;
  MarkItemRead(request: MarkItemReadRequest): Promise<MarkReadResponse>;
  RefreshRssFeed(request: GetComicsRequest): Promise<RefreshRssFeedResponse>;
}

export const ComicsServiceName = "Comics";
export class ComicsClientImpl implements Comics {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || ComicsServiceName;
    this.rpc = rpc;
    this.GetComics = this.GetComics.bind(this);
    this.GetComicFeeds = this.GetComicFeeds.bind(this);
    this.UpdateComic = this.UpdateComic.bind(this);
    this.MarkAsRead = this.MarkAsRead.bind(this);
    this.MarkItemRead = this.MarkItemRead.bind(this);
    this.RefreshRssFeed = this.RefreshRssFeed.bind(this);
  }
  GetComics(request: GetComicsRequest): Promise<GetComicsResponse> {
    const data = GetComicsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetComics", data);
    return promise.then((data) => GetComicsResponse.decode(_m0.Reader.create(data)));
  }

  GetComicFeeds(request: DescribeComicsRequest): Promise<DescribeComicsResponse> {
    const data = DescribeComicsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetComicFeeds", data);
    return promise.then((data) => DescribeComicsResponse.decode(_m0.Reader.create(data)));
  }

  UpdateComic(request: UpdateComicRequest): Promise<GetComicsResponse> {
    const data = UpdateComicRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "UpdateComic", data);
    return promise.then((data) => GetComicsResponse.decode(_m0.Reader.create(data)));
  }

  MarkAsRead(request: MarkReadRequest): Promise<MarkReadResponse> {
    const data = MarkReadRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "MarkAsRead", data);
    return promise.then((data) => MarkReadResponse.decode(_m0.Reader.create(data)));
  }

  MarkItemRead(request: MarkItemReadRequest): Promise<MarkReadResponse> {
    const data = MarkItemReadRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "MarkItemRead", data);
    return promise.then((data) => MarkReadResponse.decode(_m0.Reader.create(data)));
  }

  RefreshRssFeed(request: GetComicsRequest): Promise<RefreshRssFeedResponse> {
    const data = GetComicsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RefreshRssFeed", data);
    return promise.then((data) => RefreshRssFeedResponse.decode(_m0.Reader.create(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(globalThis.Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
