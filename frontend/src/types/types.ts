export interface IStandardResponse<T> {
  data: T;
}

export interface IRequestValidationError {
  statusCode: 422;
  message: string | null;
}
