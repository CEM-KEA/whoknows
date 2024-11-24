import { useEffect, useRef, useState } from "react";
import PageLayout from "../components/PageLayout";
import { apiGet } from "../utils/apiUtils";
import type { ISearchResponse } from "../types/search.types";
import LoadingSpinner from "../components/LoadingSpinner";
import toast from "react-hot-toast";
import { useLocation, useNavigate } from "react-router-dom";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function Search() {
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [searchResponse, setSearchResponse] = useState<ISearchResponse | null>(null);
  const query = useQuery();
  const navigate = useNavigate();
  const initialSearchPerformed = useRef(false); // useRef to persist the value across re-renders

  useEffect(() => {
    const queryParam = query.get("q");
    if (queryParam && !initialSearchPerformed.current) {
      setSearch(queryParam);
      performSearch(queryParam);
      initialSearchPerformed.current = true;
    }
  }, [query]);

  const performSearch = (query: string) => {
    setLoading(true);
    apiGet<ISearchResponse>("/search?q=" + query)
      .then((response) => {
        setSearchResponse(response);
      })
      .catch((error) => {
        toast.error(error.message);
      })
      .finally(() => setLoading(false));
  };

  const handleSearch = () => {
    if (search.trim() === "") {
      setSearchResponse(null);
      navigate(``);
      return;
    }
    // Encode the search query to handle special characters for URL
    navigate(`?q=${encodeURIComponent(search)}`);
    performSearch(search);
  };

  // highlights the part of the text that matches the search query
  const highlightText = (text: string, highlight: string) => {
    const regex = new RegExp(`(${highlight})`, "gi");
    return text.split(regex).map((part, index) =>
      regex.test(part) ? (
        <span
          key={part + index}
          className="bg-yellow-200"
        >
          {part}
        </span>
      ) : (
        part
      )
    );
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
        {loading && (
          <div className="mt-32 w-full flex justify-center">
            <LoadingSpinner size={100} />
          </div>
        )}
        {searchResponse === null && !loading && (
          <div className="mt-4 italic text-gray-400 text-center text-lg">
            Search for something...
          </div>
        )}
        {searchResponse?.data && searchResponse.data.length > 0 && (
          <div className="mt-4">
            <ul
              id="search-results" // for playwright test
            >
              {searchResponse.data.map((result, i) => (
                <li
                  key={result.title + i}
                  className="list-none border-b p-1"
                >
                  <a
                    href={result.url}
                    className="text-blue-500 hover:underline font-semibold text-lg"
                    target="_blank"
                    title={result.url}
                  >
                    {result.title}
                    <p className="text-gray-400 italic">{highlightText(result.content, search)}</p>
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
