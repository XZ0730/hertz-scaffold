namespace go manager

struct BaseRequest{
}

struct BaseResponse{
    1: required i64 code,
    2: required string message,
}

struct GetParkSpaceReq{
    1: required string garage_id
}
	# Id            int64
	# GarageId      string
	# Uid           int64
	# VehicleNumber string
	# ParkState     int
	# ParkType      int
struct ParkSpaceModel{
    1: i64 id
    2: string garage_id
    3: i64 uid
    4: string vehicle_number
    5: i64 park_state // 当前是否有车停入
    6: i64 park_type // 是否空闲车位 park_type : 0 游客车位 1 固定车位无归属 2 固定且已有归属
    7: string duration
}

struct ParkRecordModel {
    1: string time
    2: string vehicle_number
    3: string gid
    4: i64 park_id
    5: i64 type // 0 入库 1 出库
}

struct GetParkRecordListResp{
    1: i64 code
    2: string msg
    3: string garage_id 
    4: list<ParkRecordModel> prlist
    5: i64 total
}


struct GetParkSpaceResponse{
    1: i64 code
    2: string msg
    3: string garage_id
    4: list<ParkSpaceModel> plist
    5: i64 total
}

struct GetParkRecordReq{
     1: required string garage_id
     2: required string start_time
     3: required string end_time
}

struct CntVolumeReq{
    1: required string garage_id
}

struct CntVolumeResp{
    1: i64 code
    2: string msg
    3: list<i64> cnt
}

struct VehicleAuditModel {
    1: string vehicle_number
    2: string brand
    3: string model
    4: string description
    5: string garage_id
    6: i64 park_id
    7: string image
    8: i64 uid
    9: i64 id
}

struct GetVehicleAuditReq{
    1: string garage_id
}

struct GetVehicleAuditResp{
    1: i64 code
    2: string msg
    3: list<VehicleAuditModel> va_list
    4: i64 total
}

struct AuditVehicleReq{
    1: i64 aid
    2: i64 is_audit
}

service ManagerService{
    GetParkSpaceResponse GetGarageParkSpace(1:GetParkSpaceReq req)(api.get="/auth/manager/parkspace")
    GetParkRecordListResp GetParkRecordList(1:GetParkRecordReq req)(api.get="/auth/manager/park_record")

    CntVolumeResp CntVolume(1:CntVolumeReq req)(api.get="/auth/manager/cnt_volume")
    
    GetVehicleAuditResp GetVehicleAudit(1:GetVehicleAuditReq req)(api.get="/auth/manager/vehicle_audit")
    BaseResponse AuditVehicle(1:AuditVehicleReq req)(api.post="/auth/manager/vehicle_audit")
}
