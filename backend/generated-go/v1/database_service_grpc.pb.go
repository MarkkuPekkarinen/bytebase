// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: v1/database_service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DatabaseService_GetDatabase_FullMethodName          = "/bytebase.v1.DatabaseService/GetDatabase"
	DatabaseService_BatchGetDatabases_FullMethodName    = "/bytebase.v1.DatabaseService/BatchGetDatabases"
	DatabaseService_ListDatabases_FullMethodName        = "/bytebase.v1.DatabaseService/ListDatabases"
	DatabaseService_UpdateDatabase_FullMethodName       = "/bytebase.v1.DatabaseService/UpdateDatabase"
	DatabaseService_BatchUpdateDatabases_FullMethodName = "/bytebase.v1.DatabaseService/BatchUpdateDatabases"
	DatabaseService_SyncDatabase_FullMethodName         = "/bytebase.v1.DatabaseService/SyncDatabase"
	DatabaseService_BatchSyncDatabases_FullMethodName   = "/bytebase.v1.DatabaseService/BatchSyncDatabases"
	DatabaseService_GetDatabaseMetadata_FullMethodName  = "/bytebase.v1.DatabaseService/GetDatabaseMetadata"
	DatabaseService_GetDatabaseSchema_FullMethodName    = "/bytebase.v1.DatabaseService/GetDatabaseSchema"
	DatabaseService_DiffSchema_FullMethodName           = "/bytebase.v1.DatabaseService/DiffSchema"
	DatabaseService_ListChangelogs_FullMethodName       = "/bytebase.v1.DatabaseService/ListChangelogs"
	DatabaseService_GetChangelog_FullMethodName         = "/bytebase.v1.DatabaseService/GetChangelog"
	DatabaseService_GetSchemaString_FullMethodName      = "/bytebase.v1.DatabaseService/GetSchemaString"
)

