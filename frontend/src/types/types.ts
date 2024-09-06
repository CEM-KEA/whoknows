export interface ISearchRequest {
  q: string;
  language: string | null; // language code e.g. en
}

export interface ISearchResponse {
  data: {
    title: string;
    url: string;
  }[];
}

export interface IStandardResponse {
  data: object;
}

export interface IRequestValidationError {
  statusCode: 422;
  message: string | null;
}

export interface IAuthResponse {
  statusCode: number;
  message: string;
}
