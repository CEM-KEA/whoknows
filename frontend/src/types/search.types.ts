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