// DatabaseServiceClient is the client API for DatabaseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseServiceClient interface {
	// Permissions required: bb.databases.get
	GetDatabase(ctx context.Context, in *GetDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	// Permissions required: bb.databases.get
	BatchGetDatabases(ctx context.Context, in *BatchGetDatabasesRequest, opts ...grpc.CallOption) (*BatchGetDatabasesResponse, error)
	// Permissions required: bb.databases.list
	ListDatabases(ctx context.Context, in *ListDatabasesRequest, opts ...grpc.CallOption) (*ListDatabasesResponse, error)
	// Permissions required: bb.databases.update
	UpdateDatabase(ctx context.Context, in *UpdateDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	// Permissions required: bb.databases.update
	BatchUpdateDatabases(ctx context.Context, in *BatchUpdateDatabasesRequest, opts ...grpc.CallOption) (*BatchUpdateDatabasesResponse, error)
	// Permissions required: bb.databases.sync
	SyncDatabase(ctx context.Context, in *SyncDatabaseRequest, opts ...grpc.CallOption) (*SyncDatabaseResponse, error)
	// Permissions required: bb.databases.sync
	BatchSyncDatabases(ctx context.Context, in *BatchSyncDatabasesRequest, opts ...grpc.CallOption) (*BatchSyncDatabasesResponse, error)
	// Permissions required: bb.databases.getSchema
	GetDatabaseMetadata(ctx context.Context, in *GetDatabaseMetadataRequest, opts ...grpc.CallOption) (*DatabaseMetadata, error)
	// Permissions required: bb.databases.getSchema
	GetDatabaseSchema(ctx context.Context, in *GetDatabaseSchemaRequest, opts ...grpc.CallOption) (*DatabaseSchema, error)
	// Permissions required: bb.databases.get
	DiffSchema(ctx context.Context, in *DiffSchemaRequest, opts ...grpc.CallOption) (*DiffSchemaResponse, error)
	// Permissions required: bb.changelogs.list
	ListChangelogs(ctx context.Context, in *ListChangelogsRequest, opts ...grpc.CallOption) (*ListChangelogsResponse, error)
	// Permissions required: changelogs.get
	GetChangelog(ctx context.Context, in *GetChangelogRequest, opts ...grpc.CallOption) (*Changelog, error)
	// Permissions required: databases.getSchema
	GetSchemaString(ctx context.Context, in *GetSchemaStringRequest, opts ...grpc.CallOption) (*GetSchemaStringResponse, error)
}

type databaseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseServiceClient(cc grpc.ClientConnInterface) DatabaseServiceClient {
	return &databaseServiceClient{cc}
}

func (c *databaseServiceClient) GetDatabase(ctx context.Context, in *GetDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Database)
	err := c.cc.Invoke(ctx, DatabaseService_GetDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) BatchGetDatabases(ctx context.Context, in *BatchGetDatabasesRequest, opts ...grpc.CallOption) (*BatchGetDatabasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchGetDatabasesResponse)
	err := c.cc.Invoke(ctx, DatabaseService_BatchGetDatabases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) ListDatabases(ctx context.Context, in *ListDatabasesRequest, opts ...grpc.CallOption) (*ListDatabasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDatabasesResponse)
	err := c.cc.Invoke(ctx, DatabaseService_ListDatabases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) UpdateDatabase(ctx context.Context, in *UpdateDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Database)
	err := c.cc.Invoke(ctx, DatabaseService_UpdateDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) BatchUpdateDatabases(ctx context.Context, in *BatchUpdateDatabasesRequest, opts ...grpc.CallOption) (*BatchUpdateDatabasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchUpdateDatabasesResponse)
	err := c.cc.Invoke(ctx, DatabaseService_BatchUpdateDatabases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) SyncDatabase(ctx context.Context, in *SyncDatabaseRequest, opts ...grpc.CallOption) (*SyncDatabaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncDatabaseResponse)
	err := c.cc.Invoke(ctx, DatabaseService_SyncDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) BatchSyncDatabases(ctx context.Context, in *BatchSyncDatabasesRequest, opts ...grpc.CallOption) (*BatchSyncDatabasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchSyncDatabasesResponse)
	err := c.cc.Invoke(ctx, DatabaseService_BatchSyncDatabases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetDatabaseMetadata(ctx context.Context, in *GetDatabaseMetadataRequest, opts ...grpc.CallOption) (*DatabaseMetadata, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DatabaseMetadata)
	err := c.cc.Invoke(ctx, DatabaseService_GetDatabaseMetadata_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetDatabaseSchema(ctx context.Context, in *GetDatabaseSchemaRequest, opts ...grpc.CallOption) (*DatabaseSchema, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DatabaseSchema)
	err := c.cc.Invoke(ctx, DatabaseService_GetDatabaseSchema_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) DiffSchema(ctx context.Context, in *DiffSchemaRequest, opts ...grpc.CallOption) (*DiffSchemaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DiffSchemaResponse)
	err := c.cc.Invoke(ctx, DatabaseService_DiffSchema_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) ListChangelogs(ctx context.Context, in *ListChangelogsRequest, opts ...grpc.CallOption) (*ListChangelogsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListChangelogsResponse)
	err := c.cc.Invoke(ctx, DatabaseService_ListChangelogs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetChangelog(ctx context.Context, in *GetChangelogRequest, opts ...grpc.CallOption) (*Changelog, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Changelog)
	err := c.cc.Invoke(ctx, DatabaseService_GetChangelog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetSchemaString(ctx context.Context, in *GetSchemaStringRequest, opts ...grpc.CallOption) (*GetSchemaStringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetSchemaStringResponse)
	err := c.cc.Invoke(ctx, DatabaseService_GetSchemaString_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseServiceServer is the server API for DatabaseService service.
// All implementations must embed UnimplementedDatabaseServiceServer
// for forward compatibility.
type DatabaseServiceServer interface {
	// Permissions required: bb.databases.get
	GetDatabase(context.Context, *GetDatabaseRequest) (*Database, error)
	// Permissions required: bb.databases.get
	BatchGetDatabases(context.Context, *BatchGetDatabasesRequest) (*BatchGetDatabasesResponse, error)
	// Permissions required: bb.databases.list
	ListDatabases(context.Context, *ListDatabasesRequest) (*ListDatabasesResponse, error)
	// Permissions required: bb.databases.update
	UpdateDatabase(context.Context, *UpdateDatabaseRequest) (*Database, error)
	// Permissions required: bb.databases.update
	BatchUpdateDatabases(context.Context, *BatchUpdateDatabasesRequest) (*BatchUpdateDatabasesResponse, error)
	// Permissions required: bb.databases.sync
	SyncDatabase(context.Context, *SyncDatabaseRequest) (*SyncDatabaseResponse, error)
	// Permissions required: bb.databases.sync
	BatchSyncDatabases(context.Context, *BatchSyncDatabasesRequest) (*BatchSyncDatabasesResponse, error)
	// Permissions required: bb.databases.getSchema
	GetDatabaseMetadata(context.Context, *GetDatabaseMetadataRequest) (*DatabaseMetadata, error)
	// Permissions required: bb.databases.getSchema
	GetDatabaseSchema(context.Context, *GetDatabaseSchemaRequest) (*DatabaseSchema, error)
	// Permissions required: bb.databases.get
	DiffSchema(context.Context, *DiffSchemaRequest) (*DiffSchemaResponse, error)
	// Permissions required: bb.changelogs.list
	ListChangelogs(context.Context, *ListChangelogsRequest) (*ListChangelogsResponse, error)
	// Permissions required: changelogs.get
	GetChangelog(context.Context, *GetChangelogRequest) (*Changelog, error)
	// Permissions required: databases.getSchema
	GetSchemaString(context.Context, *GetSchemaStringRequest) (*GetSchemaStringResponse, error)
	mustEmbedUnimplementedDatabaseServiceServer()
}

// UnimplementedDatabaseServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDatabaseServiceServer struct{}

func (UnimplementedDatabaseServiceServer) GetDatabase(context.Context, *GetDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) BatchGetDatabases(context.Context, *BatchGetDatabasesRequest) (*BatchGetDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGetDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) ListDatabases(context.Context, *ListDatabasesRequest) (*ListDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) UpdateDatabase(context.Context, *UpdateDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) BatchUpdateDatabases(context.Context, *BatchUpdateDatabasesRequest) (*BatchUpdateDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpdateDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) SyncDatabase(context.Context, *SyncDatabaseRequest) (*SyncDatabaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) BatchSyncDatabases(context.Context, *BatchSyncDatabasesRequest) (*BatchSyncDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchSyncDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) GetDatabaseMetadata(context.Context, *GetDatabaseMetadataRequest) (*DatabaseMetadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabaseMetadata not implemented")
}
func (UnimplementedDatabaseServiceServer) GetDatabaseSchema(context.Context, *GetDatabaseSchemaRequest) (*DatabaseSchema, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabaseSchema not implemented")
}
func (UnimplementedDatabaseServiceServer) DiffSchema(context.Context, *DiffSchemaRequest) (*DiffSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiffSchema not implemented")
}
func (UnimplementedDatabaseServiceServer) ListChangelogs(context.Context, *ListChangelogsRequest) (*ListChangelogsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListChangelogs not implemented")
}
func (UnimplementedDatabaseServiceServer) GetChangelog(context.Context, *GetChangelogRequest) (*Changelog, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChangelog not implemented")
}
func (UnimplementedDatabaseServiceServer) GetSchemaString(context.Context, *GetSchemaStringRequest) (*GetSchemaStringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchemaString not implemented")
}
func (UnimplementedDatabaseServiceServer) mustEmbedUnimplementedDatabaseServiceServer() {}
func (UnimplementedDatabaseServiceServer) testEmbeddedByValue()                         {}

// UnsafeDatabaseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseServiceServer will
// result in compilation errors.
type UnsafeDatabaseServiceServer interface {
	mustEmbedUnimplementedDatabaseServiceServer()
}

func RegisterDatabaseServiceServer(s grpc.ServiceRegistrar, srv DatabaseServiceServer) {
	// If the following call pancis, it indicates UnimplementedDatabaseServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DatabaseService_ServiceDesc, srv)
}

func _DatabaseService_GetDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_GetDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabase(ctx, req.(*GetDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_BatchGetDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).BatchGetDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_BatchGetDatabases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).BatchGetDatabases(ctx, req.(*BatchGetDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_ListDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).ListDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_ListDatabases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).ListDatabases(ctx, req.(*ListDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_UpdateDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).UpdateDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_UpdateDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).UpdateDatabase(ctx, req.(*UpdateDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_BatchUpdateDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpdateDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).BatchUpdateDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_BatchUpdateDatabases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).BatchUpdateDatabases(ctx, req.(*BatchUpdateDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_SyncDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).SyncDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_SyncDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).SyncDatabase(ctx, req.(*SyncDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_BatchSyncDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchSyncDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).BatchSyncDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_BatchSyncDatabases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).BatchSyncDatabases(ctx, req.(*BatchSyncDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetDatabaseMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabaseMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_GetDatabaseMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabaseMetadata(ctx, req.(*GetDatabaseMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetDatabaseSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabaseSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_GetDatabaseSchema_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabaseSchema(ctx, req.(*GetDatabaseSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_DiffSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiffSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).DiffSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_DiffSchema_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).DiffSchema(ctx, req.(*DiffSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_ListChangelogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChangelogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).ListChangelogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_ListChangelogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).ListChangelogs(ctx, req.(*ListChangelogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetChangelog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChangelogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetChangelog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_GetChangelog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetChangelog(ctx, req.(*GetChangelogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetSchemaString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetSchemaString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_GetSchemaString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetSchemaString(ctx, req.(*GetSchemaStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DatabaseService_ServiceDesc is the grpc.ServiceDesc for DatabaseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DatabaseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bytebase.v1.DatabaseService",
	HandlerType: (*DatabaseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDatabase",
			Handler:    _DatabaseService_GetDatabase_Handler,
		},
		{
			MethodName: "BatchGetDatabases",
			Handler:    _DatabaseService_BatchGetDatabases_Handler,
		},
		{
			MethodName: "ListDatabases",
			Handler:    _DatabaseService_ListDatabases_Handler,
		},
		{
			MethodName: "UpdateDatabase",
			Handler:    _DatabaseService_UpdateDatabase_Handler,
		},
		{
			MethodName: "BatchUpdateDatabases",
			Handler:    _DatabaseService_BatchUpdateDatabases_Handler,
		},
		{
			MethodName: "SyncDatabase",
			Handler:    _DatabaseService_SyncDatabase_Handler,
		},
		{
			MethodName: "BatchSyncDatabases",
			Handler:    _DatabaseService_BatchSyncDatabases_Handler,
		},
		{
			MethodName: "GetDatabaseMetadata",
			Handler:    _DatabaseService_GetDatabaseMetadata_Handler,
		},
		{
			MethodName: "GetDatabaseSchema",
			Handler:    _DatabaseService_GetDatabaseSchema_Handler,
		},
		{
			MethodName: "DiffSchema",
			Handler:    _DatabaseService_DiffSchema_Handler,
		},
		{
			MethodName: "ListChangelogs",
			Handler:    _DatabaseService_ListChangelogs_Handler,
		},
		{
			MethodName: "GetChangelog",
			Handler:    _DatabaseService_GetChangelog_Handler,
		},
		{
			MethodName: "GetSchemaString",
			Handler:    _DatabaseService_GetSchemaString_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/database_service.proto",
}
