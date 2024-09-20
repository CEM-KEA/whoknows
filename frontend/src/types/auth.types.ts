export interface IAuthResponse {
  statusCode: number;
  message: string;
}

export interface ILoginRequest {
  email: string;
  password: string;
}

export interface ILoginSession {
  email: string;
}

export interface ILoginResponse {
  token: string;
}

export interface IRegisterRequest {
  username: string;
  email: string;
  password: string;
}