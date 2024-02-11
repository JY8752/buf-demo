import { useState } from "react";
import "./App.css";
import { createConnectTransport } from "@connectrpc/connect-web";
import { createPromiseClient } from "@connectrpc/connect";

import { WeatherService } from "@buf/jyapp_weather.connectrpc_es/jyapp/weather/v1/weather_connect";

function App() {
  const [weather, setWeather] = useState("");
  const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
  });

  const client = createPromiseClient(WeatherService, transport);

  return (
    <>
      <p>{weather}</p>
      <button
        onClick={async () => {
          const res = await client.getWeather({
            latitude: 1.0,
            longitude: 1.0,
          });
          setWeather(res.toJsonString());
        }}
      >
        Call
      </button>
    </>
  );
}

export default App;
