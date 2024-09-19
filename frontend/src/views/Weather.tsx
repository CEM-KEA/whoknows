import { useEffect, useState } from "react";
import type { IWeatherData } from "../types/weather.types";
import { apiGet } from "../utils/apiUtils";
import { IStandardResponse } from "../types/types";
import { formatTemperature, weatherToEmoji } from "../helpers/weatherHelpers";
import PageLayout from "../components/PageLayout";
import { toTitleCase } from "../helpers/formatHelpers";
import { MdLocationPin } from "react-icons/md";
import LoadingSpinner from "../components/LoadingSpinner";

function Weather() {
  const [loading, setLoading] = useState(true);
  const [weather, setWeather] = useState<IWeatherData | null>(null);

  useEffect(() => {
    apiGet<IStandardResponse<IWeatherData>>("/weather")
      .then((res: IStandardResponse<IWeatherData>) => setWeather(res.data))
      .catch((error) => {
        console.error(error);
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <PageLayout>
      <div className="flex items-center justify-center mt-8">
        {loading && (
          <div className="mt-32">
            <LoadingSpinner size={100} />
          </div>
        )}
        {weather && (
          <div
            className="border rounded p-8 bg-blue-200 bg-opacity-40 flex flex-col gap-1 shadow-lg"
            id="weather-content" // for playwright test
          >
            <div className="text-[15em]">{weatherToEmoji[weather.weather[0].description]}</div>
            <div className="text-5xl font-semibold text-center">
              {formatTemperature(weather.main.temp)}
            </div>
            <div className="text-2xl text-center">
              {toTitleCase(weather.weather[0].description)}
            </div>
            <div className="text-center text-xl">
              <MdLocationPin className="inline-block mb-1 text-blue-500" /> {weather.name},{" "}
              {weather.sys.country}
            </div>
          </div>
        )}
      </div>
    </PageLayout>
  );
}

export default Weather;
