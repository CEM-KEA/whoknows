export interface ISearchResponse {
  data: {
    title: string;
    url: string;
    content: string;
    id: number;
  }[];
}
