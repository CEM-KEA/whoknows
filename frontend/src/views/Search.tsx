import { useState } from "react";
import PageLayout from "../components/PageLayout";
import { apiPost } from "../utils/apiUtils";
import type { ISearchRequest, ISearchResponse } from "../types/types";

function Search() {
  const [search, setSearch] = useState("");
  const [searchResponse, setSearchResponse] = useState<ISearchResponse | null>(null);

  const handleSearch = () => {
    const searchBody: ISearchRequest = {
      q: search,
      language: "en"
    };
    apiPost<ISearchRequest, ISearchResponse>("/search", searchBody)
      .then((response) => {
        setSearchResponse(response);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return (
    <PageLayout>
      <div className="flex gap-2">
        <input
          id="search" // for playwright test
          type="search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSearch()}
          placeholder="Search for a programming language..."
          className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl"
        />
        <button
          className="bg-blue-500 text-white p-2 rounded ml-2 font-semibold hover:brightness-90 transition-colors duration-300"
          onClick={handleSearch}
        >
          Search
        </button>
      </div>
      <div>
        {searchResponse?.data && searchResponse.data.length > 0 && (
          <div className="mt-4">
            <ul
              id="search-results" // for playwright test
            >
              {searchResponse.data.map((result, i) => (
                <li
                  key={i}
                  className="list-none border-b p-1"
                >
                  <a
                    href={result.url}
                    className="text-blue-500 hover:underline font-semibold text-lg"
                    target="_blank"
                    title={result.url}
                  >
                    {result.title}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        )}
        {searchResponse?.data && searchResponse.data.length === 0 && (
          <div className="mt-4 italic text-gray-400 text-center text-lg">No results found...</div>
        )}
      </div>
    </PageLayout>
  );
}

export default Search;
