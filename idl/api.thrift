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

service BaseService {
    BaseResponse Ping(1:BaseRequest req)(api.get="/ping")

    BaseResponse Register(1:RegisterReq req)(api.post="/auth/register/pwd")
    LoginResponse Login(1:LoginRequest req)(api.post="/auth/login/pwd")
}