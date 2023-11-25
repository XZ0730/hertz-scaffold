namespace go base

struct BaseRequest{
}

struct BaseResponse{
    1: required i64 code,
    2: required string message,
}

struct LoginRequest{
    1: required string username,
    2: required string password,
}

struct LoginResponse{
    1: BaseResponse base,
    2: string token,
}

struct RegisterReq{
    1: required string username,
    2: required string password,
}
// 更新信息
struct UpdateUserInfoReq{
    1: required string name,
    2: required i64 age,
    3: required string sex,
    4: required string card_id,
    5: required string phone
}

struct VehicleModel{
    1: string vehicle_number
    2: string brand
    3: string model
    4: list<ParkSpaceKV> parklist
}

struct ParkSpaceKV{
    1: i64 park_id
    2: string garage_id
}

struct GetMyVehicleResponse{
    1: i64 code
    2: string msg
    3: list<VehicleModel> vlist
    4: i64 total
}
struct GetIdleParkSpaceReq{
    1: string garage_id
}
struct GetIdleParkSpaceResponse{
    1: i64 code 
    2: string msg 
    3: list<ParkSpaceKV> parklist
    4: i64 total
}
// 上传车牌照片 以及多张辅助照片
struct PostVehicleAuditReq {
    1: string vehicle_number
    2: string brand
    3: string model
    4: string description
    5: string garage_id
    6: i64 park_id
}

struct InParkRequest{
    1: string vehicle_number
    2: string garage_id
}

struct ParkRequest{
    1: required string vehicle_number
    2: required string garage_id
    3: required i64 type //入库还是出库 0 入库 1 出库
}

struct ParkResponse{
    1: i64 code
    2: string msg
    3: ParkSpaceKV psv
}


service BaseService {
    BaseResponse Ping(1:BaseRequest req)(api.get="/ping")

    BaseResponse Register(1:RegisterReq req)(api.post="/auth/register/pwd")
    LoginResponse Login(1:LoginRequest req)(api.post="/auth/login/pwd")
    
    BaseResponse PutUserInfo(1:UpdateUserInfoReq req)(api.put="/auth/user/")
    GetMyVehicleResponse GetMyVehicles(1:BaseRequest req)(api.get="/auth/user/vehicles")
    GetIdleParkSpaceResponse GetIdleParkSpace(1:GetIdleParkSpaceReq req)(api.get="/auth/user/park_space")

    BaseResponse PostVehicleAudit(1:PostVehicleAuditReq req)(api.post="/auth/user/audit")

    // 入库 / 出库
    ParkResponse Park(ParkRequest req)(api.post="/auth/user/park")
}