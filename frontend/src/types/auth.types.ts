export interface IAuthResponse {
  statusCode: number;
  message: string;
}

export interface ILoginRequest {
  username: string;
  password: string;
}

export interface ILoginSession {
  email: string;
}

export interface ILoginResponse {
  token: string;
  require_password_change: boolean;
}

export interface IRegisterRequest {
  username: string;
  email: string;
  password: string;
  password2: string;
}

export interface IChangePasswordRequest {
  username: string;
  old_password: string;
  new_password: string;
  repeat_new_password: string;
}
