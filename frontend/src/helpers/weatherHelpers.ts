export const weatherToEmoji: { [key: string]: string } = {
  "clear sky": "☀️",
  "few clouds": "🌤",
  "scattered clouds": "⛅️",
  "broken clouds": "🌥",
  "overcast clouds": "☁️",
  mist: "🌫️",
  fog: "🌫️",
  rain: "🌧️",
  "light rain": "🌦",
  "heavy rain": "🌧️",
  drizzle: "🌦",
  thunderstorm: "⛈️",
  snow: "🌨️",
  "light snow": "🌨️",
  "heavy snow": "❄️",
  hail: "🌨️",
  sleet: "🌧️❄️",
  windy: "💨",
  "shower rain": "🌧️"
};

export function formatTemperature(tempInKelvin: number): string {
  return (tempInKelvin - 273.15).toFixed(1) + " °C";
}
